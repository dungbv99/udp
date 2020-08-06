package main

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
	"udp/utils"
	//"bytes"
	"fmt"
	"net"
)

var m = make(map[string]utils.PacketClient)
var m_rm = make(map[string]int32)
var cnt uint64 = 0
var mux sync.Mutex


func server(connection *net.UDPConn){

	for{
		inputBytes  := make([]byte, 4096)
		length, addr, _ := connection.ReadFromUDP(inputBytes )



		message := &utils.PacketClient{}
		if err := proto.Unmarshal(inputBytes[:length], message); err != nil {
			response := &utils.PacketServer{Cmd: 1, Reason: err.Error(), Rescode: 400}
			res , err := proto.Marshal(response)
			if err != nil{
				fmt.Println(err)
			}
			connection.WriteToUDP(res, addr)
		}
		if message.Cmd == 1{
			mux.Lock()
			_, found := m["./data/"+strconv.FormatUint(message.Msisdn,10)+".json"]
			mux.Unlock()
			if found == true{
				response := &utils.PacketServer{Cmd: 1, Reason: "file exist", Rescode: 300}
				res , err := proto.Marshal(response)
				if err != nil{
					fmt.Println(err)
				}
				connection.WriteToUDP(res, addr)
			}else {
				file, _ := json.MarshalIndent(message, "", " ")
				err := ioutil.WriteFile("./data/"+strconv.FormatUint(message.Msisdn,10)+".json",file,0644)
				if err != nil{
					fmt.Println("err write file : ", err)
				}
				mux.Lock()
				m["./data/"+strconv.FormatUint(message.Msisdn,10)+".json"] = *message
				mux.Unlock()
				response := &utils.PacketServer{Cmd: 1, Reason: "ok create", Rescode: 200}
				res , err := proto.Marshal(response)
				if err != nil{
					fmt.Println(err)
				}
				connection.WriteToUDP(res, addr)
			}
		}else if message.Cmd == 2{
			mux.Lock()
			data, found := m["./data/"+strconv.FormatUint(message.Msisdn,10)+".json"]
			mux.Unlock()
			if found == false{
				response := &utils.PacketServer{Cmd: 2, Reason: "file " + strconv.FormatUint(message.Msisdn,10)+".json  "+"dose not exist to change ", Rescode: 500}
				res , err := proto.Marshal(response)
				if err != nil{
					fmt.Println(err)
				}
				connection.WriteToUDP(res, addr)
			}else {
				if len(message.Cmnd) > 0{
					data.Cmnd = message.Cmnd;
				}
				if len(message.Dob) > 0{
					data.Dob = message.Dob
				}
				if len(message.Name) > 0{
					data.Name = message.Name
				}
				file, _ := json.MarshalIndent(message, "", " ")
				err := ioutil.WriteFile("./data/"+strconv.FormatUint(message.Msisdn,10)+".json",file,0644)
				response := &utils.PacketServer{Cmd: 1, Reason: "ok change", Rescode: 200}
				res , err := proto.Marshal(response)
				if err != nil{
					fmt.Println(err)
				}else {
					connection.WriteToUDP(res, addr)
				}
			}
		}else if message.Cmd == 3{
			mux.Lock()
			_, found := m["./data/"+strconv.FormatUint(message.Msisdn,10)+".json"]
			mux.Unlock()
			if found == false {
				response := &utils.PacketServer{Cmd: 2, Reason: "file " + strconv.FormatUint(message.Msisdn,10)+".json  "+"dose not exist to delete ", Rescode: 400}
				res , err := proto.Marshal(response)
				if err != nil{
					fmt.Println(err)
				}
				connection.WriteToUDP(res, addr)
			}else{
				mux.Lock()
				_, f := m_rm["./data/"+strconv.FormatUint(message.Msisdn,10)+".json"]
				mux.Unlock()
				if f == false{
					err := os.Rename("./data/"+strconv.FormatUint(message.Msisdn,10)+".json", "./rm/"+strconv.FormatUint(message.Msisdn,10)+".json")
					if err != nil{
						fmt.Println("err delete : ", err)
					}
					mux.Lock()
					m_rm["./data/"+strconv.FormatUint(message.Msisdn,10)+".json"] = 1
					mux.Unlock()
				}
				response := &utils.PacketServer{Cmd: 1, Reason: "ok delete", Rescode: 200}
				res , err := proto.Marshal(response)
				if err != nil{
					fmt.Println(err)
				}else {
					connection.WriteToUDP(res, addr)
				}
				mux.Lock()
				delete(m, "./data/"+strconv.FormatUint(message.Msisdn,10)+".json")
				mux.Unlock()
			}
		}
		mux.Lock()
		cnt++
		mux.Unlock()
	}

}

func printCnt(){
	for{
		mux.Lock()
		fmt.Println("cnt : ", cnt)
		mux.Unlock()
		time.Sleep(time.Second*10)
	}

}

func main()  {
	src := "./data/"
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		file, _ := ioutil.ReadFile(path)
		data := utils.PacketClient{}
		_ = json.Unmarshal([]byte(file), &data)
		m[path] = data
		return nil
	})
	if err != nil {
		panic(err)
	}
	src_err := "./rm/"
	err = filepath.Walk(src_err, func(path string, info os.FileInfo, err error) error {
		m_rm[path] = 1
		return nil
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("m : ",len(m))
	s, err := net.ResolveUDPAddr("udp4", ":1234")
	if err != nil {
		fmt.Println(err)
		return
	}
	connection, err := net.ListenUDP("udp", s)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connection.Close()

	cnt = 0
	for i:=0; i<2; i++{
		go server(connection)
	}
	go printCnt()
	var wg sync.WaitGroup
	wg.Add(2)
	wg.Wait()
}
