package pepster

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
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
	cli.pepster.login()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">>> ")
		inp, _ := reader.ReadString('\n')
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
