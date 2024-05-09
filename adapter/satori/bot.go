package satori

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/WindowsSov8forUs/tsugu-bot-go/adapter"
	"github.com/WindowsSov8forUs/tsugu-bot-go/config"
	log "github.com/WindowsSov8forUs/tsugu-bot-go/mylog"
	satoriMessage "github.com/satori-protocol-go/satori-model-go/pkg/message"
)

type Bot struct {
	apiBase string
	config  *config.Config
}

type RequestMessageCreate struct {
	ChannelId string `json:"channel_id"` // 频道 ID
	Content   string `json:"content"`    // 消息内容
}

func (bot *Bot) messageCreate(channelId, content, selfId, platform string) error {
	var client = &http.Client{}
	var apiURL = fmt.Sprintf("%s/message.create", bot.apiBase)
	var req = &RequestMessageCreate{
		ChannelId: channelId,
		Content:   content,
	}
	var reqBody, err = json.Marshal(req)
	if err != nil {
		return err
	}
	ApiURL, _ := url.Parse(apiURL)
	header := http.Header{}
	header.Add("Authorization", fmt.Sprintf("Bearer %s", bot.config.Satori.Token))
	header.Add("Content-Type", "application/json")
	header.Add("X-Self-Id", selfId)
	header.Add("X-Platform", platform)
	request := &http.Request{
		Method: "POST",
		URL:    ApiURL,
		Body:   io.NopCloser(strings.NewReader(string(reqBody))),
		Header: header,
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		log.Warnf("<Satori> message.create 返回状态码 %d", response.StatusCode)
	}
	return nil
}

func (bot *Bot) Send(session adapter.Session, message *adapter.Message) error {
	var elements []satoriMessage.MessageElement
	var logContent string
	for _, seg := range message.Segments {
		switch seg.Type {
		case "image":
			element := &satoriMessage.MessageElementImg{
				Src: fmt.Sprintf("data:image/png;base64,%s", seg.Data),
			}
			elements = append(elements, element)
			logContent += "[图片]"
		case "text":
			element := &satoriMessage.MessageElementText{
				Content: seg.Data,
			}
			elements = append(elements, element)
			logContent += seg.Data
		}
	}
	sess := session.GetSession().(*Session)
	quote := &satoriMessage.MessageElementQuote{}
	quote.ExtendAttributes = quote.AddAttribute("id", sess.Data.Message.Id)
	elements = append([]satoriMessage.MessageElement{quote}, elements...)
	log.Infof("<Satori> 发送消息至 %s : %s", session.ChannelID(), logContent)

	content, _ := satoriMessage.Stringify(elements)

	return bot.messageCreate(sess.Data.Channel.Id, content, sess.Data.SelfId, sess.Data.Platform)
}
