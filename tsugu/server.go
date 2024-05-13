package tsugu

type ServerId int8

const (
	JP ServerId = iota // 日服 0
	EN                 // 国际服 1
	TW                 // 台服 2
	CN                 // 国服 3
	KR                 // 韩服 4
)

type BindingStatus int8

const (
	None      BindingStatus = iota // 无
	Verifying                      // 验证中
	Success                        // 成功
)

func serverNameToId(server string) ServerId {
	switch server {
	case "0":
		return JP
	case "1":
		return EN
	case "2":
		return TW
	case "3":
		return CN
	case "4":
		return KR
	case "jp":
		return JP
	case "en":
		return EN
	case "tw":
		return TW
	case "cn":
		return CN
	case "kr":
		return KR
	case "日服":
		return JP
	case "国际服":
		return EN
	case "台服":
		return TW
	case "国服":
		return CN
	case "韩服":
		return KR
	default:
		return -1
	}
}

func serverIdToShortName(server ServerId) string {
	switch server {
	case JP:
		return "jp"
	case EN:
		return "en"
	case TW:
		return "tw"
	case CN:
		return "cn"
	case KR:
		return "kr"
	default:
		return ""
	}
}

func serverIdToFullName(server ServerId) string {
	switch server {
	case JP:
		return "日服"
	case EN:
		return "国际服"
	case TW:
		return "台服"
	case CN:
		return "国服"
	case KR:
		return "韩服"
	default:
		return ""
	}
}
