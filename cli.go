package pepster

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/chzyer/readline"
)

type Cli struct {
	pepster        *Pepster
	currentChannel string
}

func NewCli(pepster *Pepster) (cli *Cli) {
	cli = new(Cli)
	cli.pepster = pepster
	cli.currentChannel = "459128353413136386"
	cli.pepster.dg.AddHandler(cli.messageHandler)
	return
}

func (cli *Cli) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	fmt.Printf("[%s] <%s: %s> %s\n", m.ChannelID, m.Author.Username, m.Author.ID, m.Content)
}

func (cli *Cli) Run() {
	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	cli.pepster.login()
	for {
		inp, err := rl.Readline()
		if err != nil {
			break
		}
		if strings.HasPrefix(inp, "!chan") {
			parts := strings.Split(inp, " ")
			cli.currentChannel = strings.Trim(parts[1], "\n ")
			fmt.Println("now sending to", cli.currentChannel)
		} else {
			cli.pepster.dg.ChannelMessageSend(cli.currentChannel, inp)
			// fmt.Println("sent to", cli.currentChannel, ":", res)
		}
	}
}
