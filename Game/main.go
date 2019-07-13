package main

import (
	"fmt"
	"log"
	"os"

	"../common/utils"
	"github.com/BurntSushi/toml"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func (game *GameServer) initializeDatabase(config Config) {
	var err error
	var db *sqlx.DB
	c := config.Database

	if config.Database.Type == 0 {
		db, err = sqlx.Connect("sqlite3", c.Filename)
		if err != nil {
			log.Fatalln(err)
		}
		if c.Wipe {
			db.MustExec(sqliteSchema)
		}
		game.DB = db
	} else if config.Database.Type == 1 {

		params := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.IP, c.Port, c.User, c.Password, c.Database)

		db, err = sqlx.Open("postgres", params)
		if err != nil {
			log.Fatalln(err)
		}
		if c.Wipe {
			db.MustExec(postgesSchema)
		}
		game.DB = db
	} else {
		log.Fatalln("Unknown database engine type")
	}
}

func (game *GameServer) initializeServer(config Config) {
	game.Address = utils.MakeAdress(config.General.IP, config.General.Port)
}

func main() {
	args, configPath := os.Args, "config.toml"
	if len(args) >= 2 {
		configPath = args[1]
	}

	// Try to load configuration file, if error then meaningless to proceed
	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		log.Fatalln("Config load error:")
	}

	gameServer := &GameServer{}
	gameServer.initializeDatabase(config)
	gameServer.initializeServer(config)
	//gameServer.Listen()
}
