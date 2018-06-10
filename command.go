package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type command func([]string, *discordgo.Session, *discordgo.Message)

// Commands is a command manager
type Commands struct {
	pepster *Pepster           // pointer to parent pepster object
	cmdmap  map[string]command // map of commands
}

// NewCommands creates a new instance of the command manager
func NewCommands(pepster *Pepster) (commands Commands) {
	commands = Commands{
		pepster: pepster,
		cmdmap:  make(map[string]command),
	}
	commands.cmdmap["help"] = helpCommand
	return
}

func helpCommand(args []string, s *discordgo.Session, m *discordgo.Message) {
	_, err := s.ChannelMessageSend(m.ChannelID, "no")
	if err != nil {
		log.Println(err)
	}
}
