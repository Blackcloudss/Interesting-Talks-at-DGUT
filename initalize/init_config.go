package initalize

import (
	"flag"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/configs"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/global"
	"github.com/Blackcloudss/Interesting-Talks-at-DGUT/log/zlog"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"time"
)

// @Title        init_config.go
// @Description
// @Create       XdpCs 2025-02-24 下午11:15
// @Update       XdpCs 2025-02-24 下午11:15
func InitConfig() {
	// 初始化时间为东八区的时间
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone

	// 默认配置文件路径
	var configPath string
	flag.StringVar(&configPath, "c", global.Path+global.DEFAULT_CONFIG_FILE_PATH, "配置文件绝对路径或相对路径")
	flag.Parse()
	zlog.Infof("配置文件路径为 %s", configPath)
	// 初始化配置文件
	viper.SetConfigFile(configPath)
	viper.WatchConfig()
	// 观察配置文件变动
	viper.OnConfigChange(func(in fsnotify.Event) {
		zlog.Warnf("配置文件发生变化")
		if err := viper.Unmarshal(&configs.Conf); err != nil {
			zlog.Errorf("无法反序列化配置文件 %v", err)
		}
		zlog.Debugf("%+v", configs.Conf)
		global.Config = configs.Conf
		Eve()
		Init()
	})
	// 将配置文件读入 viper
	if err := viper.ReadInConfig(); err != nil {
		zlog.Panicf("无法读取配置文件 err: %v", err)

	}
	// 解析到变量中
	if err := viper.Unmarshal(&configs.Conf); err != nil {
		zlog.Panicf("无法解析配置文件 err: %v", err)
	}
	zlog.Debugf("配置文件为 ： %+v", configs.Conf)
	global.Config = configs.Conf
}
