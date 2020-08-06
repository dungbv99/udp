package main

import (
	"udp/utils"
	"fmt"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"sync"
	"time"

	"net"
)

func testChange(wg *sync.WaitGroup, numberPacket uint64){
	defer wg.Done()
	s, err := net.ResolveUDPAddr("udp4", "localhost:1234")
	c, err := net.DialUDP("udp", nil, s)
	if err != nil {
		fmt.Println(err)
		return
	}
	var i uint64
	for i = 0; i < numberPacket; i++{
		data := &utils.PacketClient{Cmd : 2,Imsi: 4520411213, 	Msisdn : 84000000000 + rand.Uint64()%(numberPacket+20000), Cmnd : "aafsdgfgf"  , Name : "dung" , Dob : "1/1/1111" }
		out, err := proto.Marshal(data)
		if err != nil{
			fmt.Println(err)
		}


		c.Write(out)

		res := make([]byte,512)
		length, _, _ := c.ReadFromUDP(res)
		response := &utils.PacketServer{}
		err = proto.Unmarshal(res[:length],response)
		if err != nil {
			fmt.Println(err)
		} else{
			//fmt.Println("rescode : ", response.Rescode, "reason : ", response.Reason)
		}



		time.Sleep(time.Millisecond)
	}
}


func testAdd(wg *sync.WaitGroup, numberPacket uint64){
	defer wg.Done()
	s, err := net.ResolveUDPAddr("udp4", "localhost:1234")
	c, err := net.DialUDP("udp", nil, s)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()
	var i uint64
	for i = 0; i < numberPacket; i++{

		data := &utils.PacketClient{Cmd : 1,Imsi: 452041, 	Msisdn : 84000000000 + rand.Uint64()%(numberPacket+20000), Cmnd : "aafsdgfgf"  , Name : "dung" , Dob : "1/1/1111" }
		out, err := proto.Marshal(data)
		if err != nil{
			fmt.Println(err)
		}


		c.Write(out)

		res := make([]byte,512)
		length, _, _ := c.ReadFromUDP(res)
		response := &utils.PacketServer{}
		err = proto.Unmarshal(res[:length],response)
		if err != nil {
			fmt.Println(err)
		}else{
			//fmt.Println("rescode : ", response.Rescode, "reason : ", response.Reason)
		}




		time.Sleep(time.Millisecond)


	}


}


func testDelete(wg *sync.WaitGroup, numberPacket uint64){
	defer wg.Done()
	s, err := net.ResolveUDPAddr("udp4", "localhost:1234")
	c, err := net.DialUDP("udp", nil, s)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()
	var i uint64
	for i = 0; i < numberPacket; i++{

		data := &utils.PacketClient{Cmd : 3,Imsi: 452041, 	Msisdn : 84000000000 + rand.Uint64()%(numberPacket+20000), Cmnd : "aafsdgfgf"  , Name : "dung" , Dob : "1/1/1111" }
		out, err := proto.Marshal(data)
		if err != nil{
			fmt.Println(err)
		}


		c.Write(out)

		res := make([]byte,512)
		length, _, _ := c.ReadFromUDP(res)
		response := &utils.PacketServer{}
		err = proto.Unmarshal(res[:length],response)
		if err != nil {
			fmt.Println(err)
		}else{
			//fmt.Println("rescode : ", response.Rescode, " reason : ", response.Reason)
		}




		time.Sleep(time.Millisecond)


	}


}

func main(){
	//d := utils.Packet{Cmd: 1, Cmnd: "123453213442", Dob: "18/01/1999", Msisdn: 123211231, Name: "dung"}
	//fmt.Println("d : ", d.Encode())
	start := time.Now()

	var wg sync.WaitGroup
	var i uint64
	var numberPacket uint64
	numberPacket = 80000
	for i = 0; i < 8; i++{
		wg.Add(1)
		go testAdd(&wg,numberPacket)

	}
	for i = 0; i < 10; i++{
		wg.Add(1)
		go testChange(&wg,numberPacket)

	}
	for i = 0; i < 2; i++{
		wg.Add(1)
		go testDelete(&wg,numberPacket)

	}
	wg.Wait()
	fmt.Println(time.Since(start))








	//
	//a := &utils.Packet{Name: "asdasd", Msisdn: 12312, Dob: "asdassas",Cmnd: "012312039", Cmd: 1}
	//out, err := proto.Marshal(a)
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(out)
	//b := &utils.Packet{}
	//if err := proto.Unmarshal(out, b); err != nil {
	//	log.Fatalln("Failed to parse address book:", err)
	//}
	//fmt.Println(b)
}
