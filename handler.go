package pepster

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	humanize "github.com/dustin/go-humanize"
	"github.com/vmihailenco/msgpack"
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

	// now check for new messages
	go func() {
		key := fmt.Sprintf("tellmap:%s:%s", m.ChannelID, m.Author.ID)
		msgs := pepster.cache.LRange(key, 0, -1).Val()
		summary := fmt.Sprintf("<@%s>: while you were gone, you missed these messages:\n", m.Author.ID)
		for _, msgencode := range msgs {
			var msg MissedMessage
			err := msgpack.Unmarshal([]byte(msgencode), &msg)
			if err != nil {
				log.Println(err)
				continue
			}
			summary += fmt.Sprintf("%s: %s\n", humanize.Time(msg.Timestamp), msg.Message)
		}
		if len(msgs) > 0 {
			pepster.tellMap[m.Author.ID] = nil
			_, err := s.ChannelMessageSend(m.ChannelID, summary)
			if err == nil {
				pepster.cache.Del(key)
			}
		}
	}()
}
