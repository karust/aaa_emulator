package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"../common/utils"
	"github.com/BurntSushi/toml"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var gameServer *GameServer

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
	game.Address = utils.MakeAddress(config.General.IP, config.General.Port)

	// Convert modulus
	mod, err := hex.DecodeString(config.Crypto.Modulus)
	if err != nil {
		log.Fatalln("Wrong modulus of public key")
	}
	game.PubModulus = mod

	// Convert exponent
	byteModulus := make([]byte, 4)
	binary.BigEndian.PutUint32(byteModulus, config.Crypto.Exponent)
	game.PubExponent = append(make([]byte, 124), byteModulus...)
}

func main() {
	args, configPath := os.Args, "config.toml"
	if len(args) >= 2 {
		configPath = args[1]
	}

	// Try to load configuration file, if error then meaningless to proceed
	log.Println("Loading configuration file...")
	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		log.Fatalln("Config load error:", configPath)
	}

	// Connect to Login server
	loginConn := LoginConnection{}
	if err := loginConn.Initialize(config); err != nil {
		log.Fatalln("Cannot establish connection with Login server, check Address!")
	}
	go loginConn.Listen()

	// Initialize Database and Game Server
	gameServer = &GameServer{}
	gameServer.initializeDatabase(config)
	gameServer.initializeServer(config)
	gameServer.Listen()
}
