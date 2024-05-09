package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/WindowsSov8forUs/tsugu-bot-go/adapter/satori"
	"github.com/WindowsSov8forUs/tsugu-bot-go/config"
	log "github.com/WindowsSov8forUs/tsugu-bot-go/mylog"
	"github.com/WindowsSov8forUs/tsugu-bot-go/sys"
	"github.com/WindowsSov8forUs/tsugu-bot-go/tsugu"
)

func main() {
	// 定义 faststart 命令行标志，默认为 false
	fastStart := flag.Bool("faststart", false, "是否快速启动")

	// 解析命令行参数到定义的标志
	flag.Parse()

	// 检查是否使用了 -faststart 参数
	if !*fastStart {
		sys.InitBase()
	}

	// 检查 config.yml 是否存在
	if _, err := os.Stat("config.yml"); os.IsNotExist(err) {
		var err error
		configData := config.ConfigTemplate

		// 写入 config.yml
		err = os.WriteFile("config.yml", []byte(configData), 0644)
		if err != nil {
			log.Fatalf("写入配置文件时出错: %v", err)
			return
		}

		log.Info("已生成默认配置文件 config.yml，请修改后重启程序")
		fmt.Println("按下任意键继续...")
		fmt.Scanln()
		os.Exit(0)
	}

	// 加载配置
	conf, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("加载配置文件时出错: %v", err)
		return
	}

	if conf.Tsugu.UserDataBasePath != "" {
		log.Infof("<Tsugu> 正在加载本地用户数据库路径: %s", conf.Tsugu.UserDataBasePath)
		err := tsugu.DataBase(conf)
		if err != nil {
			log.Errorf("<Tsugu> 加载本地用户数据库时出错: %v", err)
			conf.Tsugu.UserDataBasePath = ""
		} else {
			log.Info("<Tsugu> 本地用户数据库已加载")
		}
	}

	var satoriClient *satori.Client
	if conf.Satori != nil {
		log.Info("已加载 Satori 配置")
		satoriClient, err = satori.NewClient(conf)
		if err != nil {
			log.Errorf("初始化 Satori 客户端时出错: %v", err)
		} else {
			err := satoriClient.Run()
			if err != nil {
				log.Errorf("运行 Satori 客户端时出错: %v", err)
			} else {
				log.Info("Satori 客户端已启动")
			}
		}
	}

	// 配置日志等级
	log.SetLogLevel(conf.LogLevel)

	// 使用通道来等待信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号
	<-sigCh

	log.Info("正在关闭 Tsugu 机器人客户端...")
	if satoriClient != nil {
		satoriClient.Close()
	}
}
