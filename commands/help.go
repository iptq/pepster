package commands

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func HelpCommand(args []string, s *discordgo.Session, m *discordgo.Message) error {
	if len(args) > 0 && strings.Contains(strings.ToLower(args[0]), "please") {
		description := strings.Join([]string{
			"`!color <color>` => Change your color",
			"`!help` => Help contents",
			"`!tell <user> <message>` => Tell user message",
		}, "\n")
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title:       "pepster bot",
			Description: description,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "ur dad lesbian",
			},
		})
		if err != nil {
			return err
		}
	} else {
		msg, err := s.ChannelMessageSend(m.ChannelID, "no")
		if err != nil {
			return err
		}
		go (func() {
			time.Sleep(1 * time.Second)
			s.ChannelMessageEdit(m.ChannelID, msg.ID, "lol")
			go (func() {
				time.Sleep(2 * time.Second)
				s.ChannelMessageEdit(m.ChannelID, msg.ID, "fuck u")
			})()
		})()
	}
	return nil
}
