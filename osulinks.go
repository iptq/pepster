package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	osuapi "github.com/thehowl/go-osuapi"
)

var mapsetColormap = map[osuapi.ApprovedStatus]int{
	osuapi.StatusGraveyard: colors["darkgray"],
	osuapi.StatusWIP:       colors["gray"],
	osuapi.StatusPending:   colors["gray"],

	osuapi.StatusQualified: colors["gold"],
	osuapi.StatusLoved:     colors["pink"],
	osuapi.StatusApproved:  colors["lime"],
	osuapi.StatusRanked:    colors["lime"],
}

func (pepster *Pepster) osuMapDetails(bid int, s *discordgo.Session, m *discordgo.MessageCreate) {
	// maps, err := pepster.api.GetBeatmaps(osuapi.GetBeatmapsOpts{BeatmapID: bid})
}

func (pepster *Pepster) osuMapsetDetails(sid int, s *discordgo.Session, m *discordgo.MessageCreate) {
	maps, err := pepster.api.GetBeatmaps(osuapi.GetBeatmapsOpts{BeatmapSetID: sid})
	if err != nil {
		log.Println(err)
		return
	}
	sort.Slice(maps[:], func(i, j int) bool {
		return maps[i].DifficultyRating > maps[j].DifficultyRating
	})

	firstMap := maps[0]
	description := fmt.Sprintf("Length: %s / BPM: %.1f\n", timeFormat(firstMap.TotalLength), firstMap.BPM)

	diffDetails := make([]string, 0)
	for _, bmap := range maps[:5] {
		diffDetails = append(diffDetails, formatHelper(bmap))
	}
	description += strings.Join(diffDetails, "\n")

	if len(maps) > 5 {
		remaining := len(maps) - 5
		var suffix string
		if remaining == 1 {
			suffix = "y"
		} else {
			suffix = "ies"
		}
		description += fmt.Sprintf("\n... %d more difficult%s", remaining, suffix)
	}

	embed := discordgo.MessageEmbed{
		URL:         fmt.Sprintf("https://osu.ppy.sh/s/%d", sid),
		Type:        "rich",
		Title:       fmt.Sprintf("%s - %s", firstMap.Artist, firstMap.Title),
		Description: description,
		Color:       mapsetColormap[firstMap.Approved],
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: fmt.Sprintf("https://b.ppy.sh/thumb/%dl.jpg", sid),
		},
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, &embed)
	if err != nil {
		log.Println(err)
		return
	}
}

func formatHelper(b osuapi.Beatmap) string {
	return fmt.Sprintf("**%s**: %.2f\n CS%.1f / AR%.1f / OD%.1f / HP%.1f", b.DiffName, b.DifficultyRating, b.CircleSize, b.ApproachRate, b.OverallDifficulty, b.HPDrain)
}

func timeFormat(seconds int) string {
	minutes := seconds / 60
	seconds = seconds % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
