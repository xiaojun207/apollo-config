package config

import (
	"fmt"
	"github.com/shima-park/agollo"
	"log"
)

type ApolloConfig struct {
	agollo.Agollo
}

var LogLevel = "debug"

func InitConfig(configServerURL, appID string, loadConfig func(conf ApolloConfig)) agollo.Agollo {

	log.Println("appID:", appID)
	apolloConfig, err := agollo.New(configServerURL, appID, agollo.AutoFetchOnCacheMiss())
	if err != nil {
		panic(err)
	}

	go LoadApolloWatch(apolloConfig, loadConfig)

	return apolloConfig
}

func LoadApolloWatch(apolloConfig agollo.Agollo, loadConfig func(conf ApolloConfig)) {
	var conf = ApolloConfig{apolloConfig}
	loadConfig(conf)
	errorCh := conf.Start()

	// 监听apollo配置更改事件
	// 返回namespace和其变化前后的配置,以及可能出现的error
	watchCh := conf.Watch()

	stop := make(chan bool)
	watchNamespace := "application"
	watchNSCh := conf.WatchNamespace(watchNamespace, stop)

	go func() {
		for {
			select {
			case err := <-errorCh:
				fmt.Println("Error:", err)
			case <-watchCh:
				loadConfig(conf)
			case <-watchNSCh:
				//fmt.Println("Watch Namespace", watchNamespace, resp)
			}
		}
	}()

	select {}
}

func (e *ApolloConfig) Get(key, defaultValue string) string {
	return e.Agollo.Get(key, agollo.WithDefault(defaultValue))
}
