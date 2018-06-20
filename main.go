package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/thehowl/conf"
)

// Config describes the pepster config
type Config struct {
	Token    string `description:"Bot token for authentication"`
	APIKey   string `description:"osu! API key"`
	Database string `description:"Path to sqlite database file"`
}

var defaultCfg = Config{
	Database: "pepster.db",
}

func main() {
	configFile := flag.String("conf", "pepster.conf", "config file location")
	flag.Parse()

	config := Config{}
	err := conf.Load(&config, *configFile)
	if err == conf.ErrNoFile {
		conf.Export(defaultCfg, *configFile)
		fmt.Println("Default configuration written to " + *configFile)
		os.Exit(0)
	}

	pepster := NewPepster(config)
	go pepster.Run()

	// wait for signals
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Println("shutting down")
	pepster.Close()
}
