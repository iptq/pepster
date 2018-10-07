package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"pepster"

	"github.com/thehowl/conf"
)

func main() {
	configFile := flag.String("conf", "pepster.conf", "config file location")
	cmd := flag.Bool("cmd", false, "whether to start as a command line")
	flag.Parse()

	config := pepster.Config{}
	err := conf.Load(&config, *configFile)
	if err == conf.ErrNoFile {
		conf.Export(pepster.DefaultCfg, *configFile)
		fmt.Println("Default configuration written to " + *configFile)
		os.Exit(0)
	}

	pepster, _ := pepster.NewPepster(config)
	defer pepster.Close()
	if *cmd {
		pepster.Cmd()
	} else {
		go pepster.Run()

		// wait for signals
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
	}
}
