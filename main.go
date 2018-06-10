package main

import (
	"flag"
)

func main() {
	token := flag.String("token", "", "Discord bot token.")
	flag.Parse()

	pepster := NewPepster(*token)
	pepster.Run()
}
