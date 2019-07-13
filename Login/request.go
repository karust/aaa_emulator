package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"

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
		return errors.New("Error parsing packet")
	}

	//fmt.Printf("[ChallengeResponse2] pFrom %v, pTo %v, dev %v, str1 %v, login %v, token %v\n", pFrom, pTo, dev, hex.EncodeToString(mac), login, hex.EncodeToString(token))

	user := User{}
	err := loginServer.DB.Select(&user, "SELECT * FROM users WHERE login = ?", login)
	if err != nil {
		fmt.Println("ChallengeResponse2: ", err)
	}
	fmt.Println(user, user == User{})
	// Create new user if autologin
	if loginServer.Autologin {

	}
	//should be proper check for login
	if login != "admin" && login != "test" {
		err := sess.loginDenied("User doesn't exists", 2)
		if err != nil {
			return err
		}
		return nil
	}

	sess.Username = login

	//should be proper check for token/password
	if hex.EncodeToString(token) == "0102030405060708090a0b0c0d0e0f1000000000000000000000000000000000" {
		err := sess.loginDenied("Wrong password, sir", 4)
		if err != nil {
			return err
		}
		return nil
	}

	err = sess.joinResponse()
	if err != nil {
		return err
	}
	err = sess.authResponse()
	if err != nil {
		return err
	}
	return nil
}

// ListWorld ... Request to show game servers
func (sess *Session) ListWorld(parser *packet.Reader) error {
	err := sess.worldListPacket()
	if err != nil {
		return err
	}
	return nil
}

// EnterWorld ... request to enter some server
func (sess *Session) EnterWorld(parser *packet.Reader) error {
	parser.Long()
	serverID := parser.Byte()
	sess.worldCookiePacket(rand.Uint32(), &loginServer.GameServers[serverID])
	return nil
}
