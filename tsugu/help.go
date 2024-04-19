package tsugu

import (
	"github.com/WindowsSov8forUs/tsugu-bot-go/adapter"
	log "github.com/WindowsSov8forUs/tsugu-bot-go/mylog"
)

var (
	PLAYER_STATE_HELP = `查询指定玩家的状态信息。
使用示例：
	玩家状态 : 查询你当前默认服务器的玩家状态
	玩家状态  jp : 查询日服玩家状态
	玩家状态 2 : 查询您的第二个绑定玩家的状态(*此功能需要BOT支持)`
	SWITCH_CAR_FORWARD = `开启或关闭车牌转发功能，针对个人
使用示例：
	开启车牌转发 : 开启车牌转发功能
	关闭车牌转发 : 关闭车牌转发功能`
	BIND_PLAYER = `绑定玩家ID到当前账号
使用示例：
	绑定玩家 : 绑定玩家开始绑定玩家流程(*请根据BOT的提示进行操作)`
	UNBIND = `解除绑定玩家ID
使用示例：
	解除绑定 : 解除绑定玩家开始解绑玩家流程(*请根据BOT的提示进行操作)`
	SWITCH_MAIN_SERVER = `切换主服务器
使用示例：
	主服务器 jp : 切换主服务器到日服
	国服模式 : 切换主服务器到国服`
	SET_DEFAULT_SERVER = `可以改变查询时的资源优先顺序
设置默认服务器，使用空格分隔服务器列表
使用示例：
	设置默认服务器 国服 日服 : 将国服设置为第一服务器，日服设置为第二服务器`
	CARD_ILLUSTION = `根据卡片ID查询卡片插画
使用示例：
	查卡面 1399 :返回1399号卡牌的插画`
	SEARCH_PLAYER = `查询指定ID玩家的信息。省略服务器名时，默认从你当前的主服务器查询
使用示例：
	查玩家 10000000 : 查询你当前默认服务器中，玩家ID为10000000的玩家信息
	查玩家 40474621 jp : 查询日服玩家ID为40474621的玩家信息`
	GACHA_SIMULATION = `模拟抽卡，如果没有卡池ID的话，卡池为当前活动的卡池
使用示例：
	抽卡模拟:模拟抽卡10次
	抽卡模拟 300 922 :模拟抽卡300次，卡池为922号卡池`
	SEARCH_GACHA = `根据卡池ID查询卡池信息`
	SEARCH_EVENT = `根据关键词或活动ID查询活动信息
使用示例：
	查活动 177 :返回177号活动的信息
	查活动 绿 tsugu :返回所有属性加成为pure，且活动加成角色中包括羽泽鸫的活动列表
	查活动 >255 :返回所有活动ID大于255的活动列表
	查活动 255-256 :返回所有活动ID在255到256之间的活动列表
	查活动 ppp :匹配到 PPP 乐队的活动信息`
	FIND_SONG = `根据关键词或曲目ID查询曲目信息
使用示例：
	查曲 1 :返回1号曲的信息
	查曲 ag lv27 :返回所有难度为27的ag曲列表
	查曲 1 ex :返回1号曲的expert难度曲目信息
	查曲 滑滑蛋 :匹配到 ふわふわ時間 的曲目信息`
	FIND_SCORE     = `查询指定服务器的歌曲分数表，如果没有服务器名的话，服务器为用户的默认服务器`
	FIND_CHARACTER = `根据关键词或角色ID查询角色信息
使用示例：
	查角色 10 :返回10号角色的信息
	查角色 吉他 :返回所有角色模糊搜索标签中包含吉他的角色列表`
	FIND_CHART = `根据曲目ID与难度查询铺面信息
使用示例：
	查谱面 1 :返回1号曲的所有铺面
	查谱面 1 expert :返回1号曲的expert难度铺面`
	YCX_ALL = `查询所有档位的预测线，如果没有服务器名的话，服务器为用户的默认服务器。如果没有活动ID的话，活动为当前活动
可用档线:
20, 30, 40, 50, 100, 200, 300, 400, 500, 1000, 2000, 5000, 10000, 20000, 30000, 50000, 

使用示例：
    ycxall :返回默认服务器当前活动所有档位的档线与预测线
    ycxall 177 jp:返回日服177号活动所有档位的档线与预测线`
	YCX = `查询指定档位的预测线，如果没有服务器名的话，服务器为用户的默认服务器。如果没有活动ID的话，活动为当前活动
可用档线:
20, 30, 40, 50, 100, 200, 300, 400, 500, 1000, 2000, 5000, 10000, 20000, 30000, 50000, 
使用示例：
	ycx 1000 :返回默认服务器当前活动1000档位的档线与预测线
	ycx 1000 177 jp:返回日服177号活动1000档位的档线与预测线`
	LSYCX = `与ycx的区别是，lsycx会返回与最近的4期活动类型相同的活动的档线数据
查询指定档位的预测线，与最近的4期活动类型相同的活动的档线数据，如果没有服务器名的话，服务器为用户的默认服务器。如果没有活动ID的话，活动为当前活动
可用档线:
20, 30, 40, 50, 100, 200, 300, 400, 500, 1000, 2000, 5000, 10000, 20000, 30000, 50000, 

使用示例：
	lsycx 1000 :返回默认服务器当前活动的档线与预测线，与最近的4期活动类型相同的活动的档线数据
	lsycx 1000 177 jp:返回日服177号活动1000档位档线与最近的4期活动类型相同的活动的档线数据`
	YCM = `获取所有车牌车牌
使用示例：
	ycm : 获取所有车牌`
	SEARCH_CARD = `根据关键词或卡牌ID查询卡片信息，请使用空格隔开所有参数
使用示例：
	查卡 1399 :返回1399号卡牌的信息
	查卡 绿 tsugu :返回所有属性为pure的羽泽鸫的卡牌列表
	查卡 kfes ars :返回所有为kfes的ars的卡牌列表`
)

func helpDoc(command string) string {
	switch command {
	case "玩家状态":
		return PLAYER_STATE_HELP
	case "开关车牌转发":
		return SWITCH_CAR_FORWARD
	case "绑定玩家":
		return BIND_PLAYER
	case "解除绑定":
		return UNBIND
	case "切换主服务器":
		return SWITCH_MAIN_SERVER
	case "设置默认服务器":
		return SET_DEFAULT_SERVER
	case "查卡面":
		return CARD_ILLUSTION
	case "查询玩家":
		return SEARCH_PLAYER
	case "卡池模拟":
		return GACHA_SIMULATION
	case "查卡池":
		return SEARCH_GACHA
	case "查活动":
		return SEARCH_EVENT
	case "查曲":
		return FIND_SONG
	case "查分数表":
		return FIND_SCORE
	case "查角色":
		return FIND_CHARACTER
	case "查谱面":
		return FIND_CHART
	case "ycxall":
		return YCX_ALL
	case "ycx":
		return YCX
	case "lsycx":
		return LSYCX
	case "ycm":
		return YCM
	case "查卡":
		return SEARCH_CARD
	default:
		return ""
	}
}

func helpCommand(session adapter.Session, command string, bot adapter.Bot) error {
	if command == "" {
		return nil
	}
	help := helpDoc(command)
	if help != "" {
		log.Infof("<Tsugu> 发送帮助信息: %s", command)
		message := &adapter.Message{}
		message.Text(help)
		return bot.Send(session, message)
	}
	return nil
}
