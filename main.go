package main

import (
	"flag"
	"github.com/jasonkayzk/distributed-id-generator/config"
	"github.com/jasonkayzk/distributed-id-generator/core"
	"github.com/jasonkayzk/distributed-id-generator/mysql"
	"github.com/jasonkayzk/distributed-id-generator/server"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "config.json", "the config file: config.json")
	flag.Parse()

	err := config.LoadConf(configFile)
	if err != nil {
		panic(err)
	}
	err = mysql.InitDB()
	if err != nil {
		panic(err)
	}
	err = core.InitAlloc()
	if err != nil {
		panic(err)
	}
	err = server.StartServer()
	if err != nil {
		panic(err)
	}
}
