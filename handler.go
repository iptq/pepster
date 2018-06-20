package main

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	mapPattern    = regexp.MustCompile(`.*https?://osu.ppy.sh/b/(\d+)[^/]*`)
	mapsetPattern = regexp.MustCompile(`.*https?://osu.ppy.sh/(s|beatmapsets)/(\d+)[^/]*`)
)

func (pepster *Pepster) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore bot messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	// commands
	if strings.HasPrefix(m.Content, "!") {
		// it's a command
		parts := strings.Split(m.Content[1:], " ")
		fn, ok := pepster.commands.cmdmap[parts[0]]
		if ok {
			fn(parts[1:], s, m.Message)
		}
		return
	}

	// osu link handlers
	if match := mapPattern.FindStringSubmatch(m.Content); match != nil {
		bid, err := strconv.Atoi(match[1])
		if err == nil {
			pepster.osuMapDetails(bid, s, m)
		}
	}
	if match := mapsetPattern.FindStringSubmatch(m.Content); match != nil {
		sid, err := strconv.Atoi(match[2])
		if err == nil {
			pepster.osuMapsetDetails(sid, s, m)
		}
	}
}
