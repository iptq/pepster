package pepster

import (
	"fmt"
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
	conf     Config
	commands Commands
	logger   *log.Logger
	tellMap  map[string][]string
}

// NewPepster creates and initializes a new instance of Pepster
func NewPepster(config Config) (pepster *Pepster) {
	pepster = new(Pepster)

	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal(err)
	}
	dg.AddHandler(pepster.guildCreateHandler)
	pepster.dg = dg

	// TODO: configure in-memory cache
	pepster.cache = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})

	pepster.api = osuapi.NewClient(config.APIKey)
	pepster.tellMap = make(map[string][]string)

	pepster.commands = NewCommands(pepster)
	pepster.logger = NewLogger(pepster)
	pepster.conf = config

	pepster.logger.Println("initialized")
	return
}

func (pepster *Pepster) guildCreateHandler(s *discordgo.Session, g *discordgo.GuildCreate) {
	fmt.Println(g.Name)
	for _, c := range g.Channels {
		fmt.Println(" -", c.ID, c.Type, c.Name)
	}
}

// Run is the main function of the bot
func (pepster *Pepster) Run() {
	pepster.dg.AddHandler(pepster.messageHandler)
	pepster.login()
	go pepster.QueueMonitor()
}

// Cmd runs a cmd
func (pepster *Pepster) Cmd() {
	cli := NewCli(pepster)
	cli.Run()
}

// Close shuts everything down gracefully
func (pepster *Pepster) Close() {
	pepster.logger.Println("closing")
	pepster.dg.Close()
}

func (pepster *Pepster) login() {
	err := pepster.dg.Open()
	if err != nil {
		log.Fatal(err)
	}
	pepster.logger.Println("connected")
}
