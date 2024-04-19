package operation

import (
	"github.com/WindowsSov8forUs/tsugu-bot-go/adapter/satori/event"
	"github.com/satori-protocol-go/satori-model-go/pkg/login"
)

type OpCode int

const (
	OpCodeEvent    OpCode = iota // 0
	OpCodePing                   // 1
	OpCodePong                   // 2
	OpCodeIdentify               // 3
	OpCodeReady                  // 4
)

type Operation struct {
	Op   OpCode      `json:"op"`             // 信令类型
	Body interface{} `json:"body,omitempty"` // 信令数据
}

type EventBody event.Event

type IdentifyBody struct {
	Token    string `json:"token,omitempty"`    // 鉴权令牌
	Sequence int64  `json:"sequence,omitempty"` // 序列号
}

type ReadyBody struct {
	Logins []*login.Login `json:"logins"` // 登录信息
}
