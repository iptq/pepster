package lib

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/vmihailenco/msgpack"
)

var mentionRgx = regexp.MustCompile(`.*(\d+).*`)

type command func([]string, *discordgo.Session, *discordgo.Message)

// Commands is a command manager
type Commands struct {
	pepster *Pepster           // pointer to parent pepster object
	cmdmap  map[string]command // map of commands
}

type MissedMessage struct {
	Timestamp time.Time
	Message   string
}

// NewCommands creates a new instance of the command manager
func NewCommands(pepster *Pepster) (commands Commands) {
	commands = Commands{
		pepster: pepster,
		cmdmap: map[string]command{
			"color":  colorCommand,
			"help":   helpCommand,
			"source": sourceCommand,
			"invite": inviteCommand,
			"tell":   pepster.tellCommand,
		},
	}
	return
}

func successReact(s *discordgo.Session, m *discordgo.Message) {
	err := s.MessageReactionAdd(m.ChannelID, m.ID, "\xf0\x9f\x91\x8d")
	if err != nil {
		log.Println(err)
		return
	}
}

func errorReact(s *discordgo.Session, m *discordgo.Message) {
	err := s.MessageReactionAdd(m.ChannelID, m.ID, "\xf0\x9f\x9a\xab")
	if err != nil {
		log.Println(err)
		return
	}
}

func helpCommand(args []string, s *discordgo.Session, m *discordgo.Message) {
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
			log.Println(err)
		}
	} else {
		msg, err := s.ChannelMessageSend(m.ChannelID, "no")
		if err != nil {
			log.Println(err)
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
}

func sourceCommand(args []string, s *discordgo.Session, m *discordgo.Message) {
	_, err := s.ChannelMessageSend(m.ChannelID, "https://github.com/iptq/pepster")
	if err != nil {
		log.Println(err)
	}
}

func inviteCommand(args []string, s *discordgo.Session, m *discordgo.Message) {
	_, err := s.ChannelMessageSend(m.ChannelID, "https://discord.gg/MpXXvsD")
	if err != nil {
		log.Println(err)
	}
}

func (pepster *Pepster) tellCommand(args []string, s *discordgo.Session, m *discordgo.Message) {
	var usage = "Usage: !tell <mention> <message>"
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, usage)
		return
	}

	mention := args[0]
	message := fmt.Sprintf("%s said: %s", m.Author.Username, strings.Join(args[1:], " "))
	if match := mentionRgx.FindStringSubmatch(mention); match != nil {
		did := m.Mentions[0].ID
		key := fmt.Sprintf("tellmap:%s:%s", m.ChannelID, did)
		v := MissedMessage{Timestamp: time.Now(), Message: message}
		value, err := msgpack.Marshal(v)
		if err != nil {
			log.Println(err)
			return
		}
		pepster.cache.LPush(key, value)
		successReact(s, m)
	} else {
		log.Println("error on", m.Content)
		errorReact(s, m)
	}
}

func colorCommand(args []string, s *discordgo.Session, m *discordgo.Message) {
	var usage = "Usage: !color <name>, where name is from https://www.w3schools.com/colors/colors_names.asp"
	if len(args) != 1 {
		s.ChannelMessageSend(m.ChannelID, usage)
		return
	}
	name := strings.ToLower(args[0])
	value, ok := colors[name]
	if !ok {
		s.ChannelMessageSend(m.ChannelID, usage)
		return
	}

	// get roles
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Println(err)
		return
	}
	roles, err := s.GuildRoles(channel.GuildID)
	if err != nil {
		log.Println(err)
		return
	}

	var colorRole *discordgo.Role
	var colorRoleFound = false
	roleMap := make(map[string]string)
	for _, role := range roles {
		roleMap[role.ID] = role.Name
		if role.Name == "Color: "+name {
			colorRole = role
			colorRoleFound = true
			break
		}
	}

	if !colorRoleFound {
		// create the role
		newRole, err := s.GuildRoleCreate(channel.GuildID)
		if err != nil {
			log.Println(err)
			return
		}
		role, err := s.GuildRoleEdit(channel.GuildID, newRole.ID, "Color: "+name, value, false, 0, false)
		colorRole = role
		if err != nil {
			log.Println(err)
			return
		}
	}

	member, err := s.GuildMember(channel.GuildID, m.Author.ID)
	if err != nil {
		log.Println(err)
		return
	}

	// remove existing colors
	// log.Printf("%+v\n", roleMap)
	for _, roleID := range member.Roles {
		role, ok := roleMap[roleID]
		log.Println(ok, roleID)
		if ok && strings.HasPrefix(role, "Color: ") {
			err := s.GuildMemberRoleRemove(channel.GuildID, m.Author.ID, roleID)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}

	// add current role
	err = s.GuildMemberRoleAdd(channel.GuildID, m.Author.ID, colorRole.ID)
	if err != nil {
		log.Println(err)
		return
	}

	// send emoji reply!
	successReact(s, m)
}
