package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"udemy-modern-golang/dino/communicationlayer/dinoproto3"

	"github.com/golang/protobuf/proto"
)

/*
1- We will serialize some data via proto2
2- We will send this data via TCP to a different service
3- We will deserialize the data via proto2, and print out the extracted values

(client) ---> (server)
A- A TCP client needs to be written to send the data
B- A TCP server to receive the data
*/

func main() {
	op := flag.String("op", "s", "s for server, c for client")
	flag.Parse()
	switch strings.ToLower(*op) {
	case "s":
		RunProto3Server()
	case "c":
		RunProto3Client()
	}
}

func RunProto3Server() {
	l, err := net.Listen("tcp", ":8282")
	if err != nil {
		log.Fatal(err)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		defer l.Close()
		go func(c net.Conn) {
			defer c.Close()
			data, err := ioutil.ReadAll(c)
			if err != nil {
				return
			}
			a := &dinoproto3.Animal{}
			proto.Unmarshal(data, a)
			fmt.Println(a)
		}(c)
	}
}

//client
func RunProto3Client() {
	a := &dinoproto3.Animal{
		Id:         1,
		AnimalType: "Raptor",
		Nickname:   "Rapto",
		Zone:       3,
		Age:        20,
	}
	data, err := proto.Marshal(a)
	if err != nil {
		log.Fatal(err)
	}
	sendData(data)
}

func sendData(data []byte) {
	c, err := net.Dial("tcp", "127.0.0.1:8282")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	c.Write(data)
}
