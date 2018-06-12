package main

import (
	"log"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/bwmarrin/discordgo"
)

var (
	queuePage      = "https://old.ppy.sh/forum/60"
	queueChannelID = "442472075278417925"
)

// Queue represents a modding queue
type Queue struct {
	Title       string
	LastUpdated time.Time
}

// QueueMonitor tracks modding queues
func (pepster *Pepster) QueueMonitor() {
	var queues map[string]Queue
	var first = true
	for {
		newQueues, err := fetchQueues()
		if err != nil {
			log.Println(err)
			continue
		}
		if first {
			first = false
			queues = newQueues
			continue
		}
		updates := make(map[string]Queue)
		for h, queue := range newQueues {
			if _, ok := queues[h]; ok && queue.LastUpdated.After(queues[h].LastUpdated) {
				updates[h] = queue
			}
		}
		log.Println(updates)

		embeds := make([]discordgo.MessageEmbed, 0)
		for _, queue := range newQueues {
			embed := discordgo.MessageEmbed{
				Title: queue.Title,
			}
			embeds = append(embeds, embed)
		}
		for _, embed := range embeds {
			_, err := pepster.dg.ChannelMessageSendEmbed(queueChannelID, &embed)
			if err != nil {
				log.Println(err)
			}
		}
		queues = newQueues
		time.Sleep(2 * time.Second)
	}
}

func fetchQueues() (queues map[string]Queue, err error) {
	resp, err := soup.Get(queuePage)
	if err != nil {
		return nil, err
	}
	z := soup.HTMLParse(resp)
	pageContent := z.Find("div", "id", "pagecontent")
	// table := pageContent.Pointer.FirstChild
	log.Println(pageContent)
	// log.Println(resp)
	return
}

func hashQueue() {

}
