package tsugu

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/WindowsSov8forUs/tsugu-bot-go/adapter"
	"github.com/WindowsSov8forUs/tsugu-bot-go/config"
)

type Api string

const (
	ApiEventStage                 Api = "/eventStage"
	ApiGachaSimulate              Api = "/gachaSimulate"
	ApiGetCardIllustration        Api = "/getCardIllustration"
	ApiLsycx                      Api = "/lsycx"
	ApiRoomList                   Api = "/roomList"
	ApiSearchCard                 Api = "/searchCard"
	ApiSearchCharacter            Api = "/searchCharacter"
	ApiSearchEvent                Api = "/searchEvent"
	ApiSearchGacha                Api = "/searchGacha"
	ApiSearchPlayer               Api = "/searchPlayer"
	ApiSearchSong                 Api = "/searchSong"
	ApiSongChart                  Api = "/songChart"
	ApiSongMeta                   Api = "/songMeta"
	ApiYcx                        Api = "/ycx"
	ApiYcxAll                     Api = "/ycxAll"
	ApiUserGetUserData            Api = "/user/getUserData"
	ApiUserChangeUserData         Api = "/user/changeUserData"
	ApiUserBindPlayerRequest      Api = "/user/bindPlayerRequest"
	ApiUserBindPlayerVerification Api = "/user/bindPlayerVerification"
	ApiStationSubmitRoomNumber    Api = "/station/submitRoomNumber"
	ApiStationQueryAllRoom        Api = "/station/queryAllRoom"
)

type BindingData struct {
	Platform string
	UserID   string
	BindType bool
}

type BindingRecord struct {
	mutex *sync.Mutex
	datas []*BindingData
}

func NewBindingRecord() *BindingRecord {
	return &BindingRecord{
		mutex: &sync.Mutex{},
		datas: make([]*BindingData, 0),
	}
}

func (r *BindingRecord) Add(platform string, userID string, bindType bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.datas = append(r.datas, &BindingData{
		Platform: platform,
		UserID:   userID,
		BindType: bindType,
	})
}

func (r *BindingRecord) Remove(platform string, userID string, bindType bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for i, data := range r.datas {
		if data.Platform == platform && data.UserID == userID && data.BindType == bindType {
			r.datas = append(r.datas[:i], r.datas[i+1:]...)
			break
		}
	}
}

func (r *BindingRecord) IsExist(platform string, userID string, bindType bool) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, data := range r.datas {
		if data.Platform == platform && data.UserID == userID && data.BindType == bindType {
			return true
		}
	}
	return false
}

var bindingRecord = NewBindingRecord()

func matchCommand(message string, commands []string, conf config.Config) (string, bool) {
	for _, command := range commands {
		if strings.HasPrefix(message, command) {
			if len(message) == len(command) {
				return "", true
			} else {
				if conf.Tsugu.NoSpace {
					return strings.TrimSpace(message[len(command):]), true
				} else {
					if message[len(command)] == ' ' {
						return strings.TrimSpace(message[len(command)+1:]), true
					}
				}
			}
		}
	}
	return message, false
}

func switchGachaSimulation(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, error) {
	if !conf.Tsugu.Functions.SwitchGachaSimulate {
		return false, nil
	}

	message := strings.TrimSpace(session.Message())
	var ifSwitch = false
	if message == "开启抽卡" {
		ifSwitch = true
	} else if message == "关闭抽卡" {
		ifSwitch = false
	} else {
		if arg, ok := matchCommand(message, append([]string{"抽卡"}, conf.Tsugu.CommandAlias.SwitchGachaSimulate...), *conf); ok {
			if arg == "开启" || arg == "on" {
				ifSwitch = true
			} else if arg == "关闭" || arg == "off" {
				ifSwitch = false
			} else {
				reply := TextResponse("无效指令")
				return false, sendMessage(session, bot, reply)
			}
		} else {
			return false, nil
		}
	}

	if ifSwitch {
		channelId := session.ChannelID()
		conf.Tsugu.RemoveBanList(channelId)
		reply := TextResponse("开启成功")
		return true, sendMessage(session, bot, reply)
	} else {
		channelId := session.ChannelID()
		conf.Tsugu.AddBanList(channelId)
		reply := TextResponse("关闭成功")
		return true, sendMessage(session, bot, reply)
	}
}

func openCarForward(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, error) {
	if !conf.Tsugu.Functions.SwitchCarForward {
		return false, nil
	}

	message := strings.TrimSpace(session.Message())
	for _, command := range append([]string{"开启车牌转发"}, conf.Tsugu.CommandAlias.OpenCarForward...) {
		if message == command {
			result, err := setCarForward(session.Platform(), session.UserID(), true, conf)
			if err != nil {
				if result != "" {
					reply := TextResponse(result)
					return true, sendMessage(session, bot, reply)
				} else {
					reply := TextResponse(fmt.Sprintf("错误：%s", err.Error()))
					return true, sendMessage(session, bot, reply)
				}
			}
			reply := TextResponse("已开启车牌转发")
			return true, sendMessage(session, bot, reply)
		}
	}
	return false, nil
}

func closeCarForward(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, error) {
	if !conf.Tsugu.Functions.SwitchCarForward {
		return false, nil
	}

	message := strings.TrimSpace(session.Message())
	for _, command := range append([]string{"关闭车牌转发"}, conf.Tsugu.CommandAlias.CloseCarForward...) {
		if message == command {
			result, err := setCarForward(session.Platform(), session.UserID(), false, conf)
			if err != nil {
				if result != "" {
					reply := TextResponse(result)
					return true, sendMessage(session, bot, reply)
				} else {
					reply := TextResponse(fmt.Sprintf("错误：%s", err.Error()))
					return true, sendMessage(session, bot, reply)
				}
			}
			reply := TextResponse("已关闭车牌转发")
			return true, sendMessage(session, bot, reply)
		}
	}
	return false, nil
}

func bindPlayerHandler(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, error) {
	if !conf.Tsugu.Functions.BindPlayer {
		return false, nil
	}

	message := strings.TrimSpace(session.Message())
	if arg, ok := matchCommand(message, append([]string{"绑定玩家"}, conf.Tsugu.CommandAlias.BindPlayer...), *conf); ok {
		var server ServerId
		var err error

		if arg != "" {
			server = serverNameToId(arg)
			if server == -1 {
				reply := TextResponse(fmt.Sprintf("无效的服务器：%s", arg))
				return true, sendMessage(session, bot, reply)
			}
		} else {
			server, err = getMainServer(session.Platform(), session.UserID(), conf)
			if err != nil {
				reply := TextResponse(fmt.Sprintf("获取用户数据出错：%s", err.Error()))
				return true, sendMessage(session, bot, reply)
			}
		}

		if conf.Tsugu.UserDataBasePath == "" {
			response, err := RemoteBindPlayerRequest(session.Platform(), session.UserID(), server, true, conf)
			if err != nil {
				if e, ok := err.(*UserApiFailedError); ok {
					reply := TextResponse(e.Error())
					return true, sendMessage(session, bot, reply)
				} else {
					reply := TextResponse(fmt.Sprintf("错误：%s", err.Error()))
					return true, sendMessage(session, bot, reply)
				}
			}
			reply := TextResponse(
				fmt.Sprintf(
					"正在绑定 %s 账号，请将你的\n评论(个性签名)\n或者\n你的当前使用的卡组的卡组名(乐队编队名称)\n改为以下数字后，直接发送你的玩家id\n%v",
					serverIdToFullName(server),
					response.Data.VerifyCode,
				),
			)
			registerTempHandler(bindVerifyHandler(server))
			bindingRecord.Add(session.Platform(), session.UserID(), true)
			return true, sendMessage(session, bot, reply)
		} else {
			response := bindPlayer(session.Platform(), session.UserID(), server)
			registerTempHandler(bindVerifyHandler(server))
			bindingRecord.Add(session.Platform(), session.UserID(), true)
			return true, sendMessage(session, bot, response)
		}

	} else {
		return false, nil
	}
}

func bindVerifier(session adapter.Session, bot adapter.Bot, server ServerId, conf *config.Config) (bool, bool, error) {
	if !conf.Tsugu.Functions.BindPlayer {
		return false, false, nil
	}

	if !bindingRecord.IsExist(session.Platform(), session.UserID(), true) {
		return false, false, nil
	}

	message := strings.TrimSpace(session.Message())
	if message == "" {
		return false, true, nil
	}

	playerID, err := strconv.Atoi(message)
	if err != nil {
		return true, true, sendMessage(session, bot, TextResponse("错误：无效的玩家id"))
	}

	if conf.Tsugu.UserDataBasePath == "" {
		response, err := RemoteBindVerification(session.Platform(), session.UserID(), server, playerID, true, conf)
		if err != nil {
			reply := TextResponse(fmt.Sprintf("错误：%s", err.Error()))
			return true, true, sendMessage(session, bot, reply)
		}
		bindingRecord.Remove(session.Platform(), session.UserID(), true)
		return true, false, sendMessage(session, bot, TextResponse(response.Data))
	} else {
		if verifyCodeMap[fmt.Sprintf("%s:%s", session.Platform(), session.UserID())] == nil {
			bindingRecord.Remove(session.Platform(), session.UserID(), true)
			return false, false, nil
		}

		response, retryable := bindVerify(session.Platform(), session.UserID(), playerID, conf)
		if retryable {
			return true, true, sendMessage(session, bot, response)
		}
		bindingRecord.Remove(session.Platform(), session.UserID(), true)
		return true, false, sendMessage(session, bot, response)
	}
}

func bindVerifyHandler(server ServerId) tsuguTempHandler {
	return func(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, bool, error) {
		return bindVerifier(session, bot, server, conf)
	}
}

func unbindPlayerHandler(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, error) {
	if !conf.Tsugu.Functions.BindPlayer {
		return false, nil
	}

	message := strings.TrimSpace(session.Message())
	if arg, ok := matchCommand(message, append([]string{"解除绑定"}, conf.Tsugu.CommandAlias.UnbindPlayer...), *conf); ok {

		if conf.Tsugu.UserDataBasePath == "" {
			var server ServerId
			var err error

			if arg != "" {
				server = serverNameToId(arg)
				if server == -1 {
					reply := TextResponse(fmt.Sprintf("无效的服务器：%s", arg))
					return true, sendMessage(session, bot, reply)
				}
			} else {
				server, err = getMainServer(session.Platform(), session.UserID(), conf)
				if err != nil {
					reply := TextResponse(fmt.Sprintf("获取用户数据出错：%s", err.Error()))
					return true, sendMessage(session, bot, reply)
				}
			}

			response, err := RemoteBindPlayerRequest(session.Platform(), session.UserID(), server, false, conf)
			if err != nil {
				if e, ok := err.(*UserApiFailedError); ok {
					reply := TextResponse(e.Error())
					return true, sendMessage(session, bot, reply)
				} else {
					reply := TextResponse(fmt.Sprintf("错误：%s", err.Error()))
					return true, sendMessage(session, bot, reply)
				}
			}
			reply := TextResponse(
				fmt.Sprintf(
					"正在解除绑定 %s 账号\n因为使用远程服务器，解除绑定需要验证\n请将账号的\n评论(个性签名)\n或者\n你的当前使用的卡组的卡组名(乐队编队名称)\n改为以下数字后，发送玩家id继续\n%v",
					serverIdToFullName(server),
					response.Data.VerifyCode,
				),
			)
			registerTempHandler(unbindVerifyHandler(server))
			bindingRecord.Add(session.Platform(), session.UserID(), false)
			return true, sendMessage(session, bot, reply)
		} else {
			var index int
			var err error

			if arg != "" {
				index, err = strconv.Atoi(arg)
				if err != nil {
					reply := TextResponse("无效的序号值")
					return true, sendMessage(session, bot, reply)
				}
			} else {
				index = 1
			}
			response := unbindPlayer(session.Platform(), session.UserID(), index)
			return true, sendMessage(session, bot, response)
		}

	} else {
		return false, nil
	}
}

func unbindVerifier(session adapter.Session, bot adapter.Bot, server ServerId, conf *config.Config) (bool, bool, error) {
	if !conf.Tsugu.Functions.BindPlayer {
		return false, false, nil
	}

	if !bindingRecord.IsExist(session.Platform(), session.UserID(), false) {
		return false, false, nil
	}

	message := strings.TrimSpace(session.Message())
	if message == "" {
		return false, true, nil
	}

	playerID, err := strconv.Atoi(message)
	if err != nil {
		return true, true, sendMessage(session, bot, TextResponse("错误：无效的玩家id"))
	}

	if conf.Tsugu.UserDataBasePath == "" {
		response, err := RemoteBindVerification(session.Platform(), session.UserID(), server, playerID, false, conf)
		if err != nil {
			reply := TextResponse(fmt.Sprintf("错误：%s", err.Error()))
			return true, true, sendMessage(session, bot, reply)
		}
		bindingRecord.Remove(session.Platform(), session.UserID(), false)
		return true, false, sendMessage(session, bot, TextResponse(response.Data))
	} else {
		bindingRecord.Remove(session.Platform(), session.UserID(), false)
		return false, false, nil
	}
}

func unbindVerifyHandler(server ServerId) tsuguTempHandler {
	return func(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, bool, error) {
		return unbindVerifier(session, bot, server, conf)
	}
}

func mainServerHandler(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, error) {
	if !conf.Tsugu.Functions.ChangeMainServer {
		return false, nil
	}

	message := strings.TrimSpace(session.Message())
	var server ServerId

	if arg, ok := matchCommand(message, append([]string{"主服务器"}, conf.Tsugu.CommandAlias.ChangeMainServer...), *conf); ok {
		if arg == "" {
			return true, sendMessage(session, bot, TextResponse("错误: 未指定服务器"))
		}
		server = serverNameToId(arg)
		if server == -1 {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("错误: 服务器不存在：%s", arg)))
		}
	} else {
		pattern := `^(日服|国际服|台服|国服|韩服)模式$`
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(message)
		if matches != nil {
			server = serverNameToId(matches[1])
			if server == -1 {
				return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("错误: 服务器不存在：%s", matches[1])))
			}
		} else {
			return false, nil
		}
	}

	err := setMainServer(session.Platform(), session.UserID(), server, conf)
	if err != nil {
		return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("错误：%s", err.Error())))
	}
	return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("已切换到%s模式", serverIdToFullName(server))))
}

func defaultServersHandler(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, error) {
	if !conf.Tsugu.Functions.ChangeServerList {
		return false, nil
	}

	message := strings.TrimSpace(session.Message())
	if arg, ok := matchCommand(message, append([]string{"设置默认服务器"}, conf.Tsugu.CommandAlias.ChangeServerList...), *conf); ok {
		args := strings.Split(arg, " ")
		var servers []ServerId
		for _, arg := range args {
			server := serverNameToId(arg)
			if server == -1 {
				return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("错误: 指定了不存在的服务器：%s", arg)))
			} else {
				for _, s := range servers {
					if s == server {
						return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("错误: 指定了重复的服务器：%s", arg)))
					}
				}
			}
			servers = append(servers, server)
		}
		if len(servers) == 0 {
			return true, sendMessage(session, bot, TextResponse("错误: 请指定至少一个服务器"))
		}
		err := setDefaultServers(session.Platform(), session.UserID(), servers, conf)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("错误：%s", err.Error())))
		}
		var serverNames string
		for i, server := range servers {
			if i != 0 {
				serverNames += ", "
			}
			serverNames += serverIdToFullName(server)
		}
		return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("成功切换默认服务器顺序：%s", serverNames)))
	} else {
		return false, nil
	}
}

func ycm(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, error) {
	if !conf.Tsugu.Functions.SwitchGachaSimulate {
		return false, nil
	}

	message := strings.TrimSpace(session.Message())
	if _, ok := matchCommand(message, append([]string{"ycm"}, conf.Tsugu.CommandAlias.SwitchGachaSimulate...), *conf); ok {
		queryResponse, err := queryAllRoom(conf)
		if err != nil {
			return true, err
		}
		request := &RequestRoomList{
			RoomList: queryResponse.Data,
			Compress: conf.Tsugu.Compress,
		}
		response, err := requestBackend(ApiRoomList, request, conf)
		if err != nil {
			return true, err
		}
		return true, sendMessage(session, bot, response)
	} else {
		return false, nil
	}
}

func searchPlayer(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, error) {
	if !conf.Tsugu.Functions.SearchPlayer {
		return false, nil
	}

	message := strings.TrimSpace(session.Message())
	if arg, ok := matchCommand(message, append([]string{"查玩家"}, conf.Tsugu.CommandAlias.SearchPlayer...), *conf); ok {
		args := strings.Split(arg, " ")
		if args[0] == "" {
			return true, sendMessage(session, bot, TextResponse("查询指定ID玩家的信息。省略服务器名时，默认从你当前的主服务器查询"))
		}
		playerID, err := strconv.Atoi(args[0])
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("无效的玩家ID: %s", args[0])))
		}
		request := &RequestSearchPlayer{
			PlayerID:  playerID,
			UseEasyBG: conf.Tsugu.UseEasyBG,
			Compress:  conf.Tsugu.Compress,
		}
		if len(args) > 1 {
			server := serverNameToId(args[1])
			if server == -1 {
				return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("无效的服务器：%s", args[1])))
			}
			request.Server = server
		} else {
			server, err := getMainServer(session.Platform(), session.UserID(), conf)
			if err != nil {
				return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("获取用户数据错误：%s", err.Error())))
			}
			request.Server = server
		}
		response, err := requestBackend(ApiSearchPlayer, request, conf)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("错误：%s", err.Error())))
		}
		return true, sendMessage(session, bot, response)
	} else {
		return false, nil
	}
}

func searchCard(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, error) {
	if !conf.Tsugu.Functions.SearchCard {
		return false, nil
	}

	message := strings.TrimSpace(session.Message())
	if arg, ok := matchCommand(message, append([]string{"查卡"}, conf.Tsugu.CommandAlias.SearchCard...), *conf); ok {
		defaultServers, err := getDefaultServers(session.Platform(), session.UserID(), conf)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("获取用户数据错误：%s", err.Error())))
		}
		request := &RequestSearchCard{
			DefaultServers: defaultServers,
			Text:           arg,
			UseEasyBG:      conf.Tsugu.UseEasyBG,
			Compress:       conf.Tsugu.Compress,
		}
		response, err := requestBackend(ApiSearchCard, request, conf)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("错误：%s", err.Error())))
		}
		return true, sendMessage(session, bot, response)
	} else {
		return false, nil
	}
}

func cardIllustration(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, error) {
	if !conf.Tsugu.Functions.CardIllustration {
		return false, nil
	}

	message := strings.TrimSpace(session.Message())
	if arg, ok := matchCommand(message, append([]string{"查卡面"}, conf.Tsugu.CommandAlias.CardIllustration...), *conf); ok {
		if arg == "" {
			return true, sendMessage(session, bot, TextResponse("根据卡片ID查询卡片插画"))
		}
		cardID, err := strconv.Atoi(arg)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("无效的卡片ID: %s", arg)))
		}
		request := &RequestGetCardIllustration{
			CardID: cardID,
		}
		response, err := requestBackend(ApiGetCardIllustration, request, conf)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("错误：%s", err.Error())))
		}
		return true, sendMessage(session, bot, response)
	} else {
		return false, nil
	}
}

func searchCharacter(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, error) {
	if !conf.Tsugu.Functions.SearchCharacter {
		return false, nil
	}

	message := strings.TrimSpace(session.Message())
	if arg, ok := matchCommand(message, append([]string{"查角色"}, conf.Tsugu.CommandAlias.SearchCharacter...), *conf); ok {
		defaultServers, err := getDefaultServers(session.Platform(), session.UserID(), conf)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("获取用户数据错误：%s", err.Error())))
		}
		request := &RequestSearchCharacter{
			DefaultServers: defaultServers,
			Text:           arg,
			Compress:       conf.Tsugu.Compress,
		}
		response, err := requestBackend(ApiSearchCharacter, request, conf)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("错误：%s", err.Error())))
		}
		return true, sendMessage(session, bot, response)
	} else {
		return false, nil
	}
}

func searchEvent(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, error) {
	if !conf.Tsugu.Functions.SearchEvent {
		return false, nil
	}

	message := strings.TrimSpace(session.Message())
	if arg, ok := matchCommand(message, append([]string{"查活动"}, conf.Tsugu.CommandAlias.SearchEvent...), *conf); ok {
		defaultServers, err := getDefaultServers(session.Platform(), session.UserID(), conf)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("获取用户数据错误：%s", err.Error())))
		}
		request := &RequestSearchEvent{
			DefaultServers: defaultServers,
			Text:           arg,
			UseEasyBG:      conf.Tsugu.UseEasyBG,
			Compress:       conf.Tsugu.Compress,
		}
		response, err := requestBackend(ApiSearchEvent, request, conf)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("错误：%s", err.Error())))
		}
		return true, sendMessage(session, bot, response)
	} else {
		return false, nil
	}
}

func searchSong(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, error) {
	if !conf.Tsugu.Functions.SearchSong {
		return false, nil
	}

	message := strings.TrimSpace(session.Message())
	if arg, ok := matchCommand(message, append([]string{"查曲"}, conf.Tsugu.CommandAlias.SearchSong...), *conf); ok {
		defaultServers, err := getDefaultServers(session.Platform(), session.UserID(), conf)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("获取用户数据错误：%s", err.Error())))
		}
		request := &RequestSearchSong{
			DefaultServers: defaultServers,
			Text:           arg,
			Compress:       conf.Tsugu.Compress,
		}
		response, err := requestBackend(ApiSearchSong, request, conf)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("错误：%s", err.Error())))
		}
		return true, sendMessage(session, bot, response)
	} else {
		return false, nil
	}
}

func songChart(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, error) {
	if !conf.Tsugu.Functions.SearchChart {
		return false, nil
	}

	message := strings.TrimSpace(session.Message())
	if arg, ok := matchCommand(message, append([]string{"查谱面"}, conf.Tsugu.CommandAlias.SearchChart...), *conf); ok {
		args := strings.Split(arg, " ")
		if args[0] == "" {
			return true, sendMessage(session, bot, TextResponse("根据曲目ID查询谱面信息"))
		}
		songID, err := strconv.Atoi(args[0])
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("无效的曲目ID: %s", args[0])))
		}
		defaultServers, err := getDefaultServers(session.Platform(), session.UserID(), conf)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("获取用户数据错误：%s", err.Error())))
		}
		request := &RequestSongChart{
			DefaultServers: defaultServers,
			SongID:         songID,
			Compress:       conf.Tsugu.Compress,
		}
		if len(args) > 1 {
			request.DifficultyText = args[1]
		} else {
			request.DifficultyText = "expert"
		}
		response, err := requestBackend(ApiSongChart, request, conf)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("错误：%s", err.Error())))
		}
		return true, sendMessage(session, bot, response)
	} else {
		return false, nil
	}
}

func songMeta(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, error) {
	if !conf.Tsugu.Functions.SongMeta {
		return false, nil
	}

	message := strings.TrimSpace(session.Message())
	if arg, ok := matchCommand(message, append([]string{"查询分数表"}, conf.Tsugu.CommandAlias.SongMeta...), *conf); ok {
		server, defaultServers, err := getMainServerAndDefaultServers(session.Platform(), session.UserID(), conf)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("获取用户数据错误：%s", err.Error())))
		}
		if arg != "" {
			server = serverNameToId(arg)
			if server == -1 {
				return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("无效的服务器：%s", arg)))
			}
		}
		request := &RequestSongMeta{
			DefaultServers: defaultServers,
			Server:         server,
			Compress:       conf.Tsugu.Compress,
		}
		response, err := requestBackend(ApiSongMeta, request, conf)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("错误：%s", err.Error())))
		}
		return true, sendMessage(session, bot, response)
	} else {
		return false, nil
	}
}

func eventStage(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, error) {
	if !conf.Tsugu.Functions.EventStage {
		return false, nil
	}

	message := strings.TrimSpace(session.Message())
	if arg, ok := matchCommand(message, append([]string{"查试炼"}, conf.Tsugu.CommandAlias.EventStage...), *conf); ok {
		args := strings.Split(arg, " ")
		server, err := getMainServer(session.Platform(), session.UserID(), conf)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("获取用户数据错误：%s", err.Error())))
		}
		request := &RequestEventStage{
			Server:   server,
			Compress: conf.Tsugu.Compress,
		}
		if args[0] != "" {
			eventID, err := strconv.Atoi(args[0])
			if err != nil {
				return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("无效的活动ID: %s", args[0])))
			}
			request.EventID = eventID
		}
		if len(args) > 1 {
			if args[1] == "-m" {
				request.Meta = true
			}
		}
		response, err := requestBackend(ApiEventStage, request, conf)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("错误：%s", err.Error())))
		}
		return true, sendMessage(session, bot, response)
	} else {
		return false, nil
	}
}

func searchGacha(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, error) {
	if !conf.Tsugu.Functions.SearchGacha {
		return false, nil
	}

	message := strings.TrimSpace(session.Message())
	if arg, ok := matchCommand(message, append([]string{"查卡池"}, conf.Tsugu.CommandAlias.SearchGacha...), *conf); ok {
		defaultServers, err := getDefaultServers(session.Platform(), session.UserID(), conf)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("获取用户数据错误：%s", err.Error())))
		}
		request := &RequestSearchGacha{
			DefaultServers: defaultServers,
			UseEasyBG:      conf.Tsugu.UseEasyBG,
			Compress:       conf.Tsugu.Compress,
		}
		if arg != "" {
			gachaID, err := strconv.Atoi(arg)
			if err != nil {
				return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("无效的卡池ID: %s", arg)))
			}
			request.GachaID = gachaID
		}
		response, err := requestBackend(ApiSearchGacha, request, conf)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("错误：%s", err.Error())))
		}
		return true, sendMessage(session, bot, response)
	} else {
		return false, nil
	}
}

func ycx(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, error) {
	if !conf.Tsugu.Functions.Ycx {
		return false, nil
	}

	message := strings.TrimSpace(session.Message())
	if arg, ok := matchCommand(message, append([]string{"ycx"}, conf.Tsugu.CommandAlias.Ycx...), *conf); ok {
		args := strings.Split(arg, " ")
		if args[0] == "" {
			return true, sendMessage(session, bot, TextResponse("查询指定档位的预测线，如果没有服务器名的话，服务器为用户的默认服务器。如果没有活动ID的话，活动为当前活动"))
		}
		tier, err := strconv.Atoi(args[0])
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("无效的档位: %s", args[0])))
		}
		request := &RequestYcx{
			Tier:     tier,
			Compress: conf.Tsugu.Compress,
		}
		if len(args) > 1 {
			eventID, err := strconv.Atoi(args[1])
			if err != nil {
				return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("无效的活动ID: %s", args[1])))
			}
			request.EventID = eventID
		}
		if len(args) > 2 {
			server := serverNameToId(args[2])
			if server == -1 {
				return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("无效的服务器：%s", args[2])))
			}
			request.Server = server
		} else {
			server, err := getMainServer(session.Platform(), session.UserID(), conf)
			if err != nil {
				return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("获取用户数据错误：%s", err.Error())))
			}
			request.Server = server
		}
		response, err := requestBackend(ApiYcx, request, conf)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("错误：%s", err.Error())))
		}
		return true, sendMessage(session, bot, response)
	} else {
		return false, nil
	}
}

func ycxAll(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, error) {
	if !conf.Tsugu.Functions.YcxAll {
		return false, nil
	}

	message := strings.TrimSpace(session.Message())
	if arg, ok := matchCommand(message, append([]string{"ycxall"}, conf.Tsugu.CommandAlias.YcxAll...), *conf); ok {
		args := strings.Split(arg, " ")
		request := &RequestYcxAll{
			Compress: conf.Tsugu.Compress,
		}
		if args[0] != "" {
			eventID, err := strconv.Atoi(args[0])
			if err != nil {
				return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("无效活动ID: %s", args[0])))
			}
			request.EventID = eventID
		}
		if len(args) > 1 {
			server := serverNameToId(args[1])
			if server == -1 {
				return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("无效的服务器：%s", args[1])))
			}
			request.Server = server
		} else {
			server, err := getMainServer(session.Platform(), session.UserID(), conf)
			if err != nil {
				return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("获取用户数据错误：%s", err.Error())))
			}
			request.Server = server
		}
		response, err := requestBackend(ApiYcxAll, request, conf)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("错误：%s", err.Error())))
		}
		return true, sendMessage(session, bot, response)
	} else {
		return false, nil
	}
}

func lsycx(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, error) {
	if !conf.Tsugu.Functions.Lsycx {
		return false, nil
	}

	message := strings.TrimSpace(session.Message())
	if arg, ok := matchCommand(message, append([]string{"lsycx"}, conf.Tsugu.CommandAlias.Lsycx...), *conf); ok {
		args := strings.Split(arg, " ")
		if args[0] == "" {
			return true, sendMessage(session, bot, TextResponse("查询指定档位的预测线，与最近的4期活动类型相同的活动的档线数据，如果没有服务器名的话，服务器为用户的默认服务器。如果没有活动ID的话，活动为当前活动"))
		}
		tier, err := strconv.Atoi(args[0])
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("无效的档位: %s", args[0])))
		}
		request := &RequestLsycx{
			Tier:     tier,
			Compress: conf.Tsugu.Compress,
		}
		if len(args) > 1 {
			eventID, err := strconv.Atoi(args[1])
			if err != nil {
				return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("无效的活动ID: %s", args[1])))
			}
			request.EventID = eventID
		}
		if len(args) > 2 {
			server := serverNameToId(args[2])
			if server == -1 {
				return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("无效的服务器：%s", args[2])))
			}
			request.Server = server
		} else {
			server, err := getMainServer(session.Platform(), session.UserID(), conf)
			if err != nil {
				return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("获取用户数据错误：%s", err.Error())))
			}
			request.Server = server
		}
		response, err := requestBackend(ApiLsycx, request, conf)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("错误：%s", err.Error())))
		}
		return true, sendMessage(session, bot, response)
	} else {
		return false, nil
	}
}

func gachaSimulate(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, error) {
	if !conf.Tsugu.Functions.GachaSimulate {
		return false, nil
	}
	for _, channelId := range conf.Tsugu.BanGachaSimulate {
		if channelId == session.ChannelID() {
			return true, sendMessage(session, bot, TextResponse("抽卡功能已关闭"))
		}
	}

	message := strings.TrimSpace(session.Message())
	if arg, ok := matchCommand(message, append([]string{"抽卡模拟"}, conf.Tsugu.CommandAlias.GachaSimulate...), *conf); ok {
		args := strings.Split(arg, " ")
		server, err := getMainServer(session.Platform(), session.UserID(), conf)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("获取用户数据错误：%s", err.Error())))
		}
		request := &RequestGachaSimulate{
			ServerMode: server,
			Compress:   conf.Tsugu.Compress,
		}
		if args[0] != "" {
			times, err := strconv.Atoi(args[0])
			if err != nil {
				return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("无效的抽卡次数: %s", args[0])))
			}
			request.Times = times
		}
		if len(args) > 1 {
			gachaID, err := strconv.Atoi(args[1])
			if err != nil {
				return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("无效的卡池ID: %s", args[1])))
			}
			request.GachaID = gachaID
		}
		response, err := requestBackend(ApiGachaSimulate, request, conf)
		if err != nil {
			return true, sendMessage(session, bot, TextResponse(fmt.Sprintf("错误：%s", err.Error())))
		}
		return true, sendMessage(session, bot, response)
	} else {
		return false, nil
	}
}

func init() {
	registerHandler(switchGachaSimulation)
	registerHandler(openCarForward)
	registerHandler(closeCarForward)
	registerHandler(bindPlayerHandler)
	registerHandler(unbindPlayerHandler)
	registerHandler(mainServerHandler)
	registerHandler(defaultServersHandler)
	registerHandler(ycm)
	registerHandler(searchPlayer)
	registerHandler(searchCard)
	registerHandler(cardIllustration)
	registerHandler(searchCharacter)
	registerHandler(searchEvent)
	registerHandler(searchSong)
	registerHandler(songChart)
	registerHandler(songMeta)
	registerHandler(eventStage)
	registerHandler(searchGacha)
	registerHandler(ycx)
	registerHandler(ycxAll)
	registerHandler(lsycx)
	registerHandler(gachaSimulate)
}
