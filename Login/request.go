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
	sess.AccountID = user.ID
	sess.ID = uint32(loginServer.NumConnections & 0xffffffff)
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

// X2EnterWorld ... request to enter some server
func (sess *Session) X2EnterWorld(parser *packet.Reader) error {
	parser.Long() // Flag
	gameServerID := parser.Byte()

	// RequestEnterWorld
	// If there is game server with such ID
	var gS *GameServer
	var ok bool
	if gS, ok = loginServer.GameServers[gameServerID]; !ok {
		return errors.New("Server not exist")
	}

	if gS.IsOnline == true {
		loginServer.Clients.Set(sess.ID, sess)
		loginServer.GameConn.lgPlayerEnter(sess.AccountID, sess.ID, gS.ID)
		return nil
	}
	return errors.New("Server not active")
}

// RequestReconnect ...
func (sess *Session) RequestReconnect(reader *packet.Reader) error {
	pFrom := reader.UInt()
	pTo := reader.UInt()
	accountID := reader.Long()
	gsID := reader.Byte()
	cookie := reader.UInt()
	macLength := reader.Short()
	mac := reader.BytesLen(macLength)

	sess.JoinResponse() //1, 0x480306, 0)
	sess.AuthResponse(sess.AccountID, 0)
	fmt.Println("RequestReconnect:", pFrom, pTo, accountID, gsID, cookie, mac)
	return nil
}
