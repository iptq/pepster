package pepster

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	mapPattern    = regexp.MustCompile(`.*https?://(osu|old).ppy.sh/b/(?P<id>\d+)[^/]*`)
	mapsetPattern = regexp.MustCompile(`.*https?://(osu|old).ppy.sh/(s|beatmapsets)/(?P<id>\d+)[^/]*`)
	userPattern   = regexp.MustCompile(`.*https?://(osu|old).ppy.sh/(u|users)/([^/]+).*`)
)

func (pepster *Pepster) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore bot messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	// commands
	if strings.HasPrefix(m.Content, "!") {
		line := strings.TrimLeft(m.Content, "!")

		// TODO: some kind of quote parser for this
		argv := strings.Split(line, " ")
		if len(argv) == 0 {
			return
		}

		// special help command
		if argv[0] == "help" {
			help := pepster.commands.GenerateHelp()
			s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
				Description: help,
			})
		}

		// switch on command
		command, ok := pepster.commands.Get(argv[0])
		if !ok {
			// just don't do anything
			return
		} else {
			err := command.Handle(argv, s, m.Message)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> error: %s", m.Author.ID, err))
			}
		}
	}

	// osu link handlers
	if match := mapPattern.FindStringSubmatch(m.Content); match != nil {
		bid, err := strconv.Atoi(match[2])
		if err == nil {
			pepster.osuMapDetails(bid, s, m)
		}
	}
	if match := mapsetPattern.FindStringSubmatch(m.Content); match != nil {
		sid, err := strconv.Atoi(match[3])
		if err == nil {
			pepster.osuMapsetDetails(sid, s, m)
		}
	}
	if match := userPattern.FindStringSubmatch(m.Content); match != nil {
		uid := match[3]
		pepster.osuUserDetails(uid, s, m)
	}
}
