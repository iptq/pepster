package pepster

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	osuapi "github.com/thehowl/go-osuapi"
)

type ProfileCommand struct {
	pepster *Pepster
}

func (cmd ProfileCommand) GetDescription() string {
	return "Get the information about the last play you made."
}

func (cmd ProfileCommand) Handle(args []string, s *discordgo.Session, m *discordgo.Message) error {
	key := fmt.Sprintf("config/%s/user", m.Author.ID)
	val, err := cmd.pepster.cache.Get(key).Result()
	if err != nil {
		return fmt.Errorf("I don't know who you are! Set your username with `!config user <username>`")
	}

	id, _ := strconv.Atoi(val)
	scores, err := cmd.pepster.api.GetUserRecent(osuapi.GetUserScoresOpts{UserID: id})
	if err != nil || len(scores) == 0 {
		return errorReact(s, m)
	}

	first := scores[0]
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s", first.Score),
		Description: "hi",
	})
	return nil
}
