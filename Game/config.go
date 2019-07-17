package main

// Config ... Holds structure of TOML configuration file
type Config struct {
	General  general
	Login    login
	Crypto   crypto
	Database database
}

type general struct {
	IP   string
	Port int
}

type login struct {
	IP     string
	Port   int
	Secret string
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
