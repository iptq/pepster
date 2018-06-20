package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"git.mzhang.me/michael/go-pepster/lib"
	"github.com/thehowl/conf"
)

func main() {
	configFile := flag.String("conf", "pepster.conf", "config file location")
	flag.Parse()

	config := lib.Config{}
	err := conf.Load(&config, *configFile)
	if err == conf.ErrNoFile {
		conf.Export(lib.DefaultCfg, *configFile)
		fmt.Println("Default configuration written to " + *configFile)
		os.Exit(0)
	}

	pepster := lib.NewPepster(config)
	go pepster.Run()

	// wait for signals
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Println("shutting down")
	pepster.Close()
}
