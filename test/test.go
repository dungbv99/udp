package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"

	"strconv"
	"udp/utils"
)

func main() {
	data := &utils.PacketClient{Cmd : 1,Imsi: 452041, 	Msisdn : 84000000000 + rand.Uint64()%(800+20000), Cmnd : "aafsdgfgf"  , Name : "dung" , Dob : "1/1/1111" }
	file, _ := json.MarshalIndent(data, "", " ")
	err := ioutil.WriteFile("./data/"+strconv.FormatUint(data.Msisdn,10)+".json",file,0644)
	if err != nil{
		fmt.Println("err write file : ", err)
	}
}