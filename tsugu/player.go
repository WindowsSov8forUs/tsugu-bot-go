package tsugu

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/WindowsSov8forUs/tsugu-bot-go/config"
	log "github.com/WindowsSov8forUs/tsugu-bot-go/mylog"
	"github.com/syndtr/goleveldb/leveldb"
)

// 用户数据库
type UserDataDB struct {
	DB    *leveldb.DB
	mutex sync.Mutex
}

func (db *UserDataDB) getUserData(platform, userID string) (*UserData, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	key := []byte(userID + platform)
	data, err := db.DB.Get(key, nil)
	if err != nil {
		return nil, err
	}
	var userData = &UserData{}
	err = json.Unmarshal(data, userData)
	if err != nil {
		return nil, err
	}
	return userData, nil
}

func (db *UserDataDB) saveUserData(userData *UserData) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	key := []byte(userData.UserID + userData.Platform)
	data, err := json.Marshal(userData)
	if err != nil {
		return err
	}
	return db.DB.Put(key, data, nil)
}

func GetUserData(platform, userID string) (*ResponseUserData, error) {
	userData, err := userDataBaseInstance.getUserData(platform, userID)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return &ResponseUserData{
				Status: "failure",
				Data: &UserData{
					UserID:        userID,
					Platform:      platform,
					ServerMode:    3,
					DefaultServer: []int{3},
					Car:           true,
					GameIDs:       []*GameID{},
					VerifyCode:    "",
				},
			}, nil
		} else {
			return nil, err
		}
	}
	return &ResponseUserData{
		Status: "success",
		Data:   userData,
	}, nil
}

func SaveUserData(userData *UserData) error {
	return userDataBaseInstance.saveUserData(userData)
}

var userDataBaseInstance *UserDataDB

type ResponseUserData struct {
	Status string    `json:"status"`
	Data   *UserData `json:"data"`
}

func (data *ResponseUserData) Stringify() string {
	return fmt.Sprintf("{Status: %s, Data: %s}", data.Status, data.Data.Stringify())
}

type UserData struct {
	UserID        string    `json:"user_id"`
	Platform      string    `json:"platform"`
	ServerMode    int       `json:"server_mode,omitempty"`
	DefaultServer []int     `json:"default_server,omitempty"`
	Car           bool      `json:"car,omitempty"`
	GameIDs       []*GameID `json:"game_ids,omitempty"`
	VerifyCode    string    `json:"verify_code,omitempty"`
	ServerList    []struct {
		PlayerID int `json:"playerId,omitempty"`
	} `json:"server_list,omitempty"`
}

func (data *UserData) Stringify() string {
	return fmt.Sprintf("{UserID: %s, Platform: %s, ServerMode: %d, DefaultServer: %v, Car: %v, GameIDs: %v, VerifyCode: %s, ServerList: %v}", data.UserID, data.Platform, data.ServerMode, data.DefaultServer, data.Car, data.GameIDs, data.VerifyCode, data.ServerList)
}

type GameID struct {
	Server int `json:"server"`
	GameID int `json:"game_id"`
}

type ResponsePlayer struct {
	Data struct {
		Profile struct {
			Introduction string `json:"introduction"`
			MainUserDeck struct {
				DeckName string `json:"deckName"`
			} `json:"mainUserDeck"`
		} `json:"profile"`
	} `json:"data"`
}

func (data *ResponsePlayer) Stringify() string {
	return fmt.Sprintf("{Data: %+v}", data.Data)
}

type ResponseRemote struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

type VerifyData struct {
	VerifyCode int `json:"verifyCode"`
}

func DataBase(conf *config.Config) error {
	path := conf.Tsugu.UserDataBasePath
	if path == "" {
		path = "data/user"
	}
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return err
	}
	userDataBaseInstance = &UserDataDB{
		DB: db,
	}
	return nil
}

func remoteGetUserData(platform, userID string, conf *config.Config) (*ResponseUserData, error) {
	api := "/user/getUserData"
	data := &RequestUser{
		Platform: platform,
		UserID:   userID,
	}

	response, err := requestPostUser(api, data, conf)
	if err != nil {
		return nil, err
	}

	var res = &ResponseUserData{}
	err = json.Unmarshal(response, res)
	if err != nil {
		return nil, err
	}
	log.Debugf("<Tsugu> 获取到用户数据: %s", res.Stringify())

	return res, nil
}

func serverToIngeter(server string) int {
	switch server {
	case "jp":
		return 0
	case "en":
		return 1
	case "tw":
		return 2
	case "cn":
		return 3
	case "kr":
		return 4
	case "日服":
		return 0
	case "国际服":
		return 1
	case "台服":
		return 2
	case "国服":
		return 3
	case "韩服":
		return 4
	default:
		return -1
	}
}

func serverToShortName(server int) string {
	switch server {
	case 0:
		return "jp"
	case 1:
		return "en"
	case 2:
		return "tw"
	case 3:
		return "cn"
	case 4:
		return "kr"
	default:
		return ""
	}
}

func serverToFullName(server int) string {
	switch server {
	case 0:
		return "日服"
	case 1:
		return "国际服"
	case 2:
		return "台服"
	case 3:
		return "国服"
	case 4:
		return "韩服"
	default:
		return ""
	}
}

func playerStatus(platform, userID string, server, index int, conf *config.Config) ([]*ResponseData, error) {
	var responseData []*ResponseData
	var playerID int

	userData, err := GetUserData(userID, platform)
	if err != nil {
		return nil, err
	}
	gameIDs := userData.Data.GameIDs
	if len(gameIDs) == 0 {
		return TextResponse("未绑定玩家，请发送 绑定玩家 进行绑定"), nil
	}
	if server < 0 && index <= 0 {
		// 未指定服务器和条目数，先查找默认服务器
		var found = false
		server = userData.Data.ServerMode
		for _, gameID := range gameIDs {
			if gameID.Server == server {
				found = true
				playerID = gameID.GameID
				responseData = TextResponse(fmt.Sprintf("已查找默认服务器 %s 的记录", serverToFullName(server)))
			}
		}
		if !found {
			// 查找第一个记录
			length := len(gameIDs)
			responseData = TextResponse(fmt.Sprintf("未在 %v 条记录中找到 %s 的记录，已查找第一个记录", length, serverToFullName(server)))
			server = gameIDs[0].Server
			playerID = gameIDs[0].GameID
		}
	} else if index > 0 {
		// 指定了条目数
		if index > len(gameIDs) {
			return TextResponse(fmt.Sprintf("总共绑定了 %v 条记录。", len(gameIDs))), nil
		}
		server = gameIDs[index-1].Server
		playerID = gameIDs[index-1].GameID
		responseData = TextResponse(fmt.Sprintf("已查找第 %v 条记录", index))
	} else {
		// 指定了服务器
		var found = false
		for _, gameID := range gameIDs {
			if gameID.Server == server {
				found = true
				playerID = gameID.GameID
				responseData = TextResponse(fmt.Sprintf("已查找服务器 %s 的记录", serverToFullName(server)))
			}
		}
		if !found {
			return TextResponse("未找到记录，请检查是否绑定过此服务器"), nil
		}
	}
	result, err := ApiV2Backend("player", fmt.Sprintf("%v", playerID), nil, server, conf)
	if err != nil {
		return nil, err
	}
	responseData = append(responseData, result...)
	return responseData, nil
}

func setCarForward(platform, userID string, forward bool) []*ResponseData {
	userData, err := GetUserData(userID, platform)
	if err != nil {
		log.Errorf("<Tsugu> 获取用户数据失败: %v", err)
		return TextResponse("获取用户数据失败")
	}
	userData.Data.Car = forward
	err = SaveUserData(userData.Data)
	if err != nil {
		log.Errorf("<Tsugu> 保存用户数据失败: %v", err)
		return TextResponse("保存用户数据失败")
	}
	return TextResponse("车牌转发设置成功")
}

func setDefaultServer(platform, userID, serverList string) []*ResponseData {
	userData, err := GetUserData(userID, platform)
	if err != nil {
		log.Errorf("<Tsugu> 获取用户数据失败: %v", err)
		return TextResponse("获取用户数据失败")
	}
	servers := strings.Split(serverList, " ")
	var defaultServer []int
	for _, server := range servers {
		serverID := serverToIngeter(server)
		if serverID == -1 {
			return TextResponse(fmt.Sprintf("未找到名为 %s 的服务器", server))
		}
		defaultServer = append(defaultServer, serverID)
	}
	userData.Data.DefaultServer = defaultServer
	err = SaveUserData(userData.Data)
	if err != nil {
		log.Errorf("<Tsugu> 保存用户数据失败: %v", err)
		return TextResponse("保存用户数据失败")
	}
	return TextResponse("默认服务器设置成功")
}

func setServerMode(platform, userID, server string) []*ResponseData {
	userData, err := GetUserData(userID, platform)
	if err != nil {
		log.Errorf("<Tsugu> 获取用户数据失败: %v", err)
		return TextResponse("获取用户数据失败")
	}
	serverID := serverToIngeter(server)
	if serverID == -1 {
		return TextResponse(fmt.Sprintf("未找到名为 %s 的服务器", server))
	}
	userData.Data.ServerMode = serverID
	err = SaveUserData(userData.Data)
	if err != nil {
		log.Errorf("<Tsugu> 保存用户数据失败: %v", err)
		return TextResponse("保存用户数据失败")
	}
	return TextResponse("默认服务器设置成功")
}

func bindPlayer(platform, userID string) []*ResponseData {
	userData, err := GetUserData(userID, platform)
	if err != nil {
		log.Errorf("<Tsugu> 获取用户数据失败: %v", err)
		return TextResponse("获取用户数据失败")
	}
	bindCount := len(userData.Data.GameIDs)

	var verifyCode string
	var verifyNum int
	for {
		rand.New(rand.NewSource(time.Now().UnixNano()))
		verifyNum = rand.Intn(90000) + 10000
		verifyCode = strconv.Itoa(verifyNum)
		if !strings.Contains(verifyCode, "64") && !strings.Contains(verifyCode, "89") {
			break
		}
	}

	responseData := TextResponse(
		fmt.Sprintf(
			"正在绑定第%v条记录，请将 评论（个性签名）或 当前使用的乐队编队名称 改为\n%s\n稍等片刻等待同步后，发送\n验证 玩家ID 来完成本次身份验证\n例如：验证 10000xxxx 国服",
			bindCount+1,
			verifyCode,
		),
	)

	userData.Data.VerifyCode = verifyCode
	err = SaveUserData(userData.Data)
	if err != nil {
		log.Errorf("<Tsugu> 保存用户数据失败: %v", err)
		return TextResponse("保存用户数据失败")
	}
	return responseData
}

func unbindPlayer(platform, userID string) []*ResponseData {
	userData, err := GetUserData(userID, platform)
	if err != nil {
		log.Errorf("<Tsugu> 获取用户数据失败: %v", err)
		return TextResponse("获取用户数据失败")
	}
	length := len(userData.Data.GameIDs)
	responseData := TextResponse(
		fmt.Sprintf(
			`当前有 %v 个记录，发送"解除绑定 %v"来获解除第%v个记录，以此类推`,
			length,
			length,
			length,
		),
	)
	return responseData
}

func bindVerify(platform, userID string, playerID, server int, conf *config.Config) []*ResponseData {
	userData, err := GetUserData(userID, platform)
	if err != nil {
		log.Errorf("<Tsugu> 获取用户数据失败: %v", err)
		return TextResponse("获取用户数据失败")
	}
	if server < 0 {
		server = userData.Data.ServerMode
	}
	if userData.Data.VerifyCode == "" {
		return TextResponse("请先获取验证代码")
	}
	gameIDs := userData.Data.GameIDs
	for _, gameID := range gameIDs {
		if gameID.GameID == playerID && gameID.Server == server {
			return TextResponse("请勿重复绑定")
		}
	}
	serverShort := serverToShortName(server)

	playerURL := fmt.Sprintf("https://bestdori.com/api/player/%s/%v?mode=2", serverShort, playerID)
	var transport *http.Transport
	if conf.Tsugu.VerifyPlayer.UseProxy {
		proxyURL, _ := url.Parse(conf.Tsugu.Proxy)
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	} else {
		transport = &http.Transport{
			Proxy: nil,
		}
	}
	client := &http.Client{
		Transport: transport,
	}
	if conf.Tsugu.Timeout > 0 {
		client.Timeout = time.Duration(conf.Tsugu.Timeout) * time.Second
	}
	response, err := client.Get(playerURL)
	if err != nil {
		log.Errorf("<Tsugu> 获取玩家信息失败: %v", err)
		return TextResponse("获取玩家数据失败")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Errorf("<Tsugu> 获取玩家信息失败: %v", response.StatusCode)
		return TextResponse(fmt.Sprintf("获取玩家数据失败，状态码：%v", response.StatusCode))
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Errorf("<Tsugu> 读取玩家信息失败: %v", err)
		return TextResponse("读取玩家信息失败")
	}

	var playerData = &ResponsePlayer{}
	err = json.Unmarshal(body, playerData)
	if err != nil {
		log.Errorf("<Tsugu> 解析玩家信息失败: %v", err)
		return TextResponse("解析玩家信息失败")
	}
	log.Debugf("<Tsugu> 获取玩家信息: %s", playerData.Stringify())
	verifyCode := userData.Data.VerifyCode
	if verifyCode != playerData.Data.Profile.MainUserDeck.DeckName && verifyCode != playerData.Data.Profile.Introduction {
		return TextResponse("验证失败，签名或者乐队编队名称与验证代码不匹配，可以检查后再次尝试（无需重复发送绑定玩家）")
	}
	gameID := &GameID{
		Server: server,
		GameID: playerID,
	}
	userData.Data.GameIDs = append(userData.Data.GameIDs, gameID)
	userData.Data.VerifyCode = ""
	err = SaveUserData(userData.Data)
	if err != nil {
		log.Errorf("<Tsugu> 保存用户数据失败: %v", err)
		return TextResponse("保存用户数据失败")
	}
	return TextResponse("绑定成功！现在可以使用“玩家状态”来查询玩家信息了")
}

func unbindVerify(platform, userID string, index int) []*ResponseData {
	userData, err := GetUserData(userID, platform)
	if err != nil {
		log.Errorf("<Tsugu> 获取用户数据失败: %v", err)
		return TextResponse("获取用户数据失败")
	}
	gameIDs := userData.Data.GameIDs
	if index < 1 || index > len(gameIDs) {
		length := len(gameIDs)
		return TextResponse(
			fmt.Sprintf(
				"解绑失败，当前有 %v 个记录，发送“解除绑定 %v”来获解除第%v个记录，以此类推",
				length,
				length,
				length,
			),
		)
	}
	userData.Data.GameIDs = append(gameIDs[:index-1], gameIDs[index:]...)
	err = SaveUserData(userData.Data)
	if err != nil {
		log.Errorf("<Tsugu> 保存用户数据失败: %v", err)
		return TextResponse("保存用户数据失败")
	}
	return TextResponse("解绑成功")
}

func remoteBindPlayer(platform, userID string, server int, bindType bool, conf *config.Config) (*ResponseRemote, error) {
	requestData := &RequestUser{
		Platform: platform,
		UserID:   userID,
		Server:   server,
		BindType: bindType,
	}
	response, err := requestPostUser("/user/bindPlayerRequest", requestData, conf)
	if err != nil {
		return nil, err
	}
	var res = &ResponseRemote{}
	err = json.Unmarshal(response, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func remoteBindVerify(platform, userID string, playerID, server int, bindType bool, conf *config.Config) (*ResponseRemote, error) {
	requestData := &RequestUser{
		Platform: platform,
		UserID:   userID,
		Server:   server,
		PlayerID: playerID,
		BindType: bindType,
	}
	response, err := requestPostUser("/user/bindPlayerVerification", requestData, conf)
	if err != nil {
		return nil, err
	}
	var res = &ResponseRemote{}
	err = json.Unmarshal(response, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func remotePlayerStatus(platform, userID string, server int, conf *config.Config) []*ResponseData {
	userData, err := remoteGetUserData(platform, userID, conf)
	if err != nil {
		log.Errorf("<Tsugu> 获取用户数据失败: %v", err)
		return TextResponse("获取用户数据失败")
	}
	if userData.Status != "success" {
		log.Warnf("<Tsugu> 获取用户数据失败: %s:%s", platform, userID)
		return TextResponse("获取用户数据失败")
	}
	if server < 0 {
		server = userData.Data.ServerMode
	}
	playerID := userData.Data.ServerList[server].PlayerID
	if playerID == 0 {
		return TextResponse("未绑定玩家，请使用 绑定玩家 进行绑定")
	}
	response, err := ApiV2Backend("player", fmt.Sprintf("%v", playerID), nil, server, conf)
	if err != nil {
		log.Errorf("<Tsugu> 获取玩家信息失败: %v", err)
		return TextResponse("获取玩家信息失败")
	}
	return response
}

func remoteSetCarForward(platform, userID string, forward bool, conf *config.Config) (*ResponseRemote, error) {
	requestData := &RequestUser{
		Platform: platform,
		UserID:   userID,
		Status:   forward,
	}
	response, err := requestPostUser("/user/changeUserData/setCarForwarding", requestData, conf)
	if err != nil {
		return nil, err
	}
	var res = &ResponseRemote{}
	err = json.Unmarshal(response, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func remoteSetDefaultServer(platform, userID, serverList string, conf *config.Config) (*ResponseRemote, error) {
	requestData := &RequestUser{
		Platform: platform,
		UserID:   userID,
		Text:     serverList,
	}
	response, err := requestPostUser("/user/changeUserData/setDefaultServer", requestData, conf)
	if err != nil {
		return nil, err
	}
	var res = &ResponseRemote{}
	err = json.Unmarshal(response, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func remoteSetServerMode(platform, userID, server string, conf *config.Config) (*ResponseRemote, error) {
	requestData := &RequestUser{
		Platform: platform,
		UserID:   userID,
		Text:     server,
	}
	response, err := requestPostUser("/user/changeUserData/setServerMode", requestData, conf)
	if err != nil {
		return nil, err
	}
	var res = &ResponseRemote{}
	err = json.Unmarshal(response, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
