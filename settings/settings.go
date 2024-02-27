package settings

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"time"
)

// 全局变量
var Conf = new(Appconfig)

type Appconfig struct {
	Port     int `mapstruceture:"port"`
	*Webhook `mapstruceture:"webhook"`
	*Db      `mapstruceture:"db"`
	*Watch   `mapstruceture:"watch"`
}
type Watch struct {
	Threshold      time.Duration `mapstruceture:"threshold"`
	Wait           time.Duration `mapstruceture:"wait"`
	RepeatInterval time.Duration `mapstruceture:"repeatInterval""`
}

type Webhook struct {
	Url   string `mapstruceture:"url"`
	Token string `mapstruceture:"token"`
	To    string `mapstruceture:"to"`
	Title string `mapstruceture:"title"`
}

type Db struct {
	Host   string `mapstruceture:"host"`
	Port   int64  `mapstruceture:"port"`
	User   string `mapstruceture:"user"`
	Pass   string `mapstruceture:"pass"`
	Dbname string `mapstruceture:"dbname"`
}

// Init 初始化配置文件
func Init() (err error) {
	viper.SetConfigFile("config.yaml")
	err = viper.ReadInConfig()
	if err != nil {
		log.Println("配置文件读取失败", err)
		return err
	}

	if err := viper.Unmarshal(Conf); err != nil {
		log.Println("配置文件反序列化失败", err)
		return err
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("配置文件改动了")
		if err := viper.Unmarshal(Conf); err != nil {
			log.Println("配置文件反序列化失败", err)
		}
	})
	return

}
