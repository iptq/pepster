package commands

import (
	"errors"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func ColorCommand(args []string, s *discordgo.Session, m *discordgo.Message) error {
	var usage = "Usage: !color <name>, where name is from https://www.w3schools.com/colors/colors_names.asp"
	if len(args) != 1 {
		s.ChannelMessageSend(m.ChannelID, usage)
		return errors.New("Incorrect usage.")
	}
	name := strings.ToLower(args[0])
	value, ok := Colors[name]
	if !ok {
		s.ChannelMessageSend(m.ChannelID, usage)
		return errors.New("Color not found.")
	}

	// get roles
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		return err
	}
	roles, err := s.GuildRoles(channel.GuildID)
	if err != nil {
		return err
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
			return err
		}
		role, err := s.GuildRoleEdit(channel.GuildID, newRole.ID, "Color: "+name, value, false, 0, false)
		colorRole = role
		if err != nil {
			return err
		}
	}

	member, err := s.GuildMember(channel.GuildID, m.Author.ID)
	if err != nil {
		return err
	}

	// remove existing colors
	// log.Printf("%+v\n", roleMap)
	for _, roleID := range member.Roles {
		role, ok := roleMap[roleID]
		log.Println(ok, roleID)
		if ok && strings.HasPrefix(role, "Color: ") {
			err := s.GuildMemberRoleRemove(channel.GuildID, m.Author.ID, roleID)
			if err != nil {
				return err
			}
		}
	}

	// add current role
	err = s.GuildMemberRoleAdd(channel.GuildID, m.Author.ID, colorRole.ID)
	if err != nil {
		return err
	}

	// send emoji reply!
	return successReact(s, m)
}
