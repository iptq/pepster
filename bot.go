package pepster

import (
	"fmt"
	"log"

	"pepster/models"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	osuapi "github.com/thehowl/go-osuapi"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// Pepster describes an instance of pepster bot
type Pepster struct {
	cache    *redis.Client
	dg       *discordgo.Session
	db       *xorm.Engine
	api      *osuapi.Client
	conf     Config
	commands Commands
	logger   *log.Logger
	tellMap  map[string][]string
}

// NewPepster creates and initializes a new instance of Pepster
func NewPepster(config Config) (pepster *Pepster, err error) {
	pepster = new(Pepster)

	// set up discord client
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return nil, err
	}
	pepster.dg = dg

	// set up database
	engine, err := xorm.NewEngine(config.DatabaseProvider, config.Database)
	if err != nil {
		log.Fatal(err)
	}
	engine.SetMapper(core.GonicMapper{})
	err = engine.Sync(new(models.User))
	if err != nil {
		log.Fatal(err)
	}
	pepster.db = engine

	fmt.Println("%+v", config)

	// TODO: configure in-memory cache
	pepster.cache = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		DB:       config.RedisDB,
		Password: config.RedisPassword,
	})

	pepster.api = osuapi.NewClient(config.APIKey)
	pepster.tellMap = make(map[string][]string)

	pepster.commands = NewCommands(pepster)
	pepster.commands.Register("color", ColorCommand{})
	pepster.commands.Register("devchan", DevChanCommand{})
	pepster.commands.Register("source", SourceCommand{})
	pepster.commands.Register("last", LastCommand{pepster})
	pepster.commands.Register("config", ConfigCommand{pepster})

	pepster.logger = NewLogger(pepster)
	pepster.conf = config

	pepster.logger.Println("initialized")
	return pepster, nil
}

// Run is the main function of the bot
func (pepster *Pepster) Run() {
	pepster.dg.AddHandler(pepster.messageHandler)
	pepster.login()
	// go pepster.QueueMonitor()
}

// Cmd runs a cmd
func (pepster *Pepster) Cmd() {
	cli := NewCli(pepster)
	pepster.dg.AddHandler(pepster.guildCreateHandler)
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
