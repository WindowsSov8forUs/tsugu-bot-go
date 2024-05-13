package tsugu

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/WindowsSov8forUs/tsugu-bot-go/adapter"
	"github.com/WindowsSov8forUs/tsugu-bot-go/config"
	log "github.com/WindowsSov8forUs/tsugu-bot-go/mylog"
)

func submitCarMessage(session adapter.Session, bot adapter.Bot, conf *config.Config) (bool, error) {
	message := session.Message()

	// 检测 Car 中的关键词
	found := false
	for _, keyward := range conf.Tsugu.CarConfig.Car {
		if strings.Contains(message, keyward) {
			found = true
			break
		}
	}
	if !found {
		return false, nil
	}

	// 检测 Fake 中的关键词
	for _, keyward := range conf.Tsugu.CarConfig.Fake {
		if strings.Contains(message, keyward) {
			return false, nil
		}
	}

	// 匹配车牌号
	carPattern := regexp.MustCompile(`^\d{5}|^\d{6}`)
	if !carPattern.MatchString(message) {
		return false, nil
	}

	platform := session.Platform()
	car, err := getCarForward(platform, session.UserID(), conf)
	if err != nil {
		log.Warnf("<Tsugu> 获取用户数据失败: %v", err)
	} else {
		if !car {
			return false, nil
		}
	}

	// 获取车牌号
	carID, err := strconv.Atoi(carPattern.FindString(message))
	if err != nil {
		return true, err
	}

	// 构建请求
	request := &RequestSubmitRoomNumber{
		Number:     carID,
		RawMessage: message,
		Platform:   session.Platform(),
		UserID:     session.UserID(),
		UserName:   session.UserName(),
		Time:       time.Now().Unix(),
	}
	if conf.Tsugu.CarStation.BandoriStationToken != "" {
		request.BandoriStationToken = conf.Tsugu.CarStation.BandoriStationToken
	}

	response, err := submitRoomNumber(request, conf)

	if err != nil {
		return true, err
	}

	if response.Status == STATUS_FAILED {
		log.Warnf("<Tsugu> 提交车牌号失败，%v", response.Data)
	} else {
		log.Infof("<Tsugu> 提交车牌号成功: %v", carID)
		if conf.Tsugu.CarStation.ForwardResponse {
			if conf.Tsugu.CarStation.ResponseContent == "" {
				forwardMessage := &adapter.Message{}
				forwardMessage.Text("车牌号提交成功")
				bot.Send(session, forwardMessage)
			} else {
				forwardMessage := &adapter.Message{}
				forwardMessage.Text(conf.Tsugu.CarStation.ResponseContent)
				bot.Send(session, forwardMessage)
			}
		}
	}

	return true, nil
}
