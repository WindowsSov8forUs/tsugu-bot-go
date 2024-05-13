package tsugu

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/WindowsSov8forUs/tsugu-bot-go/config"
	log "github.com/WindowsSov8forUs/tsugu-bot-go/mylog"
)

type Type string

const (
	TYPE_STRING Type = "string"
	TYPE_BASE64 Type = "base64"
)

type Status string

const (
	STATUS_SUCCESS Status = "success"
	STATUS_FAILED  Status = "failed"
)

type RequestEventStage struct {
	Server   ServerId `json:"server"`             // 服务器
	EventID  int      `json:"event_id,omitempty"` // 活动 ID
	Meta     bool     `json:"meta,omitempty"`     // 是否查询活动信息
	Compress bool     `json:"compress,omitempty"` // 是否压缩
}

type RequestGachaSimulate struct {
	ServerMode ServerId `json:"server_mode"`        // 服务器
	Times      int      `json:"times,omitempty"`    // 次数
	GachaID    int      `json:"gacha_id,omitempty"` // 卡池 ID
	Compress   bool     `json:"compress,omitempty"` // 是否压缩
}

type RequestGetCardIllustration struct {
	CardID int `json:"cardId"` // 卡片 ID
}

type RequestLsycx struct {
	Server   ServerId `json:"server"`             // 服务器
	Tier     int      `json:"tier"`               // 档位
	EventID  int      `json:"event_id,omitempty"` // 活动 ID
	Compress bool     `json:"compress,omitempty"` // 是否压缩
}

type RequestRoomList struct {
	RoomList []*Room `json:"roomList"`           // 房间列表
	Compress bool    `json:"compress,omitempty"` // 是否压缩
}

type RequestSearchCard struct {
	DefaultServers []ServerId `json:"default_servers"`    // 默认服务器
	Text           string     `json:"text"`               // 文本
	UseEasyBG      bool       `json:"useEasyBG"`          // 是否使用简易背景
	Compress       bool       `json:"compress,omitempty"` // 是否压缩
}

type RequestSearchCharacter struct {
	DefaultServers []ServerId `json:"default_servers"`    // 默认服务器
	Text           string     `json:"text"`               // 文本
	Compress       bool       `json:"compress,omitempty"` // 是否压缩
}

type RequestSearchEvent struct {
	DefaultServers []ServerId `json:"default_servers"`    // 默认服务器
	Text           string     `json:"text"`               // 文本
	UseEasyBG      bool       `json:"useEasyBG"`          // 是否使用简易背景
	Compress       bool       `json:"compress,omitempty"` // 是否压缩
}

type RequestSearchGacha struct {
	DefaultServers []ServerId `json:"default_servers"`    // 默认服务器
	GachaID        int        `json:"gacha_id"`           // 卡池 ID
	UseEasyBG      bool       `json:"useEasyBG"`          // 是否使用简易背景
	Compress       bool       `json:"compress,omitempty"` // 是否压缩
}

type RequestSearchPlayer struct {
	PlayerID  int      `json:"playerId"`           // 玩家 ID
	Server    ServerId `json:"server"`             // 服务器
	UseEasyBG bool     `json:"useEasyBG"`          // 是否使用简易背景
	Compress  bool     `json:"compress,omitempty"` // 是否压缩
}

type RequestSearchSong struct {
	DefaultServers []ServerId `json:"default_servers"`    // 默认服务器
	Text           string     `json:"text"`               // 文本
	Compress       bool       `json:"compress,omitempty"` // 是否压缩
}

type RequestSongChart struct {
	DefaultServers []ServerId `json:"default_servers"`    // 默认服务器
	SongID         int        `json:"songId"`             // 歌曲 ID
	DifficultyText string     `json:"difficultyText"`     // 难度文本
	Compress       bool       `json:"compress,omitempty"` // 是否压缩
}

type RequestSongMeta struct {
	DefaultServers []ServerId `json:"default_servers"`    // 默认服务器
	Server         ServerId   `json:"server"`             // 服务器
	Compress       bool       `json:"compress,omitempty"` // 是否压缩
}

type RequestYcx struct {
	Server   ServerId `json:"server"`             // 服务器
	Tier     int      `json:"tier"`               // 档位
	EventID  int      `json:"event_id,omitempty"` // 活动 ID
	Compress bool     `json:"compress,omitempty"` // 是否压缩
}

type RequestYcxAll struct {
	Server   ServerId `json:"server"`             // 服务器
	EventID  int      `json:"event_id,omitempty"` // 活动 ID
	Compress bool     `json:"compress,omitempty"` // 是否压缩
}

type RequestSubmitRoomNumber struct {
	Number              int    `json:"number"`                        // 房间号
	RawMessage          string `json:"rawMessage"`                    // 原始消息
	Platform            string `json:"platform"`                      // 平台
	UserID              string `json:"user_id"`                       // 用户 ID
	UserName            string `json:"userName"`                      // 用户名
	Time                int64  `json:"time"`                          // 时间
	BandoriStationToken string `json:"bandoriStationToken,omitempty"` // 车站令牌
}

type RequestUserData struct {
	Platform string           `json:"platform"`
	UserID   string           `json:"user_id"`
	Update   *TsuguUserUpdate `json:"update,omitempty"`
}

type RequestBindPlayerRequest struct {
	Platform string   `json:"platform"`
	UserID   string   `json:"user_id"`
	Server   ServerId `json:"server"`
	BindType bool     `json:"bindType"`
}

type RequestBindPlayerVerification struct {
	Platform string   `json:"platform"`
	UserID   string   `json:"user_id"`
	Server   ServerId `json:"server"`
	PlayerID int      `json:"playerId"`
	BindType bool     `json:"bindType"`
}

type ResponseData struct {
	Type   Type   `json:"type"`
	String string `json:"string"`
}

type Player struct {
	ID     int      `json:"id"`
	Server ServerId `json:"server"`
}

type Room struct {
	Number     int    `json:"number"`
	RawMessage string `json:"rawMessage"`
	Source     string `json:"source"`
	UserID     string `json:"user_id"`
	Time       int64  `json:"time"`
	Player     Player `json:"player"`
	Avanter    string `json:"avanter,omitempty"`
	UserName   string `json:"userName,omitempty"`
}

type ResponseQueryAllRoomSuccess struct {
	Status Status  `json:"status"`
	Data   []*Room `json:"data"`
}

type ResponseQueryAllRoomFailed struct {
	Status Status `json:"status"`
	Data   string `json:"data"`
}

type ResponseSubmitRoomNumber ResponseQueryAllRoomFailed

func TextResponse(text string) []*ResponseData {
	return []*ResponseData{
		{
			Type:   "string",
			String: text,
		},
	}
}

func requestPostUser(requestAPI Api, data interface{}, conf *config.Config) ([]byte, error) {
	log.Tracef("<Tsugu> 请求用户数据后端地址: %s", conf.Tsugu.UserDataBackend.Url)
	log.Debugf("<Tsugu> 请求用户数据: api: %s, data: {%v}", requestAPI, data)
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

func requestBackend(requestAPI Api, data interface{}, conf *config.Config) ([]*ResponseData, error) {
	log.Tracef("<Tsugu> 请求后端地址: %s", conf.Tsugu.Backend.Url)
	log.Debugf("<Tsugu> 请求后端数据: api: %s, data: {%v}", requestAPI, data)
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
	body, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	log.Tracef("<Tsugu> 后端返回数据: %s", string(body))

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusBadRequest {
		return nil, fmt.Errorf("%s", string(body))
	}

	var responseData = make([]*ResponseData, 0)
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

func requestGetBackend(requestAPI Api, params *map[string]string, conf *config.Config) (*http.Response, error) {
	log.Tracef("<Tsugu> 请求后端地址: %s", conf.Tsugu.Backend.Url)
	log.Debugf("<Tsugu> 请求后端数据: api: %s, params: {%v}", requestAPI, params)
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

	var requestURL string
	if params == nil {
		requestURL = fmt.Sprintf("%s%s", conf.Tsugu.Backend.Url, requestAPI)
	} else {
		urlParams := url.Values{}
		for key, value := range *params {
			urlParams.Add(key, value)
		}
		log.Tracef("<Tsugu> 发送后端数据: %v", urlParams)
		requestURL = fmt.Sprintf("%s%s?%s", conf.Tsugu.Backend.Url, requestAPI, urlParams.Encode())
	}
	return client.Get(requestURL)
}

func requestPostBackend(requestAPI Api, data interface{}, conf *config.Config) (*http.Response, error) {
	log.Tracef("<Tsugu> 请求后端地址: %s", conf.Tsugu.Backend.Url)
	log.Debugf("<Tsugu> 请求后端数据: api: %s, data: {%v}", requestAPI, data)
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
	return client.Post(requestURL, "application/json", bytes.NewBuffer(requestBody))
}

func submitRoomNumber(data *RequestSubmitRoomNumber, conf *config.Config) (*ResponseSubmitRoomNumber, error) {
	response, err := requestPostBackend(ApiStationSubmitRoomNumber, data, conf)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var body []byte
	body, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	log.Tracef("<Tsugu> 后端返回数据: %s", string(body))

	if response.StatusCode == http.StatusOK || response.StatusCode == http.StatusBadRequest {
		var responseData = &ResponseSubmitRoomNumber{}
		err = json.Unmarshal(body, responseData)
		if err != nil {
			return nil, err
		}
		return responseData, nil
	} else {
		return nil, fmt.Errorf("%s", string(body))
	}
}

func queryAllRoom(conf *config.Config) (*ResponseQueryAllRoomSuccess, error) {
	response, err := requestGetBackend(ApiStationQueryAllRoom, nil, conf)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var body []byte
	body, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	log.Tracef("<Tsugu> 后端返回数据: %s", string(body))

	if response.StatusCode == http.StatusOK {
		var responseData = &ResponseQueryAllRoomSuccess{}
		err = json.Unmarshal(body, responseData)
		if err != nil {
			return nil, err
		}
		return responseData, nil
	} else {
		if response.StatusCode == http.StatusBadRequest {
			var responseData = &ResponseQueryAllRoomFailed{}
			err = json.Unmarshal(body, responseData)
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("%s", responseData.Data)
		} else {
			return nil, fmt.Errorf("%s", string(body))
		}
	}
}
