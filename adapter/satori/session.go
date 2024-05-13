package satori

import (
	"fmt"

	"github.com/WindowsSov8forUs/tsugu-bot-go/adapter/satori/event"
	"github.com/satori-protocol-go/satori-model-go/pkg/channel"
	"github.com/satori-protocol-go/satori-model-go/pkg/message"
)

type Session struct {
	Data *event.Event // 事件
}

func (s *Session) Message() string {
	elements, err := message.Parse(s.Data.Message.Content)
	if err != nil {
		return ""
	}
	var msg string
	for _, element := range elements {
		switch e := element.(type) {
		case *message.MessageElementText:
			msg += e.Content
		default:
			continue
		}
	}
	return msg
}

func (s *Session) UserID() string {
	return s.Data.User.Id
}

func (s *Session) UserName() string {
	if s.Data.Member != nil && s.Data.Member.Nick != "" {
		return s.Data.Member.Nick
	} else {
		if s.Data.User != nil {
			return s.Data.User.Name
		} else {
			return "Unknown"
		}
	}
}

func (s *Session) Platform() string {
	return "red"
}

func (s *Session) ChannelID() string {
	return s.Data.Channel.Id
}

func (s *Session) GetSession() interface{} {
	return s
}

func (s *Session) Log() string {
	var logContent string
	guild := s.Data.Guild
	chnl := s.Data.Channel
	user := s.Data.User
	member := s.Data.Member
	if guild != nil {
		if guild.Id != chnl.Id {
			logContent += fmt.Sprintf("%s(%s)-", guild.Name, guild.Id)
		}
	}
	if chnl.Type != channel.CHANNEL_TYPE_DIRECT {
		logContent += fmt.Sprintf("%s(%s)-", chnl.Name, chnl.Id)
	}
	if member != nil && member.Nick != "" {
		logContent += member.Nick
	} else {
		logContent += user.Name
	}
	logContent += fmt.Sprintf("(%s): ", user.Id)

	elements, err := message.Parse(s.Data.Message.Content)
	if err != nil {
		return ""
	}
	for _, element := range elements {
		switch e := element.(type) {
		case *message.MessageElementText:
			logContent += e.Content
		default:
			continue
		}
	}

	return logContent
}
