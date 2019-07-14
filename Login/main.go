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

var loginServer *LoginServer

//
func (login *LoginServer) initializeDatabase(config Config) {
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
		login.DB = db
	} else if config.Database.Type == 1 {

		params := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.IP, c.Port, c.User, c.Password, c.Database)

		db, err = sqlx.Open("postgres", params)
		if err != nil {
			log.Fatalln(err)
		}
		if c.Wipe {
			db.MustExec(postgesSchema)
		}
		login.DB = db
	} else {
		log.Fatalln("Unknown database engine type")
	}
}

// Take config and initialize login server with games servers that it serves
func (login *LoginServer) initializeServer(config Config) {
	login.Address = utils.MakeAdress(config.General.IP, config.General.Port)
	login.Autologin = config.General.Autologin

	login.GameServers = make(map[byte]*GameServer)
	for _, serv := range config.Servers {
		gameServer := GameServer{
			Name:     serv.Name,
			ID:       byte(serv.ID),
			Type:     serv.Type,
			Color:    serv.Color,
			Load:     serv.Load,
			IsOnline: utils.BoolToByte(serv.IsOnline),
			IP:       utils.ConvertIPtoBytes(serv.IP),
			Port:     uint16(serv.Port)}
		login.GameServers[byte(serv.ID)] = &gameServer
	}
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

	loginServer = &LoginServer{}
	loginServer.initializeDatabase(config)
	loginServer.initializeServer(config)
	loginServer.Listen()
}
