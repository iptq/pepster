package lib

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	osuapi "github.com/thehowl/go-osuapi"
)

// Pepster describes an instance of pepster bot
type Pepster struct {
	cache    *redis.Client
	dg       *discordgo.Session
	api      *osuapi.Client
	commands Commands
	tellMap  map[string][]string
}

// NewPepster creates and initializes a new instance of Pepster
func NewPepster(config Config) (pepster *Pepster) {
	pepster = new(Pepster)

	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal(err)
	}
	pepster.dg = dg

	// TODO: configure in-memory cache
	pepster.cache = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})

	pepster.api = osuapi.NewClient(config.APIKey)
	pepster.tellMap = make(map[string][]string)

	pepster.commands = NewCommands(pepster)
	return
}

// Run is the main function of the bot
func (pepster *Pepster) Run() {
	pepster.dg.AddHandler(pepster.messageHandler)
	pepster.login()
	go pepster.QueueMonitor()
}

// Close shuts everything down gracefully
func (pepster *Pepster) Close() {
	pepster.dg.Close()
}

func (pepster *Pepster) login() {
	err := pepster.dg.Open()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected")
}
