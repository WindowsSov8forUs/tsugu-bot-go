package api

const TsuguBackend string = "http://tsugubot.com:8080"

type Config struct {
	Proxy                string // 代理服务器
	Timeout              int    // 超时时间
	BackendUrl           string // 后端地址
	BackendProxy         bool   // 后端是否使用代理
	DatabaseBackendUrl   string // 数据库后端地址
	DatabaseBackendProxy bool   // 数据库后端是否使用代理
	UseEasyBG            bool   // 是否使用简易背景
	Compress             bool   // 是否压缩
}
