package main

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

// Pepster describes an instance of pepster bot
type Pepster struct {
	api *discordgo.Session
}

// NewPepster creates and initializes a new instance of Pepster
func NewPepster(token string) (pepster *Pepster) {
	pepster = new(Pepster)

	dg, err := discordgo.New(token)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	pepster.api = dg

	return
}

// Run is the main function of the bot
func (pepster *Pepster) Run() {
	pepster.login()
}

func (pepster *Pepster) login() {
	err := pepster.api.Open()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Println("connected")
}
