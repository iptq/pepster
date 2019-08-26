package pepster

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type PlayerCard struct {
	ID            int
	Username      string
	Country       string
	Mode          int
	GlobalRank    int
	CountryRank   int
	Level         int
	LevelProgress float64
	PP            float64
	Accuracy      float64
	PlayCount     int
}

func (pepster *Pepster) PrintPlayerCard(s *discordgo.Session, targetChannel string, playerCard *PlayerCard) (err error) {
	var modeString string
	switch playerCard.Mode {
	case 0:
		modeString = "osu! Standard"
	case 1:
		modeString = "Taiko"
	case 2:
		modeString = "Catch the Beat"
	case 3:
		modeString = "Mania"
	}

	embed := discordgo.MessageEmbed{
		Type:  "rich",
		URL:   fmt.Sprintf("https://osu.ppy.sh/u/%d", playerCard.ID),
		Title: fmt.Sprintf("%s Profile for %s", modeString, playerCard.Username),
	}

	_, err = s.ChannelMessageSendEmbed(targetChannel, &embed)
	return err
}
