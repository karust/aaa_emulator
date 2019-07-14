package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"../common/packet"
)

// ChallengeResponse2 ... Checks login and password information from client
func (sess *Session) ChallengeResponse2(parser *packet.Reader) error {
	parser.Int()   //pFrom := parser.Int()
	parser.Int()   //pTo := parser.Int()
	parser.Bool()  //dev := parser.Bool()
	parser.Bytes() //mac := parser.Bytes()
	login := parser.String()
	token := parser.Bytes()

	if parser.Err {
		return errors.New("Error parsing ChallengeResponse2")
	}

	// Get user from db with provided login
	user := User{}
	stmt, err := loginServer.DB.Preparex(`SELECT * FROM users WHERE login=$1`)
	err = stmt.Get(&user, login)

	// If user not exist and autologin is off
	if err != nil && !loginServer.Autologin {
		sess.LoginDenied("User does not exist", 0)
		// If user not exist and autologin is on
	} else if user.Login == "" && loginServer.Autologin {
		stmt, err = loginServer.DB.Preparex(`INSERT INTO users (login, token, email, last_login, last_ip, created_at, updated_at) VALUES ($1, $2, '', $3, $4, $5, $6);`)
		timeNow := time.Now().Unix()
		clientAddr := sess.Client.RemoteAddr().String()
		_, err = stmt.Exec(login, hex.EncodeToString(token), timeNow, clientAddr, timeNow, timeNow)

	} else if user.Token != hex.EncodeToString(token) {
		sess.LoginDenied("Wrong password, sir", 4)
		err = errors.New("Wrong password")
	}

	// If some error occured during login
	if err != nil {
		fmt.Println("ChallengeResponse2 ", err)
		return err
	}

	sess.Username = login
	err = sess.JoinResponse()
	err = sess.AuthResponse(user.ID, 0)
	return err
}

// ListWorld ... Request to show game servers
func (sess *Session) ListWorld(parser *packet.Reader) error {
	parser.Long() // Flag
	// TODO: send info about characters on servers
	err := sess.WorldListPacket()
	return err
}

// EnterWorld ... request to enter some server
func (sess *Session) EnterWorld(parser *packet.Reader) error {
	parser.Long() // Flag
	gameServerID := parser.Byte()

	if gS, ok := loginServer.GameServers[gameServerID]; ok {
		// TODO: Ask GS enter permission. Get cookie from GS
		//cookie := rand.Uint32()
		sess.WorldCookiePacket(0, gS)
	} else {
		// TODO: ACEnterWorldDeniedPacket
		return errors.New("Chosen server not exist")
	}
	return nil
}
