package pepster

import "github.com/bwmarrin/discordgo"

type DevChanCommand struct{}

func (cmd DevChanCommand) GetDescription() string {
	return "Get an invite link to the dev discord for this bot."
}

func (cmd DevChanCommand) Handle(args []string, s *discordgo.Session, m *discordgo.Message) error {
	_, err := s.ChannelMessageSend(m.ChannelID, "https://discord.gg/MpXXvsD")
	return err
}

type SourceCommand struct{}

func (cmd SourceCommand) GetDescription() string {
	return "Get a link to the source code of this bot."
}

func (cmd SourceCommand) Handle(args []string, s *discordgo.Session, m *discordgo.Message) error {
	_, err := s.ChannelMessageSend(m.ChannelID, "https://git.mzhang.me/pepster.git")
	return err
}
