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
	pepster.osuDetailHelper(-1, bid, s, m)
}

func (pepster *Pepster) osuMapsetDetails(sid int, s *discordgo.Session, m *discordgo.MessageCreate) {
	pepster.osuDetailHelper(sid, -1, s, m)
}

func (pepster *Pepster) osuDetailHelper(sid int, bid int, s *discordgo.Session, m *discordgo.MessageCreate) {
	opts := osuapi.GetBeatmapsOpts{}
	if sid != -1 {
		opts.BeatmapSetID = sid
	} else if bid != -1 {
		opts.BeatmapID = bid
	} else {
		// bad function call
		return
	}
	maps, err := pepster.api.GetBeatmaps(opts)
	if err != nil {
		log.Println(err)
		return
	}
	sort.Slice(maps[:], func(i, j int) bool {
		return maps[i].DifficultyRating > maps[j].DifficultyRating
	})

	if len(maps) == 0 {
		return
	}
	firstMap := maps[0]
	description := fmt.Sprintf("Length: %s / BPM: %.1f\n", timeFormat(firstMap.TotalLength), firstMap.BPM)

	diffDetails := make([]string, 0)
	previewLength := 5
	if len(maps) < 5 {
		previewLength = len(maps)
	}
	for _, bmap := range maps[:previewLength] {
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
	description += "\n"
	description += fmt.Sprintf("mapped by [%s](https://osu.ppy.sh/u/%s)", firstMap.Creator, firstMap.Creator)

	embed := discordgo.MessageEmbed{
		URL:         fmt.Sprintf("https://osu.ppy.sh/s/%d", sid),
		Type:        "rich",
		Title:       fmt.Sprintf("%s - %s", firstMap.Artist, firstMap.Title),
		Description: description,
		Color:       mapsetColormap[firstMap.Approved],
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: fmt.Sprintf("https://b.ppy.sh/thumb/%dl.jpg", firstMap.BeatmapSetID),
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
