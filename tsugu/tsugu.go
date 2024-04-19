package tsugu

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/WindowsSov8forUs/tsugu-bot-go/adapter"
	"github.com/WindowsSov8forUs/tsugu-bot-go/config"
	log "github.com/WindowsSov8forUs/tsugu-bot-go/mylog"
)

func Handler(session adapter.Session, bot adapter.Bot, conf *config.Config) error {
	message := session.Message()

	// 帮助信息
	if !conf.Tsugu.Functions.Help {
		if strings.HasPrefix(message, "帮助") {
			return helpCommand(session, strings.TrimPrefix(message, "帮助"), bot)
		} else if strings.HasPrefix(message, "help") {
			return helpCommand(session, strings.TrimPrefix(message, "help"), bot)
		} else if strings.HasSuffix(message, "-h") {
			return helpCommand(session, strings.TrimSuffix(message, "-h"), bot)
		}
	}

	// 车牌转发
	if conf.Tsugu.Functions.CarForward {
		forwarded, err := submitCarMessage(session, bot, conf)
		if err != nil {
			log.Errorf("<Tsugu> 车牌转发失败: %v", err)
		}
		if forwarded {
			return nil
		}
	}

	// 进行命令匹配
	command, api := matchCommand(session.Message(), conf)
	if command != "" {
		response, err := ApiV2Command(session.Message(), command, api, session.Platform(), session.UserID(), session.ChannelID(), conf)
		if err != nil {
			log.Errorf("<Tsugu> 执行命令失败: %v", err)
			response = TextResponse("执行命令失败")
		}
		return sendMessage(session, bot, response)
	}

	if conf.Tsugu.UserDataBasePath != "" {
		return localExtraHandler(session, bot, conf)
	} else {
		return remoteExtraHandler(session, bot, conf)
	}
}

func sendMessage(session adapter.Session, bot adapter.Bot, response []*ResponseData) error {
	message := &adapter.Message{}
	for _, data := range response {
		switch data.Type {
		case "string":
			message.Text(data.String)
		case "base64":
			message.Image(data.String)
		}
	}
	return bot.Send(session, message)
}

func localExtraHandler(session adapter.Session, bot adapter.Bot, conf *config.Config) error {
	message := session.Message()
	if conf.Tsugu.Functions.PlayerStatus {
		if strings.HasSuffix(message, "服玩家状态") {
			serverName := strings.Trim(strings.TrimSuffix(message, "玩家状态"), " ")
			server := serverToIngeter(serverName)
			if server == -1 {
				response := TextResponse("未找到被指定的服务器")
				return sendMessage(session, bot, response)
			}
			response, err := playerStatus(session.Platform(), session.UserID(), server, 0, conf)
			if err != nil {
				return err
			}
			return sendMessage(session, bot, response)
		} else if strings.HasPrefix(message, "玩家状态") {
			arg := strings.Trim(strings.TrimPrefix(message, "玩家状态"), " ")
			if arg == "" {
				response, err := playerStatus(session.Platform(), session.UserID(), -1, 0, conf)
				if err != nil {
					return err
				}
				return sendMessage(session, bot, response)
			} else {
				if index, err := strconv.Atoi(arg); err != nil {
					// 服务器名
					server := serverToIngeter(arg)
					response, err := playerStatus(session.Platform(), session.UserID(), server, 0, conf)
					if err != nil {
						return err
					}
					return sendMessage(session, bot, response)
				} else {
					// 序列值
					response, err := playerStatus(session.Platform(), session.UserID(), -1, index, conf)
					if err != nil {
						return err
					}
					return sendMessage(session, bot, response)
				}
			}
		}
	}
	if conf.Tsugu.Functions.BindPlayer {
		if strings.HasPrefix(message, "绑定玩家") {
			serverName := strings.Trim(strings.TrimPrefix(message, "绑定玩家"), " ")
			if server := serverToIngeter(serverName); server == -1 {
				response := TextResponse(
					fmt.Sprintf(
						"未找到名为 %s 的服务器信息，请确保输入的是服务器名而不是玩家ID，通常情况发送“绑定玩家”即可",
						serverName,
					),
				)
				return sendMessage(session, bot, response)
			}
			response := bindPlayer(session.Platform(), session.UserID())
			return sendMessage(session, bot, response)
		} else if strings.HasPrefix(message, "验证") {
			arg := strings.Trim(strings.TrimPrefix(message, "验证"), " ")

			// 正则获取数字
			args := regexp.MustCompile(`\d+`).FindAllString(arg, -1)
			if len(args) == 0 || len(args) > 1 {
				response := TextResponse("请确保输入正确(例如: 验证 10000xxxxx cn)")
				return sendMessage(session, bot, response)
			}
			playerID, _ := strconv.Atoi(args[0])
			serverName := strings.Trim(strings.TrimSuffix(arg, args[0]), " ")
			server := serverToIngeter(serverName)
			response := bindVerify(session.Platform(), session.UserID(), playerID, server, conf)
			return sendMessage(session, bot, response)
		} else if strings.HasPrefix(message, "解除绑定") {
			arg := strings.Trim(strings.TrimPrefix(message, "解除绑定"), " ")
			if arg == "" {
				response := unbindPlayer(session.Platform(), session.UserID())
				return sendMessage(session, bot, response)
			} else {
				args := strings.Split(arg, " ")
				if len(args) < 1 {
					response := unbindPlayer(session.Platform(), session.UserID())
					return sendMessage(session, bot, response)
				} else if index, err := strconv.Atoi(args[len(args)-1]); err != nil {
					response := unbindPlayer(session.Platform(), session.UserID())
					return sendMessage(session, bot, response)
				} else {
					response := unbindVerify(session.Platform(), session.UserID(), index)
					return sendMessage(session, bot, response)
				}
			}
		}
	}
	if conf.Tsugu.Functions.SwitchCarForward {
		if message == "开启车牌转发" || message == "开启个人车牌转发" {
			response := setCarForward(session.Platform(), session.UserID(), true)
			return sendMessage(session, bot, response)
		} else if message == "关闭车牌转发" || message == "关闭个人车牌转发" {
			response := setCarForward(session.Platform(), session.UserID(), false)
			return sendMessage(session, bot, response)
		}
	}
	if conf.Tsugu.Functions.ChangeMainServer {
		if strings.HasSuffix(message, "服模式") {
			serverName := strings.Trim(strings.TrimSuffix(message, "模式"), " ")
			response := setServerMode(session.Platform(), session.UserID(), serverName)
			return sendMessage(session, bot, response)
		} else if strings.HasPrefix(message, "主服务器") {
			serverName := strings.Trim(strings.TrimPrefix(message, "主服务器"), " ")
			response := setServerMode(session.Platform(), session.UserID(), serverName)
			return sendMessage(session, bot, response)
		}
	}
	if conf.Tsugu.Functions.ChangeServerList {
		if strings.HasPrefix(message, "设置默认服务器") {
			serverList := strings.Trim(strings.TrimPrefix(message, "设置默认服务器"), " ")
			response := setDefaultServer(session.Platform(), session.UserID(), serverList)
			return sendMessage(session, bot, response)
		}
	}
	return nil
}

func remoteExtraHandler(session adapter.Session, bot adapter.Bot, conf *config.Config) error {
	message := session.Message()
	if conf.Tsugu.Functions.PlayerStatus {
		if strings.HasSuffix(message, "服玩家状态") {
			serverName := strings.Trim(strings.TrimSuffix(message, "玩家状态"), " ")
			server := serverToIngeter(serverName)
			if server == -1 {
				response := TextResponse("未找到被指定的服务器")
				return sendMessage(session, bot, response)
			}
			response := remotePlayerStatus(session.Platform(), session.UserID(), server, conf)
			return sendMessage(session, bot, response)
		} else if strings.HasPrefix(message, "玩家状态") {
			arg := strings.Trim(strings.TrimPrefix(message, "玩家状态"), " ")
			if arg == "" {
				response := remotePlayerStatus(session.Platform(), session.UserID(), -1, conf)
				return sendMessage(session, bot, response)
			} else {
				// 服务器名
				server := serverToIngeter(arg)
				response := remotePlayerStatus(session.Platform(), session.UserID(), server, conf)
				return sendMessage(session, bot, response)
			}
		}
	}
	if conf.Tsugu.Functions.BindPlayer {
		if strings.HasPrefix(message, "绑定玩家") {
			serverName := strings.Trim(strings.TrimPrefix(message, "绑定玩家"), " ")
			var server int
			if serverName == "" {
				userData, err := remoteGetUserData(session.Platform(), session.UserID(), conf)
				if err != nil {
					response := TextResponse("获取用户信息失败")
					return sendMessage(session, bot, response)
				}
				server = userData.Data.ServerMode
			} else {
				server = serverToIngeter(serverName)
				if server == -1 {
					response := TextResponse(
						fmt.Sprintf(
							"未找到名为 %s 的服务器信息，请确保输入的是服务器名而不是玩家ID，通常情况发送“绑定玩家”即可",
							serverName,
						),
					)
					return sendMessage(session, bot, response)
				}
			}
			res, err := remoteBindPlayer(session.Platform(), session.UserID(), server, true, conf)
			if err != nil {
				return err
			}
			if res.Status != "success" {
				response := TextResponse(fmt.Sprintf("%v", res.Data))
				return sendMessage(session, bot, response)
			}
			data, err := json.Marshal(res.Data)
			if err != nil {
				return err
			}
			var verify = &VerifyData{}
			err = json.Unmarshal(data, verify)
			if err != nil {
				return err
			}
			response := TextResponse(
				fmt.Sprintf(
					"正在绑定账号，请将 评论(个性签名) 或者 当前使用的 乐队编队名称改为\n%v\n稍等片刻等待同步后，发送\n验证 + 空格 + 玩家ID 来完成本次身份验证\n验证 10000xxxx 国服",
					verify.VerifyCode,
				),
			)
			return sendMessage(session, bot, response)
		} else if strings.HasPrefix(message, "验证解绑") {
			arg := strings.Trim(strings.TrimPrefix(message, "验证解绑"), " ")
			if arg == "" {
				response := TextResponse("请输入解绑时提供的serverID(数字)")
				return sendMessage(session, bot, response)
			}
			server := serverToIngeter(arg)
			if server == -1 {
				response := TextResponse("请输入解绑时提供的serverID(数字)")
				return sendMessage(session, bot, response)
			}
			userData, err := remoteGetUserData(session.Platform(), session.UserID(), conf)
			if err != nil {
				response := TextResponse("获取用户信息失败")
				return sendMessage(session, bot, response)
			}
			playerID := userData.Data.GameIDs[server].GameID
			res, err := remoteBindVerify(session.Platform(), session.UserID(), playerID, server, false, conf)
			if err != nil {
				return err
			}
			if res.Status != "success" {
				response := TextResponse(fmt.Sprintf("%v", res.Data))
				return sendMessage(session, bot, response)
			}
			response := TextResponse("解绑成功")
			return sendMessage(session, bot, response)
		} else if strings.HasPrefix(message, "验证") {
			arg := strings.Trim(strings.TrimPrefix(message, "验证"), " ")

			// 正则获取数字
			args := regexp.MustCompile(`\d+`).FindAllString(arg, -1)
			if len(args) == 0 || len(args) > 1 {
				response := TextResponse("请确保输入正确(例如: 验证 10000xxxxx cn)")
				return sendMessage(session, bot, response)
			}
			playerID, _ := strconv.Atoi(args[0])
			serverName := strings.Trim(strings.TrimPrefix(arg, args[0]), " ")
			server := serverToIngeter(serverName)
			if server == -1 {
				userData, err := remoteGetUserData(session.Platform(), session.UserID(), conf)
				if err != nil {
					response := TextResponse("获取用户信息失败")
					return sendMessage(session, bot, response)
				}
				server = userData.Data.ServerMode
			}
			res, err := remoteBindVerify(session.Platform(), session.UserID(), playerID, server, true, conf)
			if err != nil {
				return err
			}
			if res.Status != "success" {
				var response []*ResponseData
				if res.Data != "" {
					response = TextResponse(fmt.Sprintf("%v", res.Data))
				} else {
					response = TextResponse("错误: 未请求绑定或解除绑定玩家")
				}
				return sendMessage(session, bot, response)
			}
			response := TextResponse("绑定成功！现在可以使用“玩家状态”来查询玩家信息")
			return sendMessage(session, bot, response)
		} else if strings.HasPrefix(message, "解除绑定") {
			arg := strings.Trim(strings.TrimPrefix(message, "解除绑定"), " ")
			var server int
			if arg == "" {
				userData, err := remoteGetUserData(session.Platform(), session.UserID(), conf)
				if err != nil {
					response := TextResponse("获取用户信息失败")
					return sendMessage(session, bot, response)
				}
				server = userData.Data.ServerMode
			} else {
				server = serverToIngeter(arg)
				if server == -1 {
					response := TextResponse(fmt.Sprintf("未找到名为 %s 的服务器信息，请确保输入的是服务器名而不是玩家ID", arg))
					return sendMessage(session, bot, response)
				}
			}
			res, err := remoteBindPlayer(session.Platform(), session.UserID(), server, false, conf)
			if err != nil {
				return err
			}
			if res.Status != "success" {
				response := TextResponse(fmt.Sprintf("%v", res.Data))
				return sendMessage(session, bot, response)
			}
			data, err := json.Marshal(res.Data)
			if err != nil {
				return err
			}
			var verify = &VerifyData{}
			err = json.Unmarshal(data, verify)
			if err != nil {
				return err
			}
			response := TextResponse(
				fmt.Sprintf(
					"正在解除绑定账号，请将 评论(个性签名) 或者 当前使用的 乐队编队名称改为\n%v\n稍等片刻等待同步后，发送\n验证解绑 %v 来完成本次身份验证",
					verify.VerifyCode,
					server,
				),
			)
			return sendMessage(session, bot, response)
		}
	}
	if conf.Tsugu.Functions.SwitchCarForward {
		if message == "开启车牌转发" || message == "开启个人车牌转发" {
			res, err := remoteSetCarForward(session.Platform(), session.UserID(), true, conf)
			if err != nil {
				return err
			}
			if res.Status == "success" {
				response := TextResponse("已开启车牌转发")
				return sendMessage(session, bot, response)
			} else {
				return nil
			}
		} else if message == "关闭车牌转发" || message == "关闭个人车牌转发" {
			res, err := remoteSetCarForward(session.Platform(), session.UserID(), false, conf)
			if err != nil {
				return err
			}
			if res.Status == "success" {
				response := TextResponse("已关闭车牌转发")
				return sendMessage(session, bot, response)
			} else {
				return nil
			}
		}
	}
	if conf.Tsugu.Functions.ChangeMainServer {
		if strings.HasSuffix(message, "服模式") {
			serverName := strings.Trim(strings.TrimSuffix(message, "模式"), " ")
			res, err := remoteSetServerMode(session.Platform(), session.UserID(), serverName, conf)
			if err != nil {
				return err
			}
			if res.Status != "success" {
				response := TextResponse(fmt.Sprintf("%v", res.Data))
				return sendMessage(session, bot, response)
			}
			response := TextResponse(fmt.Sprintf("已切换至 %v 模式", serverName))
			return sendMessage(session, bot, response)
		} else if strings.HasPrefix(message, "主服务器") {
			serverName := strings.Trim(strings.TrimPrefix(message, "主服务器"), " ")
			res, err := remoteSetServerMode(session.Platform(), session.UserID(), serverName, conf)
			if err != nil {
				return err
			}
			if res.Status != "success" {
				response := TextResponse(fmt.Sprintf("%v", res.Data))
				return sendMessage(session, bot, response)
			}
			response := TextResponse(fmt.Sprintf("主服务器已切换为 %v", serverName))
			return sendMessage(session, bot, response)
		}
	}
	if conf.Tsugu.Functions.ChangeServerList {
		if strings.HasPrefix(message, "设置默认服务器") {
			serverList := strings.Trim(strings.TrimPrefix(message, "设置默认服务器"), " ")
			res, err := remoteSetDefaultServer(session.Platform(), session.UserID(), serverList, conf)
			if err != nil {
				return err
			}
			if res.Status != "success" {
				response := TextResponse(fmt.Sprintf("%v", res.Data))
				return sendMessage(session, bot, response)
			}
			response := TextResponse(fmt.Sprintf("默认服务器已设置为 %v", serverList))
			return sendMessage(session, bot, response)
		}
	}
	return nil
}
