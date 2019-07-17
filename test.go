package main

import (
	"encoding/hex"
	"fmt"

	"./common/crypt"
)

func main() {
	/*
		1000dd05 8d37a9c797704010e0b0815186b6
		0800dd05 12a1de501140
	*/

	// AES:  fab67330c36ccc117425a4831a6eb2fc , XOR:  2756154973
	aesKey, _ := hex.DecodeString("fab67330c36ccc117425a4831a6eb2fc")
	cr := crypt.ClientCrypt(aesKey, 0x2756154973)

	packet, _ := hex.DecodeString("30648929c79e2d755756108b3c0b9fea074837841c8cb7fde128000b10d13bb8e4267af2b6e7da2cacf225e7b6202144c042d1bb5bf319018d239036a98ec53378cbcd1db17e193f07c7dce3fb5ee53cb0b0872653a7b691ad248cb45796b5b9b4bd6c516ef95f1d625dc1007f51d46ecb65045bcf3d357e07dfcf316c1356d559bf4fc761b51065b3f746a7cae55ed80d630817004d2c24c90ad64dc35623a932f0cf288d17f434148e5cc3deebf81046313262d06dd17d8ea14342d52e9e2909af7ab0a307f640197c7502ce4e0a069ffb5ae0190a02a7b529adf43dacf566043dd61a14d87c3fec55fd76f2e513c91e628057786c9c1e9afd2df59c51f631be5147a57c2309e9321815f6de1a8d3579c544918269bcab35d1c7b6a61e5fb1e2b7bdc23c516d39ae698806bb7204dd5044f4940cb179011e6a874a1c27d54d4b")
	decr := cr.Decrypt(packet, len(packet))
	fmt.Println(decr)
	//pack1, _ := hex.DecodeString("12a1de501140")
	//dec1 := crypt.ToClientEncr(pack1)
	//fmt.Println(string(dec1))
	//fmt.Println(hex.EncodeToString(dec1))

	// pack2, _ := hex.DecodeString("cd2089c382625271b0cb11257381d3e43840dba6f90b2c06a99043486eefcd4b6745f31f35e70d0fd0441d85e93e1d834c23f4c494643e05d5a5754517e7b7875627f7c797704010e0b0115081b0")
	// dec2 := crypt.ToClientEncr(pack2)
	// fmt.Println(string(dec2))
	// fmt.Println(hex.EncodeToString(dec2))
	// return

	/*
		aesk, _ := hex.DecodeString("77a89d5c1df25fb2e44477f7e4b1752a")
		cryptAes := crypt.ClientCrypt(aesk, 2401647235)

		pack1, _ := hex.DecodeString("394e9cc2791e5c13bb55415ea6b374df29")
		dec = cryptAes.Decrypt(pack1, 19)
		fmt.Println(hex.EncodeToString(dec))

		f, _ := os.Create("Game/etc/big_bad_dec")
		defer f.Close()

		var seq byte = 0x17

		data, _ := ioutil.ReadFile("Game/etc/big_bad2")
		offset := 0
		for offset < len(data) {
			pkLen := int(binary.LittleEndian.Uint16(data[offset : offset+2]))
			offset += 2
			_type := binary.LittleEndian.Uint16(data[offset : offset+2])
			offset += 2
			pData := data[offset : offset+pkLen-2]
			if _type == 0x5dd || true {
				newData := crypt.ToClientEncr(pData)
				f.Write([]byte(strconv.Itoa(int(pkLen)) + "\t" + hex.EncodeToString(newData) + "\n"))
				//print("\n" + hex.EncodeToString(newData[:4]) + "\n")
				newData[0] = 0
				newData[1] = seq
				newData[0] = crypt.Crc8(newData)
				newData = crypt.ToClientEncr(newData)
				//print(hex.EncodeToString(newData[:4]) + "\n")
				j := 0
				for i := offset; i < offset+pkLen-2; i++ {
					data[i] = newData[j]
					j++
				}
				seq++
			}
			offset += pkLen - 2
		}

		//ioutil.WriteFile("Game/etc/big_bad3", data, 0644)

		pack1, _ = hex.DecodeString("3ad2ba070cd6804220a2a8e7044275a1ca132f7b91d2e73f4699aa94e28407d5a5754516e6b6865627f740fc19e7f5a999764211e131bd53220d3d6c9c3903d3a474440ce5cbc85525b09896662e04d7a7774717f0c34d623001d1a1714112e2b2825323f3c393643404d4a4754514e5b5865626f6c697673707d7a0f07f10e0b181512171fd9262b23dd3a3734313e4b4840425f5f597cf3406d6a676470ae7b787e00ffed0204e4111618e8252a2cdc293e30c03d3244b4414648a855525f59c239ec95c625f887b33e04fe06e690e419161c702d2b3ae4318e3b3835433f4c4946535f6f6a5764616dbb6875727f7c7a0704010e0b1815121f1c292623203d3a3724314e4b4845425f5c595653606d6a6764717e7b7875020f0c090613101d1a1724212e2b2835323f3c394643404d4a5754515e6b6865626f7c797673710e0b0805021f1c191613202d2a2724313e3b38354247c5494653578d5a589b916e6b6865727f7c79061e4b8d1a19a5011e2b282522296c393633304d4a4744425e5b5855526f6c69666c807d7a777af20f0c0909e3101d1a18d4212e2b37c5323f3c46b643404d55a754515e649865626f738976737fe2fa07040111eb0805122f2c292623303d3a3734414e4b4845525f5c596663606d6a7774717e7c0906030fed1a1714111e2b2825222f3c393633304d4a4744415e5b5855526f6c696663707d7a7774010e0b0815121f1dfae50322ab6a4734213e3b4845424f4c595653505d6a6764616e7b7875727fed0a0704011e1b1815122f2c292623303d3a3744414e4b4855524f4c5966636ddb0a7774717e7b1815028f0b3f5503101d0a9603813e3b3835236f4c49464ef63d523c54416e7b7865327718697704135d9b18151")

		dec = crypt.ToClientEncr(pack1)
		print(string(dec))
		print("\n")
		fmt.Println(hex.EncodeToString(dec))

		pack1, _ = hex.DecodeString("adedba070cd6804220a2a8e7044275a1ca132f7b91d2e73f4699aa94e28407d5a5754516e6b6865627f740fc19e7f5a999764211e131bd53220d3d6c9c3903d3a474440ce5cbc85525b09896662e04d7a7774717f0c34d623001d1a1714112e2b2825323f3c393643404d4a4754514e5b5865626f6c697673707d7a0f07f10e0b181512171fd9262b23dd3a3734313e4b4840425f5f597cf3406d6a676470ae7b787e00ffed0204e4111618e8252a2cdc293e30c03d3244b4414648a855525f59c239ec95c625f887b33e04fe06e690e419161c702d2b3ae4318e3b3835433f4c4946535f6f6a5764616dbb6875727f7c7a0704010e0b1815121f1c292623203d3a3724314e4b4845425f5c595653606d6a6764717e7b7875020f0c090613101d1a1724212e2b2835323f3c394643404d4a5754515e6b6865626f7c797673710e0b0805021f1c191613202d2a2724313e3b38354247c5494653578d5a589b916e6b6865727f7c79061e4b8d1a19a5011e2b282522296c393633304d4a4744425e5b5855526f6c69666c807d7a777af20f0c0909e3101d1a18d4212e2b37c5323f3c46b643404d55a754515e649865626f738976737fe2fa07040111eb0805122f2c292623303d3a3734414e4b4845525f5c596663606d6a7774717e7c0906030fed1a1714111e2b2825222f3c393633304d4a4744415e5b5855526f6c696663707d7a7774010e0b0815121f1dfae50322ab6a4734213e3b4845424f4c595653505d6a6764616e7b7875727fed0a0704011e1b1815122f2c292623303d3a3744414e4b4855524f4c5966636ddb0a7774717e7b1815028f0b3f5503101d0a9603813e3b3835236f4c49464ef63d523c54416e7b7865327718697704135d9b18151")

		dec = crypt.ToClientEncr(pack1)
		print(string(dec))
		print("\n")
		fmt.Println(hex.EncodeToString(dec))
	*/
}
