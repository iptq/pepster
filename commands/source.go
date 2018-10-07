package commands

import (
	"github.com/bwmarrin/discordgo"
)

func SourceCommand(args []string, s *discordgo.Session, m *discordgo.Message) error {
	_, err := s.ChannelMessageSend(m.ChannelID, "https://git.mzhang.me/pepster.git")
	return err
}
