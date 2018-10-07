package pepster

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	osuapi "github.com/thehowl/go-osuapi"
)

type ConfigCommand struct {
	pepster *Pepster
}

func (cmd ConfigCommand) getKey(user string, key string) (string, error) {
	raw, err := cmd.pepster.cache.Get(fmt.Sprintf("config/%s/%s", user, key)).Result()
	if err != nil {
		return "", err
	}
	val := ""
	switch key {
	case "user":
		id, _ := strconv.Atoi(raw)
		user, err := cmd.pepster.api.GetUser(osuapi.GetUserOpts{UserID: id})
		if err != nil {
			return "", err
		}
		val = user.Username
	default:
		return "", fmt.Errorf("`%s` is not a valid key", key)
	}
	return val, nil
}

func (cmd ConfigCommand) setKey(user string, key string, val string) error {
	raw := ""
	switch key {
	case "user":
		user, err := cmd.pepster.api.GetUser(osuapi.GetUserOpts{Username: val})
		if err != nil {
			return err
		}
		raw = fmt.Sprintf("%d", user.UserID)
	default:
		return fmt.Errorf("`%s` is not a valid key", key)
	}
	_, err := cmd.pepster.cache.Set(fmt.Sprintf("config/%s/%s", user, key), raw, 0).Result()
	return err
}

func (cmd ConfigCommand) GetDescription() string {
	return "`!config help` for help"
}

func (cmd ConfigCommand) Handle(args []string, s *discordgo.Session, m *discordgo.Message) error {
	usage := "usage: `!config <key> [<value>]`, type `!config help` for help"
	if len(args) < 2 {
		return errors.New(usage)
	}

	if args[1] == "help" {
		s.ChannelMessageSend(m.ChannelID, "`!config <key>` to check `<key>`. `!config <key> <value>` to set `<key>` to `<value>`\nFor example, `!config user deadcode` to set `user` to `deadcode`.")
		return nil
	}

	if len(args) < 3 {
		// checking
		val, err := cmd.getKey(m.Author.ID, args[1])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>: `%s` not set", m.Author.ID, args[1]))
		} else {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>: `%s` = `%s`", m.Author.ID, args[1], val))
		}
	} else {
		val := strings.Join(args[2:], " ")
		err := cmd.setKey(m.Author.ID, args[1], val)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>: %s", m.Author.ID, err))
		} else {
			successReact(s, m)
		}
	}
	return nil
}
