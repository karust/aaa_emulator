package main

import "../common/packet"

// LoginDenied ... Shows error :message with a title defined with :reason
/*
   0 = "login_unknown";
   1 = "bad_account";
   2 = "bad_response";
   3 = "duplicate_login";
   4 = "service_time";
   5 = "try_trade_cash_temporal";
   6 = "try_trade_cash_forever";
   7 = "traded_cash_temporal";
   8 = "traded_cash_forever";
   9 = "try_trade_item_servers";
   10 = "traded_item_servers";
   11 = "traded_account";
   12 = "try_cheat_temporal";
   13 = "try_cheat_forever";
   14 = "cheated";
   15 = "gamble_temporal";
   16 = "gamble_forever";
   17 = "abuse_bug_forever";
   18 = "abuse_bug_temporal";
   19 = "use_bot_forever";
   20 = "use_bot_temporal";
   21 = "use_bad_sw_temporal";
   22 = "use_bad_sw_forever";
   23 = "bad_user_workplace";
   24 = "bad_user_proxy_ip";
   25 = "steal_info";
   26 = "foul_lang_temporal";
   27 = "foul_lang_forever";
   28 = "bad_game_name";
   29 = "disturb_play";
   30 = "abnormal_play";
   31 = "disturb_gm";
   32 = "fraudful_report";
   33 = "fake_gm";
   34 = "wait_cert";
   35 = "steal_account_temporal";
   36 = "steal_account_forever";
   37 = "fraudful_steal_report";
   38 = "steal_person";
   39 = "request_by_self";
   40 = "request_by_parent";
   41 = "ads";
   42 = "request_by_authority";
   43 = "defraud_pay";
   44 = "unpaid_account";
   45 = "bulk_blocked_account";
   46 = "unpaid_pcbang";
   47 = "congested_server";
   48 = "invalid_mac";
*/
func (sess *Session) LoginDenied(message string, reason byte) error {
	serial := packet.CreateWriter(12)
	serial.Byte(reason)
	serial.Short(0)
	serial.String(message)
	serial.Send(sess.Client)
	err := serial.Send(sess.Client)
	return err
}

// JoinResponse ...
// TODO: Define params
func (sess *Session) JoinResponse() error {
	serial := packet.CreateWriter(0)
	// serial.Byte(1)       // AuthID
	// serial.Short(0)      // Reason
	// serial.Long(4719366) // "afs" from archerage
	serial.Short(1) // Reason
	serial.Byte(0)
	serial.UInt(0x480306) // afs
	serial.Short(0)
	serial.Byte(0)
	serial.Byte(0) // Slot count
	err := serial.Send(sess.Client)
	return err
}

// AuthResponse ... TODO: What does WSK?
func (sess *Session) AuthResponse(accID uint, slotCount byte) error {
	serial := packet.CreateWriter(3)
	serial.Long(uint64(accID)) // Account ID
	//wsk, _ := utils.RandomHex(16)
	serial.String("65CCBF5AF8DB8B633D3C03C5A8735601") // WSK
	serial.Byte(slotCount)                            // Slot Count
	err := serial.Send(sess.Client)
	return err
}

// WorldListPacket ... Returns information about servers and characters on them
// TODO: Characters info
func (sess *Session) WorldListPacket() error {
	serial := packet.CreateWriter(8)
	serial.Byte(byte(len(loginServer.GameServers)))
	for i := range loginServer.GameServers {
		serial.Byte(loginServer.GameServers[i].ID)
		serial.Byte(loginServer.GameServers[i].Type)
		serial.Byte(loginServer.GameServers[i].Color)
		serial.String(loginServer.GameServers[i].Name)
		serial.Byte(loginServer.GameServers[i].IsOnline)
		serial.Byte(loginServer.GameServers[i].Load)
		serial.Byte(3) // None
		serial.Byte(0) // Nuian
		serial.Byte(3) // Fairy
		serial.Byte(0) // Dwarf
		serial.Byte(0) // Elf
		serial.Byte(0) // Hariharan
		serial.Byte(0) // Ferre
		serial.Byte(3) // Returned
		serial.Byte(0) // Warlocks
	}
	serial.Byte(0) // Char Count
	/*
		serial.Long(character.AccountId)
		serial.Byte(character.GsId)
		serial.UInt(character.Id)
		serial.String(character.Name)
		serial.Byte(character.Race)
		serial.Byte(character.Gender)
		guid, _ := randomHex(8)
		serial.String(guid)
		serial.Long(0) //v
	*/
	err := serial.Send(sess.Client)

	return err
}

// WorldCookiePacket ... Sends IP of chosen game server and cookie to enter
func (sess *Session) WorldCookiePacket(cookie uint32, gameServer *GameServer) error {
	serial := packet.CreateWriter(0xA)
	serial.UInt(cookie)
	//serial.Bytes(gameServer.IP) need reverse it
	serial.Byte(gameServer.IP[3])
	serial.Byte(gameServer.IP[2])
	serial.Byte(gameServer.IP[1])
	serial.Byte(gameServer.IP[0])
	serial.Short(gameServer.Port)
	serial.Long(0)
	serial.Long(0)
	serial.Short(0)

	err := serial.Send(sess.Client)
	return err
}
