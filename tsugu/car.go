package tsugu

import (
	"net/http"
	"net/url"
	"regexp"
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
	carPattern := regexp.MustCompile(`^\d{5}(\D|$)|^\d{6}(\D|$)`)
	if !carPattern.MatchString(message) {
		return false, nil
	}

	platform := session.Platform()
	var userData *ResponseUserData
	var err error
	if conf.Tsugu.UserDataBasePath == "" {
		userData, err = remoteGetUserData(platform, session.UserID(), conf)
	} else {
		userData, err = GetUserData(platform, session.UserID())
	}
	if err != nil {
		log.Warnf("<Tsugu> 获取用户数据失败: %v", err)
	} else {
		if !userData.Data.Car {
			return false, nil
		}
	}

	// 获取车牌号
	carID := carPattern.FindString(message)

	// 构建 URL
	params := url.Values{}
	params.Add("function", "submit_room_number")
	params.Add("number", carID)
	params.Add("user_id", session.UserID())
	params.Add("raw_message", session.Message())
	params.Add("source", conf.Tsugu.CarStation.TokenName)
	params.Add("token", conf.Tsugu.CarStation.BandoriStationToken)
	stationURL := "https://api.bandoristation.com/index.php?" + params.Encode()

	var transport *http.Transport
	if conf.Tsugu.CarStation.UseProxy {
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

	response, err := client.Get(stationURL)
	if err != nil {
		return true, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Warnf("<Tsugu> 提交车牌号失败: %v", response.StatusCode)
	} else {
		log.Infof("<Tsugu> 提交车牌号成功: %s", carID)
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
