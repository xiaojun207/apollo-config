package config

import (
	"fmt"
	"github.com/shima-park/agollo"
	"log"
)

type Conf interface {
	agollo.Agollo
}

var LogLevel = "debug"

func InitConfig(configServerURL, appID string, loadConfig func(Config Conf)) agollo.Agollo {

	log.Println("appID:", appID)
	apolloConfig, err := agollo.New(configServerURL, appID, agollo.AutoFetchOnCacheMiss())
	if err != nil {
		panic(err)
	}

	go LoadApolloWatch(apolloConfig, loadConfig)

	return apolloConfig
}

func LoadApolloWatch(Config agollo.Agollo, loadConfig func(Config Conf)) {
	loadConfig(Config)
	errorCh := Config.Start()

	// 监听apollo配置更改事件
	// 返回namespace和其变化前后的配置,以及可能出现的error
	watchCh := Config.Watch()

	stop := make(chan bool)
	watchNamespace := "application"
	watchNSCh := Config.WatchNamespace(watchNamespace, stop)

	go func() {
		for {
			select {
			case err := <-errorCh:
				fmt.Println("Error:", err)
			case <-watchCh:
				loadConfig(Config)
			case <-watchNSCh:
				//fmt.Println("Watch Namespace", watchNamespace, resp)
			}
		}
	}()

	select {}
}

func GetConf(Config Conf, key, defaultValue string) string {
	return Config.Get(key, agollo.WithDefault(defaultValue))
}
