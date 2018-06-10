package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// Pepster describes an instance of pepster bot
type Pepster struct {
	api      *discordgo.Session
	commands Commands
}

// NewPepster creates and initializes a new instance of Pepster
func NewPepster(token string) (pepster *Pepster) {
	pepster = new(Pepster)

	dg, err := discordgo.New(token)
	if err != nil {
		log.Fatal(err)
	}
	pepster.api = dg

	pepster.commands = NewCommands(pepster)
	return
}

// Run is the main function of the bot
func (pepster *Pepster) Run() {
	// handlers
	pepster.api.AddHandler(pepster.messageHandler)

	pepster.login()
}

// Close shuts everything down gracefully
func (pepster *Pepster) Close() {
	pepster.api.Close()
}

func (pepster *Pepster) login() {
	err := pepster.api.Open()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected")
}
