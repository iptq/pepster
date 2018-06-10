package main

import (
	"flag"
)

func main() {
	token := *flag.String("token", "", "Discord bot token.")

	pepster := NewPepster(token)
	pepster.Run()
}
