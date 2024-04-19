package tsugu

import (
	"strings"

	"github.com/WindowsSov8forUs/tsugu-bot-go/config"
)

type command struct {
	api   string
	names []string
}

var Commands = []*command{
	{
		api:   "cardIllustration",
		names: []string{"查插画", "查卡面"},
	},
	{
		api:   "player",
		names: []string{"查玩家", "查询玩家"},
	},
	{
		api:   "gachaSimulate",
		names: []string{"抽卡模拟", "卡池模拟"},
	},
	{
		api:   "gacha",
		names: []string{"查卡池"},
	},
	{
		api:   "event",
		names: []string{"查活动"},
	},
	{
		api:   "song",
		names: []string{"查歌曲", "查曲"},
	},
	{
		api:   "songMeta",
		names: []string{"查询分数表", "查分数表"},
	},
	{
		api:   "character",
		names: []string{"查角色"},
	},
	{
		api:   "chart",
		names: []string{"查铺面", "查谱面"},
	},
	{
		api:   "ycxAll",
		names: []string{"ycxall", "ycx all"},
	},
	{
		api:   "ycx",
		names: []string{"ycx", "预测线"},
	},
	{
		api:   "lsycx",
		names: []string{"lsycx"},
	},
	{
		api:   "ycm",
		names: []string{"ycm", "车来"},
	},
	{
		api:   "card",
		names: []string{"查卡"},
	},
}

func matchCommand(message string, conf *config.Config) (string, string) {
	for _, c := range Commands {
		for _, name := range c.names {
			if strings.HasPrefix(message, name) {
				switch c.api {
				case "cardIllustration":
					if conf.Tsugu.Functions.CardIllustration {
						return name, c.api
					}
				case "player":
					if conf.Tsugu.Functions.Player {
						return name, c.api
					}
				case "gachaSimulate":
					if conf.Tsugu.Functions.GachaSimulate {
						return name, c.api
					}
				case "gacha":
					if conf.Tsugu.Functions.Gacha {
						return name, c.api
					}
				case "event":
					if conf.Tsugu.Functions.Event {
						return name, c.api
					}
				case "song":
					if conf.Tsugu.Functions.Song {
						return name, c.api
					}
				case "songMeta":
					if conf.Tsugu.Functions.SongMeta {
						return name, c.api
					}
				case "character":
					if conf.Tsugu.Functions.Character {
						return name, c.api
					}
				case "chart":
					if conf.Tsugu.Functions.Chart {
						return name, c.api
					}
				case "ycxAll":
					if conf.Tsugu.Functions.YcxAll {
						return name, c.api
					}
				case "ycx":
					if conf.Tsugu.Functions.Ycx {
						return name, c.api
					}
				case "lsycx":
					if conf.Tsugu.Functions.Lsycx {
						return name, c.api
					}
				case "ycm":
					if conf.Tsugu.Functions.Ycm {
						return name, c.api
					}
				case "card":
					if conf.Tsugu.Functions.Card {
						return name, c.api
					}
				}
			}
		}
	}
	return "", ""
}
