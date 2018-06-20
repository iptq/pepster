package main

import (
	"fmt"
	"log"
	"sort"

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

var diffEmoji = map[string]string{
	"easy":   "<:easy:459040242050007042>",
	"normal": "<:normal:459040207031894037>",
	"hard":   "<:hard:459040169472032788>",
	"insane": "<:insane:459040134181158933>",
	"expert": "<:expert:459038617688473600>",
}

var maxDiffsToShow = 4

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
	description := fmt.Sprintf("\\\U0001f558 **Length:** %s - \U0001d160 **BPM:** %.1f\n", timeFormat(firstMap.TotalLength), firstMap.BPM)

	previewLength := maxDiffsToShow
	if len(maps) < maxDiffsToShow {
		previewLength = len(maps)
	}
	fields := make([]*discordgo.MessageEmbedField, previewLength)
	for i, bmap := range maps[:previewLength] {
		fields[i] = formatHelper(bmap)
		// diffDetails = append(diffDetails, formatHelper(bmap))
	}
	// description += strings.Join(diffDetails, "\n")

	if len(maps) > maxDiffsToShow {
		remaining := len(maps) - maxDiffsToShow
		var suffix string
		if remaining == 1 {
			suffix = "y"
		} else {
			suffix = "ies"
		}
		description += fmt.Sprintf("(%d difficult%s not shown)", remaining, suffix)
	}

	embed := discordgo.MessageEmbed{
		URL:         fmt.Sprintf("https://osu.ppy.sh/s/%d", sid),
		Type:        "rich",
		Title:       fmt.Sprintf("%s - %s (mapped by %s)", firstMap.Artist, firstMap.Title, firstMap.Creator),
		Description: description,
		Color:       mapsetColormap[firstMap.Approved],
		Fields:      fields,
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

func formatHelper(b osuapi.Beatmap) *discordgo.MessageEmbedField {
	var emoji string
	switch {
	case b.DifficultyRating < 1.5:
		emoji = "easy"
	case b.DifficultyRating < 2.25:
		emoji = "normal"
	case b.DifficultyRating < 3.75:
		emoji = "hard"
	case b.DifficultyRating < 5.25:
		emoji = "insane"
	default:
		emoji = "expert"
	}
	line := fmt.Sprintf("**Difficulty:** %.2f - **Max Combo:** %dx\n CS%.1f / AR%.1f / OD%.1f / HP%.1f", b.DifficultyRating, b.MaxCombo, b.CircleSize, b.ApproachRate, b.OverallDifficulty, b.HPDrain)

	return &discordgo.MessageEmbedField{
		Name:  fmt.Sprintf("%s %s", diffEmoji[emoji], b.DiffName),
		Value: line,
	}
}

func timeFormat(seconds int) string {
	minutes := seconds / 60
	seconds = seconds % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
