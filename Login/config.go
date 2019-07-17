package main

// Config ... Holds structure of TOML configuration file
type Config struct {
	General      general
	GameListener gameListener
	Database     database
	Servers      map[string]server
}

type general struct {
	IP        string
	Port      int
	Autologin bool
}

type gameListener struct {
	IP   string
	Port int
}

type database struct {
	Wipe     bool
	Type     byte
	Filename string `toml:"sqlite_filename"`
	User     string
	Password string
	Database string
	IP       string
	Port     int
}

type server struct {
	ID       int
	Name     string
	IP       string
	Port     int
	IsHidden bool `toml:"is_hidden"`
	//IsOnline bool `toml:"is_online"`
	Type   byte
	Color  byte
	Load   byte
	Secret string
}
