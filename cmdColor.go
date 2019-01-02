package pepster

import (
	"errors"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type ColorCommand struct{}

func (cmd ColorCommand) GetDescription() string {
	return "change your [color](https://www.w3schools.com/colors/colors_names.asp) (ex: `!color DodgerBlue` or `!color none`)"
}

func (cmd ColorCommand) Handle(args []string, s *discordgo.Session, m *discordgo.Message) error {
	var usage = "!color <name>, where name is from https://www.w3schools.com/colors/colors_names.asp"
	if len(args) < 2 {
		return errors.New(usage)
	}
	name := strings.ToLower(args[1])
	value, ok := colors[name]
	if !ok {
		return errors.New("color was not found. pick from https://www.w3schools.com/colors/colors_names.asp")
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
	// log.Printf("%+v\n", member.Roles)
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
