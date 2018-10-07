package pepster

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Handle([]string, *discordgo.Session, *discordgo.Message) error
	GetDescription() string
}

type Commands struct {
	pepster  *Pepster
	commands map[string]Command
}

func NewCommands(pepster *Pepster) Commands {
	commands := make(map[string]Command)
	return Commands{
		pepster,
		commands,
	}
}

func (cmd *Commands) Get(key string) (Command, bool) {
	val, ok := cmd.commands[key]
	return val, ok
}

func (cmd *Commands) Register(key string, val Command) {
	cmd.commands[key] = val
}

func (cmd *Commands) GenerateHelp() string {
	lines := make([]string, 0)
	lines = append(lines, "**!help**: display this help")
	for name, command := range cmd.commands {
		lines = append(lines, fmt.Sprintf("**!%s**: %s", name, command.GetDescription()))
	}
	return strings.Join(lines, "\n")
}
