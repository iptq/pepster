package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	token := flag.String("token", "", "Discord bot token.")
	flag.Parse()

	pepster := NewPepster("Bot " + *token)
	go pepster.Run()

	// wait for signals
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Println("shutting down")
	pepster.Close()
}
