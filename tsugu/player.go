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

// 本地用户数据操作

type verifyData struct {
	Server     ServerId
	VerifyCode string
}

type localVerifyCode map[string]*verifyData

var verifyCodeMap = make(localVerifyCode)

// 用户数据库
type UserDataDB struct {
	DB    *leveldb.DB
	mutex sync.Mutex
}

type GameID struct {
	Server ServerId `json:"server"`
	GameID int      `json:"game_id"`
}

type UserData struct {
	UserID        string     `json:"user_id"`
	Platform      string     `json:"platform"`
	ServerMode    ServerId   `json:"server_mode,omitempty"`
	DefaultServer []ServerId `json:"default_server,omitempty"`
	Car           bool       `json:"car,omitempty"`
	GameIDs       []*GameID  `json:"game_ids,omitempty"`
}

type ResponseBestdoriPlayer struct {
	Data struct {
		Profile struct {
			Introduction string `json:"introduction"`
			MainUserDeck struct {
				DeckName string `json:"deckName"`
			} `json:"mainUserDeck"`
		} `json:"profile"`
	} `json:"data"`
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

func GetUserData(platform, userID string) (*UserData, error) {
	userData, err := userDataBaseInstance.getUserData(platform, userID)
	if err != nil {
		if err == leveldb.ErrNotFound {
			userData := &UserData{
				UserID:        userID,
				Platform:      platform,
				ServerMode:    CN,
				DefaultServer: []ServerId{CN, JP},
				Car:           true,
				GameIDs:       []*GameID{},
			}
			err = SaveUserData(userData)
			if err != nil {
				return nil, err
			}
			return userData, nil
		} else {
			return nil, err
		}
	}
	return userData, nil
}

func SaveUserData(userData *UserData) error {
	return userDataBaseInstance.saveUserData(userData)
}

var userDataBaseInstance *UserDataDB

func bindPlayer(platform, userID string, server ServerId) []*ResponseData {
	userData, err := GetUserData(userID, platform)
	if err != nil {
		log.Errorf("<Tsugu> 获取用户数据失败: %v", err)
		return TextResponse("获取用户数据失败")
	}
	bindCount := len(userData.GameIDs)

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

	verifyCodeMap[fmt.Sprintf("%s:%s", platform, userID)] = &verifyData{
		Server:     server,
		VerifyCode: verifyCode,
	}
	err = SaveUserData(userData)
	if err != nil {
		log.Errorf("<Tsugu> 保存用户数据失败: %v", err)
		return TextResponse("保存用户数据失败")
	}
	return responseData
}

func unbindPlayer(platform, userID string, index int) []*ResponseData {
	userData, err := GetUserData(userID, platform)
	if err != nil {
		log.Errorf("<Tsugu> 获取用户数据失败: %v", err)
		return TextResponse("获取用户数据失败")
	}
	length := len(userData.GameIDs)
	if index < 1 || index > length {
		if length == 0 {
			return TextResponse("目前没有绑定记录")
		}
		return TextResponse(
			fmt.Sprintf(
				"参数错误，当前有 %v 个记录，发送“解除绑定 %v”来获解除第%v个记录，以此类推",
				length,
				length,
				length,
			),
		)
	}

	// 本地用户数据库将不需求解绑验证

	responseData := TextResponse(
		fmt.Sprintf(
			"已解绑第%v条记录，该记录为于 %s 的 %v",
			index,
			serverIdToFullName(userData.GameIDs[index-1].Server),
			userData.GameIDs[index-1].GameID,
		),
	)

	userData.GameIDs = append(userData.GameIDs[:index-1], userData.GameIDs[index:]...)
	err = SaveUserData(userData)
	if err != nil {
		log.Errorf("<Tsugu> 保存用户数据失败: %v", err)
		return TextResponse("保存用户数据失败")
	}
	return responseData
}

func bindVerify(platform, userID string, playerID int, conf *config.Config) ([]*ResponseData, bool) {
	userData, err := GetUserData(userID, platform)
	if err != nil {
		log.Errorf("<Tsugu> 获取用户数据失败: %v", err)
		return TextResponse("获取用户数据失败"), true
	}
	if verifyCodeMap[fmt.Sprintf("%s:%s", platform, userID)] == nil {
		return TextResponse("请先获取验证代码"), false
	}
	verify := verifyCodeMap[fmt.Sprintf("%s:%s", platform, userID)]
	gameIDs := userData.GameIDs
	for _, gameID := range gameIDs {
		if gameID.GameID == playerID && gameID.Server == verify.Server {
			return TextResponse("请勿重复绑定"), false
		}
	}
	serverShort := serverIdToShortName(verify.Server)

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
		return TextResponse("获取玩家数据失败"), true
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Errorf("<Tsugu> 获取玩家信息失败: %v", response.StatusCode)
		return TextResponse(fmt.Sprintf("获取玩家数据失败，状态码：%v", response.StatusCode)), true
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Errorf("<Tsugu> 读取玩家信息失败: %v", err)
		return TextResponse("读取玩家信息失败"), true
	}

	var playerData = &ResponseBestdoriPlayer{}
	err = json.Unmarshal(body, playerData)
	if err != nil {
		log.Errorf("<Tsugu> 解析玩家信息失败: %v", err)
		return TextResponse("解析玩家信息失败"), true
	}
	log.Debugf("<Tsugu> 获取玩家信息: %v", playerData)
	verifyCode := verify.VerifyCode
	if verifyCode != playerData.Data.Profile.MainUserDeck.DeckName && verifyCode != playerData.Data.Profile.Introduction {
		return TextResponse("验证失败，签名或者乐队编队名称与验证代码不匹配，可以检查后再次尝试（无需重复发送绑定玩家）"), true
	}
	gameID := &GameID{
		Server: verify.Server,
		GameID: playerID,
	}
	userData.GameIDs = append(userData.GameIDs, gameID)
	delete(verifyCodeMap, fmt.Sprintf("%s:%s", platform, userID))
	err = SaveUserData(userData)
	if err != nil {
		log.Errorf("<Tsugu> 保存用户数据失败: %v", err)
		return TextResponse("保存用户数据失败"), true
	}
	return TextResponse("绑定成功！现在可以使用“玩家状态”来查询玩家信息了"), false
}

// 后端用户数据操作

type UserApiFailedError struct {
	Data string
}

func (err *UserApiFailedError) Error() string {
	return err.Data
}

type TsuguUserServer struct {
	PlayerID      int           `json:"playerId,omitempty"`
	BindingStatus BindingStatus `json:"bindingStatus,omitempty"`
	VerifyCode    int           `json:"verifyCode,omitempty"`
}

type TsuguUser struct {
	Id            string             `json:"_id,omitempty"`
	UserID        string             `json:"user_id"`
	Platform      string             `json:"platform"`
	ServerMode    ServerId           `json:"server_mode,omitempty"`
	DefaultServer []ServerId         `json:"default_server,omitempty"`
	Car           bool               `json:"car,omitempty"`
	ServerList    []*TsuguUserServer `json:"server_list,omitempty"`
}

type TsuguUserUpdate struct {
	UserID        string             `json:"user_id,omitempty"`
	Platform      string             `json:"platform,omitempty"`
	ServerMode    *ServerId          `json:"server_mode,omitempty"`
	DefaultServer []ServerId         `json:"default_server,omitempty"`
	Car           bool               `json:"car,omitempty"`
	ServerList    []*TsuguUserServer `json:"server_list,omitempty"`
}

type ResponseUserFailed struct {
	Status Status `json:"status"`
	Data   string `json:"data"`
}

type ResponseGetUserData struct {
	Status Status     `json:"status"`
	Data   *TsuguUser `json:"data"`
}

type ResponseChangeUserData struct {
	Status Status `json:"status"`
	Data   string `json:"data,omitempty"`
}

type VerifyCode struct {
	VerifyCode int `json:"verifyCode"`
}

type ResponseBindPlayerRequest struct {
	Status Status     `json:"status"`
	Data   VerifyCode `json:"data"`
}

type ResponseBindPlayerVerification struct {
	Status Status `json:"status"`
	Data   string `json:"data"`
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

func RemoteGetUserData(platform, userID string, conf *config.Config) (*ResponseGetUserData, error) {
	data := &RequestUserData{
		Platform: platform,
		UserID:   userID,
	}

	response, err := requestPostUser(ApiUserGetUserData, data, conf)
	if err != nil {
		return nil, err
	}

	var res = &ResponseGetUserData{}
	err = json.Unmarshal(response, res)
	if err != nil {
		fmt.Printf("%s\n", err)
		var errRes = &ResponseUserFailed{}
		err = json.Unmarshal(response, &errRes)
		if err != nil {
			return nil, err
		} else {
			return nil, &UserApiFailedError{Data: errRes.Data}
		}
	}
	log.Debugf("<Tsugu> 获取到用户数据: %v", res.Data)

	return res, nil
}

func RemoteChangeUserData(platform, userID string, update *TsuguUserUpdate, conf *config.Config) (*ResponseChangeUserData, error) {
	data := &RequestUserData{
		Platform: platform,
		UserID:   userID,
		Update:   update,
	}

	response, err := requestPostUser(ApiUserChangeUserData, data, conf)
	if err != nil {
		return nil, err
	}

	var res = &ResponseChangeUserData{}
	err = json.Unmarshal(response, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func RemoteBindPlayerRequest(platform, userID string, server ServerId, bindType bool, conf *config.Config) (*ResponseBindPlayerRequest, error) {
	requestData := &RequestBindPlayerRequest{
		Platform: platform,
		UserID:   userID,
		Server:   server,
		BindType: bindType,
	}
	response, err := requestPostUser(ApiUserBindPlayerRequest, requestData, conf)
	if err != nil {
		return nil, err
	}
	var res = &ResponseBindPlayerRequest{}
	err = json.Unmarshal(response, res)
	if err != nil {
		var errRes = &ResponseUserFailed{}
		err = json.Unmarshal(response, &errRes)
		if err != nil {
			return nil, err
		} else {
			return nil, &UserApiFailedError{Data: errRes.Data}
		}
	}

	return res, nil
}

func RemoteBindVerification(platform, userID string, server ServerId, playerID int, bindType bool, conf *config.Config) (*ResponseBindPlayerVerification, error) {
	requestData := &RequestBindPlayerVerification{
		Platform: platform,
		UserID:   userID,
		Server:   server,
		PlayerID: playerID,
		BindType: bindType,
	}
	response, err := requestPostUser(ApiUserBindPlayerVerification, requestData, conf)
	if err != nil {
		return nil, err
	}
	var res = &ResponseBindPlayerVerification{}
	err = json.Unmarshal(response, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 统合用户数据操作

func getCarForward(platform, userID string, conf *config.Config) (bool, error) {
	if conf.Tsugu.UserDataBasePath == "" {
		response, err := RemoteGetUserData(platform, userID, conf)
		if err != nil {
			return false, err
		}
		return response.Data.Car, nil
	} else {
		userData, err := GetUserData(platform, userID)
		if err != nil {
			return false, err
		}
		return userData.Car, nil
	}
}

func setCarForward(platform, userID string, forward bool, conf *config.Config) (string, error) {
	if conf.Tsugu.UserDataBasePath == "" {
		update := &TsuguUserUpdate{
			Car: forward,
		}
		response, err := RemoteChangeUserData(platform, userID, update, conf)
		if err != nil {
			return "", err
		}
		if response.Status == STATUS_FAILED {
			return response.Data, fmt.Errorf("%s", response.Data)
		} else {
			return "", nil
		}
	} else {
		userData, err := GetUserData(platform, userID)
		if err != nil {
			return "", err
		}
		userData.Car = forward
		err = SaveUserData(userData)
		if err != nil {
			return "", err
		}
		return "", nil
	}
}

func getMainServer(platform, userID string, conf *config.Config) (ServerId, error) {
	if conf.Tsugu.UserDataBasePath == "" {
		response, err := RemoteGetUserData(platform, userID, conf)
		if err != nil {
			return -1, err
		}
		return response.Data.ServerMode, nil
	} else {
		userData, err := GetUserData(platform, userID)
		if err != nil {
			return -1, err
		}
		return userData.ServerMode, nil
	}
}

func setMainServer(platform, userID string, server ServerId, conf *config.Config) error {
	if conf.Tsugu.UserDataBasePath == "" {
		update := &TsuguUserUpdate{
			ServerMode: &server,
		}
		response, err := RemoteChangeUserData(platform, userID, update, conf)
		if err != nil {
			return err
		}
		if response.Status == STATUS_FAILED {
			return fmt.Errorf("%s", response.Data)
		}
		return nil
	} else {
		userData, err := GetUserData(platform, userID)
		if err != nil {
			return err
		}
		userData.ServerMode = server
		err = SaveUserData(userData)
		if err != nil {
			return err
		}
		return nil
	}
}

func getDefaultServers(platform, userID string, conf *config.Config) ([]ServerId, error) {
	if conf.Tsugu.UserDataBasePath == "" {
		response, err := RemoteGetUserData(platform, userID, conf)
		if err != nil {
			return nil, err
		}
		return response.Data.DefaultServer, nil
	} else {
		userData, err := GetUserData(platform, userID)
		if err != nil {
			return nil, err
		}
		return userData.DefaultServer, nil
	}
}

func setDefaultServers(platform, userID string, servers []ServerId, conf *config.Config) error {
	if conf.Tsugu.UserDataBasePath == "" {
		update := &TsuguUserUpdate{
			DefaultServer: servers,
		}
		response, err := RemoteChangeUserData(platform, userID, update, conf)
		if err != nil {
			return err
		}
		if response.Status == STATUS_FAILED {
			return fmt.Errorf("%s", response.Data)
		}
		return nil
	} else {
		userData, err := GetUserData(platform, userID)
		if err != nil {
			return err
		}
		userData.DefaultServer = servers
		err = SaveUserData(userData)
		if err != nil {
			return err
		}
		return nil
	}
}

func getMainServerAndDefaultServers(platform, userID string, conf *config.Config) (ServerId, []ServerId, error) {
	if conf.Tsugu.UserDataBasePath == "" {
		response, err := RemoteGetUserData(platform, userID, conf)
		if err != nil {
			return -1, nil, err
		}
		return response.Data.ServerMode, response.Data.DefaultServer, nil
	} else {
		userData, err := GetUserData(platform, userID)
		if err != nil {
			return -1, nil, err
		}
		return userData.ServerMode, userData.DefaultServer, nil
	}
}
