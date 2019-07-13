package main

var sqliteSchema = `
DROP TABLE IF EXISTS user;
CREATE TABLE user (
	id UNSIGNED NOT NULL,
	login varchar(32) NOT NULL UNIQUE,
	token text NOT NULL,
	email varchar(128) NOT NULL,
	last_login UNSIGNED NOT NULL DEFAULT '0',
	last_ip varchar(128) NOT NULL,
	created_at UNSIGNED NOT NULL DEFAULT '0',
	updated_at UNSIGNED NOT NULL DEFAULT '0'
);`

var postgesSchema = `
DROP TABLE IF EXISTS users;
CREATE TABLE users (
	id SERIAL NOT NULL,
	login varchar(32) NOT NULL UNIQUE,
	token text NOT NULL,
	email varchar(128) NOT NULL,
	last_login BIGINT NOT NULL DEFAULT '0',
	last_ip varchar(128) NOT NULL,
	created_at BIGINT NOT NULL DEFAULT '0',
	updated_at BIGINT NOT NULL DEFAULT '0'
);`

// User ... User table
type User struct {
	ID          uint
	Login       string
	Token       string
	Email       string
	LastLogin   int64  `db:"last_login"`
	LastIP      string `db:"last_ip"`
	CreatedTime uint   `db:"created_at"`
	UpdatedTime uint   `db:"updated_at"`
}
