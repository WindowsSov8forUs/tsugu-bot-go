package tsugu

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/WindowsSov8forUs/tsugu-bot-go/config"
	log "github.com/WindowsSov8forUs/tsugu-bot-go/mylog"
)

type RequestBackend struct {
	DefaultServers []int  `json:"default_servers"`
	Server         int    `json:"server"`
	Text           string `json:"text"`
	UseEasyBG      bool   `json:"useEasyBG"`
	Compress       bool   `json:"compress"`
}

type RequestUser struct {
	Platform string `json:"platform"`
	Text     string `json:"text,omitempty"`
	UserID   string `json:"user_id"`
	Server   int    `json:"server"`
	PlayerID int    `json:"playerId,omitempty"`
	BindType bool   `json:"bindType"`
	Status   bool   `json:"status"`
}

func (data *RequestBackend) Stringify() string {
	return fmt.Sprintf("{DefaultServers: %v, Server: %d, Text: %s, UseEasyBG: %v, Compress: %v}", data.DefaultServers, data.Server, data.Text, data.UseEasyBG, data.Compress)
}

func (data *RequestUser) Stringify() string {
	return fmt.Sprintf("{Platform: %s, Text: %s, UserID: %s, Server: %d, PlayerID: %d, BindType: %v, Status: %v}", data.Platform, data.Text, data.UserID, data.Server, data.PlayerID, data.BindType, data.Status)
}

type ResponseData struct {
	Type   string `json:"type"`
	String string `json:"string"`
}

func TextResponse(text string) []*ResponseData {
	return []*ResponseData{
		{
			Type:   "string",
			String: text,
		},
	}
}

func requestPostUser(requestAPI string, data *RequestUser, conf *config.Config) ([]byte, error) {
	log.Tracef("<Tsugu> 请求用户数据后端地址: %s", conf.Tsugu.UserDataBackend.Url)
	log.Debugf("<Tsugu> 请求用户数据: api: %s, data:%s", requestAPI, data.Stringify())
	var transport *http.Transport
	if conf.Tsugu.UserDataBackend.UseProxy {
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

	requestBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	log.Tracef("<Tsugu> 发送用户数据后端数据: %s", string(requestBody))
	requestURL := fmt.Sprintf("%s%s", conf.Tsugu.UserDataBackend.Url, requestAPI)
	response, err := client.Post(requestURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK || response.StatusCode == http.StatusBadRequest {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		log.Tracef("<Tsugu> 用户数据后端返回数据: %s", string(body))
		return body, nil
	}

	return nil, errors.New("服务器出现问题，清稍后再试。")
}

func requestPostBackend(requestAPI string, data *RequestBackend, conf *config.Config) ([]*ResponseData, error) {
	log.Tracef("<Tsugu> 请求后端地址: %s", conf.Tsugu.Backend.Url)
	log.Debugf("<Tsugu> 请求后端数据: api: %s, data:%s", requestAPI, data.Stringify())
	var transport *http.Transport
	if conf.Tsugu.Backend.UseProxy {
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

	requestBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	log.Tracef("<Tsugu> 发送后端数据: %s", string(requestBody))
	requestURL := fmt.Sprintf("%s%s", conf.Tsugu.Backend.Url, requestAPI)
	response, err := client.Post(requestURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var body []byte
	if response.StatusCode == http.StatusOK {
		body, err = io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
	}
	log.Tracef("<Tsugu> 后端返回数据: %s", string(body))
	var responseData = make([]*ResponseData, 0)
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

func ApiV2Backend(api, text string, defaultServers []int, server int, conf *config.Config) ([]*ResponseData, error) {
	if len(defaultServers) == 0 {
		defaultServers = append(defaultServers, 3, 0)
	}
	path := fmt.Sprintf("/v2/%s", api)
	data := &RequestBackend{
		DefaultServers: defaultServers,
		Server:         server,
		Text:           text,
		UseEasyBG:      conf.Tsugu.UseEasyBG,
		Compress:       conf.Tsugu.Compress,
	}
	return requestPostBackend(path, data, conf)
}

func ApiV2Command(message, command, api, platform, userID, channelID string, conf *config.Config) ([]*ResponseData, error) {
	text := strings.Trim(strings.TrimPrefix(message, command), " ")

	if api == "cardIllustration" || api == "ycm" {
		// 不需要验证服务器信息
		return ApiV2Backend(api, text, nil, 3, conf)
	} else if api == "gachaSimulate" {
		for _, v := range conf.Tsugu.BanGachaSimulate {
			if v == channelID {
				return TextResponse("本群抽卡模拟功能已禁用"), nil
			}
		}
	}

	// 获取用户数据
	var userData *ResponseUserData
	var err error
	if conf.Tsugu.UserDataBasePath == "" {
		userData, err = remoteGetUserData(platform, userID, conf)
	} else {
		userData, err = GetUserData(userID, platform)
	}
	if err != nil {
		return nil, err
	}
	if userData.Status != "success" {
		log.Warnf("<Tsugu> 获取用户数据失败: %s:%s", platform, userID)
		return TextResponse("获取用户数据失败：内部错误"), nil
	} else {
		return ApiV2Backend(api, text, userData.Data.DefaultServer, userData.Data.ServerMode, conf)
	}
}
