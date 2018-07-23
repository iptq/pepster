package lib

import (
	"fmt"
	"log"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	humanize "github.com/dustin/go-humanize"
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

var statusMap = map[osuapi.ApprovedStatus]string{
	osuapi.StatusGraveyard: "Graveyard",
	osuapi.StatusWIP:       "Work in Progress",
	osuapi.StatusPending:   "Pending",

	osuapi.StatusQualified: "Qualified",
	osuapi.StatusLoved:     "Loved",
	osuapi.StatusApproved:  "Approved",
	osuapi.StatusRanked:    "Ranked",
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
	var URL string
	if sid != -1 {
		opts.BeatmapSetID = sid
		URL = fmt.Sprintf("https://osu.ppy.sh/s/%d", sid)
	} else if bid != -1 {
		opts.BeatmapID = bid
		URL = fmt.Sprintf("https://osu.ppy.sh/b/%d", bid)
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

	footer := &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("Status: %s | \u2605 Favorites: %s", statusMap[firstMap.Approved], humanize.Comma(int64(firstMap.FavouriteCount))),
	}

	embed := discordgo.MessageEmbed{
		URL:         URL,
		Type:        "rich",
		Title:       fmt.Sprintf("%s - %s (mapped by %s)", firstMap.Artist, firstMap.Title, firstMap.Creator),
		Description: description,
		Color:       mapsetColormap[firstMap.Approved],
		Fields:      fields,
		Footer:      footer,
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

func (pepster *Pepster) osuUserDetails(uid string, s *discordgo.Session, m *discordgo.MessageCreate) {
	opts := osuapi.GetUserOpts{}
	if id, err := strconv.Atoi(uid); err != nil {
		s, _ := url.PathUnescape(uid)
		opts.Username = s
	} else {
		opts.UserID = id
	}
	user, err := pepster.api.GetUser(opts)
	if err != nil {
		log.Println(err)
		return
	}
	fields := make([]*discordgo.MessageEmbedField, 0)
	description := fmt.Sprintf("\u25b8 **Rank:** #%s  (%spp)   \u25b8 **%s rank:** #%s\n", humanize.Comma(int64(user.Rank)), humanize.Commaf(user.PP), user.Country, humanize.Comma(int64(user.CountryRank)))
	description += "haha sux at osu lol"

	stats := make([]string, 0)
	stats = append(stats, fmt.Sprintf("\u25b8 **Playcount**: %s", humanize.Comma(int64(user.Playcount))))
	stats = append(stats, fmt.Sprintf("\u25b8 **Level**: %.2f", user.Level))
	stats = append(stats, fmt.Sprintf("\u25b8 **Accuracy**: %.2f%%", user.Accuracy))
	overallStats := discordgo.MessageEmbedField{
		Name:  "Overall Stats",
		Value: strings.Join(stats, "\n"),
	}
	fields = append(fields, &overallStats)

	scoreOpts := osuapi.GetUserScoresOpts{
		UserID: user.UserID,
		Limit:  5,
	}
	scores, err := pepster.api.GetUserBest(scoreOpts)
	if err == nil {
		playList := make([]string, 0)
		for _, score := range scores {
			beatmaps, err := pepster.api.GetBeatmaps(osuapi.GetBeatmapsOpts{
				BeatmapID: score.BeatmapID,
			})
			if err != nil || len(beatmaps) < 1 {
				continue
			}
			b := beatmaps[0]
			mods := ""
			if score.Mods > 0 {
				mods = " +" + score.Mods.String()
			}
			playList = append(playList, fmt.Sprintf("\u25b8 %.2fpp\t[%s - %s \\[%s\\]](https://osu.ppy.sh/b/%d)%s", score.PP, b.Artist, b.Title, b.DiffName, score.BeatmapID, mods))
		}
		plays := strings.Join(playList, "\n")
		topPlaysField := discordgo.MessageEmbedField{
			Name:  "Top Plays",
			Value: plays,
		}
		fields = append(fields, &topPlaysField)
	}

	embed := discordgo.MessageEmbed{
		URL:         fmt.Sprintf("https://osu.ppy.sh/u/%d", user.UserID),
		Type:        "rich",
		Title:       fmt.Sprintf(":flag_%s: %s", strings.ToLower(user.Country), user.Username),
		Description: description,
		Fields:      fields,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: fmt.Sprintf("https://a.ppy.sh/%d", user.UserID),
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
	line := fmt.Sprintf("   \u25b8 **Difficulty:** %.2f\u2605   \u25b8 **Max Combo:** %dx\n   \u25b8 **CS:** %.1f \u25b8 **AR:** %.1f \u25b8 **OD:** %.1f \u25b8 **HP:** %.1f", b.DifficultyRating, b.MaxCombo, b.CircleSize, b.ApproachRate, b.OverallDifficulty, b.HPDrain)

	return &discordgo.MessageEmbedField{
		Name:  fmt.Sprintf("%s __%s__", diffEmoji[emoji], b.DiffName),
		Value: line,
	}
}

func timeFormat(seconds int) string {
	minutes := seconds / 60
	seconds = seconds % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
