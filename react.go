package pepster

import "github.com/bwmarrin/discordgo"

func successReact(s *discordgo.Session, m *discordgo.Message) error {
	return s.MessageReactionAdd(m.ChannelID, m.ID, "\xf0\x9f\x91\x8d")
}

func errorReact(s *discordgo.Session, m *discordgo.Message) error {
	return s.MessageReactionAdd(m.ChannelID, m.ID, "\xf0\x9f\x9a\xab")
}
