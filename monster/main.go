package main

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"monster/db"
	"monster/web"

	_ "net/http/pprof"
)

func main() {
	viper.SetConfigType("toml")
	viper.SetConfigFile("./configs/config.toml")
	viper.ReadInConfig()

	log.SetFormatter(&log.TextFormatter{DisableColors: true})
	if viper.GetBool("core.debug") {
		log.SetLevel(log.DebugLevel)
	}
	logger := log.WithField("source", "main")

	dbIns := db.NewDBModule()
	err := dbIns.Init()
	if err != nil {
		logger.Fatalf("init db module err:%v", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		web.Startup()
	}() // 开启web服务器

	wg.Wait()
}
