package api

import "fmt"

type Server int8

type ServerList []Server

const (
	ServerJP Server = iota
	ServerEN
	ServerTW
	ServerCN
	ServerKR
)

type ApiEndpoint string

type ApiResponseType string

const (
	TypeString ApiResponseType = "string"
	TypeBase64 ApiResponseType = "base64"
)

type ErrBadStatus struct {
	StatusCode int
	Message    string
}

func (e *ErrBadStatus) Error() string {
	return fmt.Sprintf("status code: %d, message: %s", e.StatusCode, e.Message)
}

type ApiResponse struct {
	Type   ApiResponseType `json:"type"`
	String string          `json:"string"`
}

type ApiResponseStatus string

const (
	StatusFailed  ApiResponseStatus = "failed"
	StatusSuccess ApiResponseStatus = "success"
)

type ResponseFailed struct {
	Status ApiResponseStatus `json:"status"`
	Data   string            `json:"data"`
}

type Player struct {
	Id     int    `json:"id"`
	Server Server `json:"server"`
}

type Room struct {
	Number     int     `json:"number"`
	RawMessage string  `json:"rawMessage"`
	Source     string  `json:"source"`
	UserId     string  `json:"userId"`
	Time       int64   `json:"time"`
	Player     *Player `json:"player"`
	Avanter    string  `json:"avanter,omitempty"`
	UserName   string  `json:"userName,omitempty"`
}

type tsuguUser struct {
	UserId        string                   `json:"user_id"`
	Platform      string                   `json:"platform"`
	ServerMode    Server                   `json:"server_mode"`
	DefaultServer []Server                 `json:"default_server"`
	Car           bool                     `json:"car"`
	ServerList    []*tsuguUserServerInList `json:"server_list"`
}

type tsuguUserServerInList struct {
	PlayerId      int           `json:"playerId"`
	BindingStatus BindingStatus `json:"bindingStatus"`
	VerifyCode    int           `json:"verifyCode,omitempty"`
}

type BindingStatus int8

const (
	BindingStatusNone BindingStatus = iota
	BindingStatusVerifying
	BindingStatusSuccess
)

type PartialTsuguUser struct {
	UserId        string                   `json:"user_id,omitempty"`
	Platform      string                   `json:"platform,omitempty"`
	ServerMode    Server                   `json:"server_mode,omitempty"`
	DefaultServer []Server                 `json:"default_server,omitempty"`
	Car           bool                     `json:"car,omitempty"`
	ServerList    []*tsuguUserServerInList `json:"server_list,omitempty"`
}
