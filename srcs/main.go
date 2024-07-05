package main

import (
	"flag"
	"os"
	"os/signal"
	"sub/api"
	"sub/config"
	"sub/db"
	"sub/log"
	"sub/submitter"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "./config.yml", "Specifies the config file")
	flag.Parse()

	conf, err := config.LoadConfig(configFile)
	if err != nil {
		log.Panic(err)
	}
	api.FirstTick = conf.FirstRound

	log.SetLogFile("submitter.log")
	if conf.LogLevel != "" {
		log.SetLogLevel(conf.LogLevel)
	}
	log.Debugf("%+v\n", conf)

	conf.DB = db.ConnectMongo(conf.Database)
	defer conf.DB.Disconnect()
	err = conf.DB.CreateFlagsCollection()
	if err != nil {
		log.Panic(err)
	}

	go api.ServeAPI(&conf)
	go submitter.Loop(&conf)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Notice("Submitter Running. Press CTRL-C to exit.\n")
	<-stop
}
