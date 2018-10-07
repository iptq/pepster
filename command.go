package pepster

import (
	"regexp"
	"time"

	cmd "pepster/commands"

	"github.com/bwmarrin/discordgo"
)

var mentionRgx = regexp.MustCompile(`.*(\d+).*`)

type command func([]string, *discordgo.Session, *discordgo.Message) error

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
			"color":   cmd.ColorCommand,
			"help":    cmd.HelpCommand,
			"source":  cmd.SourceCommand,
			"invite":  cmd.DevChanCommand,
			"last":    cmd.LastCommand,
			"compare": cmd.CompareCommand,
			"get":     cmd.GetCommand,
			"set":     cmd.SetCommand,
		},
	}
	return
}
