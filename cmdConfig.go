package pepster

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type ConfigCommand struct {
	pepster *Pepster
}

func (cmd ConfigCommand) GetDescription() string {
	return "`!config help` for help"
}

func (cmd ConfigCommand) Handle(args []string, s *discordgo.Session, m *discordgo.Message) error {
	usage := "usage: `!config <key> [<value>]`, type `!config help` for help"
	if len(args) < 2 {
		return errors.New(usage)
	}

	switch args[1] {
	case "help":
		s.ChannelMessageSend(m.ChannelID, "`!config <key>` to check `<key>`. `!config <key> <value>` to set `<key>` to `<value>`\nFor example, `!config user deadcode` to set `user` to `deadcode`.")
	case "user":
		fmt.Println("shiet")
	}
	return nil
}
