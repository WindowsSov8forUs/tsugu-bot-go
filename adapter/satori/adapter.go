package satori

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/WindowsSov8forUs/tsugu-bot-go/adapter/satori/event"
	"github.com/WindowsSov8forUs/tsugu-bot-go/adapter/satori/operation"
	"github.com/WindowsSov8forUs/tsugu-bot-go/config"
	log "github.com/WindowsSov8forUs/tsugu-bot-go/mylog"
	"github.com/WindowsSov8forUs/tsugu-bot-go/tsugu"
	"github.com/gorilla/websocket"
	"github.com/satori-protocol-go/satori-model-go/pkg/message"
)

type Client struct {
	Host            string          // 主机地址
	Port            int             // 端口
	Path            string          // 路径
	Version         string          // 版本
	Token           string          // 鉴权令牌
	Config          *config.Config  // 配置
	Connection      *websocket.Conn // 连接
	Sequence        int64           // 序列号
	bot             *Bot            // 机器人
	heartbeatSignal chan bool       // 心跳信号
}

func (c *Client) Run() error {
	wsURL := url.URL{
		Scheme: "ws",
		Host:   fmt.Sprintf("%s:%d", c.Host, c.Port),
		Path:   fmt.Sprintf("%s/%s/events", c.Path, c.Version),
	}

	connection, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		return err
	}
	c.Connection = connection

	if err = c.identify(); err != nil {
		return err
	}

	_, message, err := c.Connection.ReadMessage()
	if err != nil {
		return err
	}
	op := &operation.Operation{}
	if err = json.Unmarshal(message, op); err != nil {
		return err
	}
	if body, err := json.Marshal(op.Body); err == nil {
		var readyBody = operation.ReadyBody{}
		if err = json.Unmarshal(body, &readyBody); err == nil {
			for _, login := range readyBody.Logins {
				log.Infof("<Satori> 连接至账号：%s(%s) in %s", login.User.Name, login.SelfId, login.Platform)
			}
		} else {
			return err
		}
	} else {
		return err
	}

	go c.sendHeartbeat()
	go c.receiveEvent()

	return nil
}

func (c *Client) Close() error {
	c.heartbeatSignal <- true
	return c.Connection.Close()
}

func (c *Client) resume() error {
	if c.Connection != nil {
		c.Connection.Close()
		c.Connection = nil
	}

	wsURL := url.URL{
		Scheme: "ws",
		Host:   fmt.Sprintf("%s:%d", c.Host, c.Port),
		Path:   fmt.Sprintf("%s/%s/events", c.Path, c.Version),
	}

	connection, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		return err
	}
	c.Connection = connection

	if err = c.identify(); err != nil {
		return err
	}

	_, message, err := c.Connection.ReadMessage()
	if err != nil {
		return err
	}
	op := &operation.Operation{}
	if err = json.Unmarshal(message, op); err != nil {
		return err
	}
	if readyBody, ok := op.Body.(operation.ReadyBody); ok {
		for _, login := range readyBody.Logins {
			log.Infof("<Satori> 连接至账号：%s(%s) in %s", login.User.Name, login.SelfId, login.Platform)
		}
	} else {
		return fmt.Errorf("无法解析 ReadyBody")
	}

	go c.sendHeartbeat()
	go c.receiveEvent()

	return nil
}

func (c *Client) sendHeartbeat() {
	// 发送心跳

	// 开始一个 0s 的计时器
	timer := time.NewTimer(0 * time.Second)
	for {
		select {
		case <-timer.C:
			// 发送心跳
			op := &operation.Operation{
				Op: operation.OpCodePing,
			}
			if err := c.Connection.WriteJSON(op); err != nil {
				log.Errorf("<Satori> 发送心跳失败：%s", err)
				c.resume()
				return
			}
			timer.Reset(10 * time.Second)
		case <-c.heartbeatSignal:
			// 停止心跳
			return
		}
	}
}

func (c *Client) receiveEvent() {
	for {
		// 读取信令
		_, message, err := c.Connection.ReadMessage()
		if err != nil {
			log.Errorf("<Satori> 读取信令失败：%s", err)
			c.heartbeatSignal <- true
			c.resume()
			return
		}
		// 解析信令
		op := &operation.Operation{}
		if err = json.Unmarshal(message, op); err != nil {
			log.Warnf("<Satori> 解析信令失败：%s", err)
			continue
		}
		if op.Op == operation.OpCodeEvent {
			// 处理事件
			body, err := json.Marshal(op.Body)
			if err != nil {
				log.Warnf("<Satori> 解析事件失败：%s", err)
				continue
			}
			e := &event.Event{}
			if err = json.Unmarshal(body, e); err != nil {
				log.Warnf("<Satori> 解析事件失败：%s", err)
				continue
			}
			c.Sequence = e.Id
			go c.handleEvent(e)
		}
	}
}

func (c *Client) handleEvent(e *event.Event) {
	if e.Type == event.EventTypeMessageCreated {
		// 处理消息创建事件
		session := &Session{
			Data: e,
		}
		log.Infof("<Satori> %s", session.Log())
		elements, err := message.Parse(e.Message.Content)
		if err != nil {
			log.Errorf("<Satori> %s", err)
			return
		}
		var at = false
		for _, element := range elements {
			if elmt, ok := element.(*message.MessageElementAt); ok {
				if elmt.Id == e.SelfId {
					at = true
				}
			}
		}
		if c.Config.Tsugu.RequireAt {
			if !at {
				return
			}
		}
		if e.User.Id == e.SelfId {
			return
		}
		err = tsugu.Handler(session, c.bot, c.bot.config)
		if err != nil {
			log.Errorf("<Satori> %s", err)
		}
	}
}

func (c *Client) identify() error {
	op := &operation.Operation{
		Op: operation.OpCodeIdentify,
		Body: &operation.IdentifyBody{
			Token:    c.Token,
			Sequence: c.Sequence,
		},
	}
	return c.Connection.WriteJSON(op)
}

func NewClient(conf *config.Config) (*Client, error) {
	var version, path string
	if conf.Satori.Version == 0 {
		version = "v1"
	} else {
		version = fmt.Sprintf("v%d", conf.Satori.Version)
	}
	if strings.HasPrefix(conf.Satori.Path, "/") {
		path = conf.Satori.Path
	} else {
		if conf.Satori.Path == "" {
			path = ""
		} else {
			path = "/" + conf.Satori.Path
		}
	}
	apiBase := fmt.Sprintf("http://%s:%d%s/%s", conf.Satori.Host, conf.Satori.Port, path, version)
	bot := &Bot{
		apiBase: apiBase,
		config:  conf,
	}
	return &Client{
		Host:            conf.Satori.Host,
		Port:            conf.Satori.Port,
		Path:            path,
		Version:         version,
		Token:           conf.Satori.Token,
		Config:          conf,
		bot:             bot,
		heartbeatSignal: make(chan bool),
	}, nil
}
