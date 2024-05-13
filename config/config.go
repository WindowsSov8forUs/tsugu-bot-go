package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"

	log "github.com/WindowsSov8forUs/tsugu-bot-go/mylog"
	"gopkg.in/yaml.v3"
)

var (
	instance *Config
	mutex    sync.Mutex
)

// Config 配置
type Config struct {
	LogLevel log.LogLevel `yaml:"log_level"` // 日志等级
	Tsugu    *Tsugu       `yaml:"tsugu"`
	Satori   *Satori      `yaml:"satori"`
}

// Tsugu Tsugu 配置
type Tsugu struct {
	RequireAt        bool            `yaml:"require_at"` // 是否需要 @
	Reply            bool            `yaml:"reply"`      // 是否回复用户
	At               bool            `yaml:"at"`         // 是否 @ 用户
	NoSpace          bool            `yaml:"no_space"`   // 是否无需空格
	Timeout          int             `yaml:"timeout"`    // 超时
	Proxy            string          `yaml:"proxy"`      // 代理
	UseEasyBG        bool            `yaml:"use_easy_bg"`
	Compress         bool            `yaml:"compress"`
	UserDataBasePath string          `yaml:"user_database_path"` // 用户数据路径
	BanGachaSimulate []string        `yaml:"ban_gacha_simulate"`
	CarStation       CarStation      `yaml:"car_station"`   // 车站配置
	VerifyPlayer     VerifyPlayer    `yaml:"verify_player"` // 验证玩家配置
	Backend          Backend         `yaml:"backend"`
	UserDataBackend  UserDataBackend `yaml:"user_data_backend"`
	Functions        Functions       `yaml:"functions"`
	CommandAlias     CommandAlias    `yaml:"command_alias"`
	CarConfig        CarConfig       `yaml:"car_config"`
}

func (tsugu *Tsugu) RemoveBanList(channelId string) {
	mutex.Lock()
	defer mutex.Unlock()

	for i, v := range tsugu.BanGachaSimulate {
		if v == channelId {
			tsugu.BanGachaSimulate = append(tsugu.BanGachaSimulate[:i], tsugu.BanGachaSimulate[i+1:]...)
			return
		}
	}
}

func (tsugu *Tsugu) AddBanList(channelId string) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, v := range tsugu.BanGachaSimulate {
		if v == channelId {
			return
		}
	}
	tsugu.BanGachaSimulate = append(tsugu.BanGachaSimulate, channelId)
}

// Backend 后端配置
type Backend struct {
	Url      string `yaml:"url"`
	UseProxy bool   `yaml:"use_proxy"`
}

// UserDataBackend 用户数据后端配置
type UserDataBackend struct {
	Url      string `yaml:"url"`
	UseProxy bool   `yaml:"use_proxy"`
}

// CarStation 车站配置
type CarStation struct {
	BandoriStationToken string `yaml:"bandori_station_token"` // 车站令牌
	ForwardResponse     bool   `yaml:"forward_response"`      // 转发响应
	ResponseContent     string `yaml:"response_content"`      // 响应内容
}

// VerifyPlayer 验证玩家配置
type VerifyPlayer struct {
	UseProxy bool `yaml:"use_proxy"` // 使用代理
}

// Functions 功能开关配置
type Functions struct {
	Help                bool `yaml:"help"`                  // 帮助文档
	CarForward          bool `yaml:"car_forward"`           // 车牌转发
	SwitchGachaSimulate bool `yaml:"switch_gacha_simulate"` // 开关本群抽卡模拟
	SwitchCarForward    bool `yaml:"switch_car_forward"`    // 是否允许指令开启车牌转发
	BindPlayer          bool `yaml:"bind_player"`           // 绑定玩家
	ChangeMainServer    bool `yaml:"change_main_server"`    // 切换主服务器
	ChangeServerList    bool `yaml:"change_server_list"`    // 切换服务器列表
	PlayerStatus        bool `yaml:"player_status"`         // 玩家状态
	Ycm                 bool `yaml:"ycm"`                   // 获取车牌
	SearchPlayer        bool `yaml:"search_player"`         // 玩家信息
	SearchCard          bool `yaml:"search_card"`           // 查卡
	CardIllustration    bool `yaml:"card_illustration"`     // 查卡面
	SearchCharacter     bool `yaml:"search_character"`      // 查角色
	SearchEvent         bool `yaml:"search_event"`          // 查活动
	SearchSong          bool `yaml:"search_song"`           // 查歌曲
	SearchChart         bool `yaml:"search_chart"`          // 查谱面
	SongMeta            bool `yaml:"song_meta"`             // 查询分数表
	EventStage          bool `yaml:"event_stage"`           // 查活动试炼
	SearchGacha         bool `yaml:"search_gacha"`          // 查卡池
	Ycx                 bool `yaml:"ycx"`                   // 预测线
	YcxAll              bool `yaml:"ycx_all"`               // 全部预测线
	Lsycx               bool `yaml:"lsycx"`                 // 历史预测线
	GachaSimulate       bool `yaml:"gacha_simulate"`        // 模拟抽卡
}

type CommandAlias struct {
	SwitchGachaSimulate []string `yaml:"switch_gacha_simulate"` // 开关本群抽卡模拟
	OpenCarForward      []string `yaml:"open_car_forward"`      // 开启车牌转发
	CloseCarForward     []string `yaml:"close_car_forward"`     // 关闭车牌转发
	BindPlayer          []string `yaml:"bind_player"`           // 绑定玩家
	UnbindPlayer        []string `yaml:"unbind_player"`         // 解绑玩家
	ChangeMainServer    []string `yaml:"change_main_server"`    // 切换主服务器
	ChangeServerList    []string `yaml:"change_server_list"`    // 设置默认服务器
	Ycm                 []string `yaml:"ycm"`                   // 有车吗
	SearchPlayer        []string `yaml:"search_player"`         // 玩家信息
	SearchCard          []string `yaml:"search_card"`           // 查卡
	CardIllustration    []string `yaml:"card_illustration"`     // 查卡面
	SearchCharacter     []string `yaml:"search_character"`      // 查角色
	SearchEvent         []string `yaml:"search_event"`          // 查活动
	SearchSong          []string `yaml:"search_song"`           // 查歌曲
	SearchChart         []string `yaml:"search_chart"`          // 查谱面
	SongMeta            []string `yaml:"song_meta"`             // 查询分数表
	EventStage          []string `yaml:"event_stage"`           // 查活动试炼
	SearchGacha         []string `yaml:"search_gacha"`          // 查卡池
	Ycx                 []string `yaml:"ycx"`                   // ycx
	YcxAll              []string `yaml:"ycx_all"`               // ycxall
	Lsycx               []string `yaml:"lsycx"`                 // lsycx
	GachaSimulate       []string `yaml:"gacha_simulate"`        // 抽卡模拟
}

// CarConfig 车牌设置
type CarConfig struct {
	Car  []string `yaml:"car"`  // 有效关键词列表
	Fake []string `yaml:"fake"` // 无效关键词列表
}

// Satori Satori 协议配置
type Satori struct {
	Host    string `yaml:"host"`    // 主机地址
	Port    int    `yaml:"port"`    // 端口
	Path    string `yaml:"path"`    // 路径
	Version int    `yaml:"version"` // 版本
	Token   string `yaml:"token"`   // 鉴权令牌
}

// LoadConfig 加载配置
func LoadConfig(path string) (*Config, error) {
	mutex.Lock()
	defer mutex.Unlock()

	// 如果已经加载过配置，直接返回
	if instance != nil {
		return instance, nil
	}

	configData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err = yaml.Unmarshal(configData, config); err != nil {
		return nil, err
	}

	// 确保配置完整性
	if err = ensureConfigComplete(path); err != nil {
		return nil, err
	}

	instance = config
	return instance, nil
}

// ensureConfigComplete 检查配置是否完整
func ensureConfigComplete(path string) error {
	// 读取配置文件
	configData, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// 解析到结构体中
	currentConfig := &Config{}
	if err = yaml.Unmarshal(configData, currentConfig); err != nil {
		return err
	}

	// 解析默认配置模板
	defaultConfig := &Config{}
	if err = yaml.Unmarshal([]byte(ConfigTemplate), defaultConfig); err != nil {
		return err
	}

	// 使用反射找出缺失设置
	missingSettings, err := getMissingSettingsByReflection(currentConfig, defaultConfig)
	if err != nil {
		return err
	}

	// 使用文本比对找出缺失设置
	missingSettingsByText, err := getMissingSettingsByText(ConfigTemplate, string(configData))
	if err != nil {
		return err
	}

	// 合并缺失设置
	missingSettings = mergeMissingSettings(missingSettings, missingSettingsByText)

	// 如果有缺失设置，处理缺失配置行
	if len(missingSettings) > 0 {
		// 更新配置文件
		if err = recreateToConfigFile(path); err != nil {
			return err
		}

		fmt.Printf("配置文件已更新，原配置文件已被命名为 config_backup.yml ，请重新启动程序。")
		os.Exit(0)
	}

	return nil
}

// getMissingSettingsByReflection 使用反射来对比结构体并找出缺失的设置
func getMissingSettingsByReflection(currentConfig, defaultConfig *Config) (map[string]string, error) {
	missingSettings := make(map[string]string)
	currentVal := reflect.ValueOf(currentConfig).Elem()
	defaultVal := reflect.ValueOf(defaultConfig).Elem()

	for i := 0; i < currentVal.NumField(); i++ {
		field := currentVal.Type().Field(i)
		yamlTag := field.Tag.Get("yaml")
		if yamlTag == "" || field.Type.Kind() == reflect.Int || field.Type.Kind() == reflect.Bool {
			continue // 跳过没有yaml标签的字段，或者字段类型为int或bool
		}
		yamlKeyName := strings.SplitN(yamlTag, ",", 2)[0]
		if isZeroOfUnderlyingType(currentVal.Field(i).Interface()) && !isZeroOfUnderlyingType(defaultVal.Field(i).Interface()) {
			missingSettings[yamlKeyName] = "missing"
		}
	}

	return missingSettings, nil
}

func isZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

// getMissingSettingsByText compares settings in two strings line by line, looking for missing keys.
func getMissingSettingsByText(templateContent, currentConfigContent string) (map[string]string, error) {
	templateKeys := extractKeysFromString(templateContent)
	currentKeys := extractKeysFromString(currentConfigContent)

	missingSettings := make(map[string]string)
	for key := range templateKeys {
		if _, found := currentKeys[key]; !found {
			missingSettings[key] = "missing"
		}
	}

	return missingSettings, nil
}

// extractKeysFromString reads a string and extracts the keys (text before the colon).
func extractKeysFromString(content string) map[string]bool {
	keys := make(map[string]bool)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, ":") {
			key := strings.TrimSpace(strings.Split(line, ":")[0])
			keys[key] = true
		}
	}
	return keys
}

// mergeMissingSettings 合并由反射和文本比对找到的缺失设置
func mergeMissingSettings(reflectionSettings, textSettings map[string]string) map[string]string {
	for k, v := range textSettings {
		reflectionSettings[k] = v
	}
	return reflectionSettings
}

func recreateToConfigFile(path string) error {
	// 将原配置文件重命名为 config_backup.yml
	err := os.Rename(path, "config_backup.yml")
	if err != nil {
		return err
	}

	// 将配置模板写入配置文件
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(ConfigTemplate)
	if err != nil {
		return err
	}

	return nil
}
