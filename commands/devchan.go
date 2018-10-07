package commands

import (
	"github.com/bwmarrin/discordgo"
)

func DevChanCommand(args []string, s *discordgo.Session, m *discordgo.Message) error {
	_, err := s.ChannelMessageSend(m.ChannelID, "https://discord.gg/MpXXvsD")
	return err
}
