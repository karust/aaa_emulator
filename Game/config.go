package main

// Config ... Holds structure of TOML configuration file
type Config struct {
	General  general
	Crypto   crypto
	Database database
}

type general struct {
	IP   string
	Port int
}

type crypto struct {
	Modulus  string
	Exponent uint32
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
