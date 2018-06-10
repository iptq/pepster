package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (pepster *Pepster) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore bot messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!") {
		// it's a command
		parts := strings.Split(m.Content[1:], " ")
		fn, ok := pepster.commands.cmdmap[parts[0]]
		if ok {
			fn(parts[1:], s, m.Message)
		}
	}
}
