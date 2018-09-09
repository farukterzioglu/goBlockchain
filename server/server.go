package server

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/farukterzioglu/goBlockchain"
	"io/ioutil"
	"log"
	"net"
)

const protocol = "tcp"
const nodeVersion = 1
const commandLength = 12

var nodeAddress = "localhost:%s"
var knownNodes = []string { "localhost:3000"}
var miningAddress string

func StartServer(nodeId string, minerAddress string){
	nodeAddress = fmt.Sprintf(nodeAddress, nodeId)
	miningAddress = minerAddress

	ln, err := net.Listen(protocol, nodeAddress)
	defer ln.Close()
	panicErr(err, "Couldn't start server " + nodeId)

	bc, err := goBlockchain.NewBlockchain(nodeId)
	panicErr(err,"NewBlockchain failed")

	if nodeAddress != knownNodes[0] {
		sendVersion(knownNodes[0], bc)
	}

	for {
		conn , err := ln.Accept()
		log.Printf(err.Error())
		go handleConnection(conn, bc)
	}
}

//Helpers
func commandToBytes(command string) []byte {
	var bytes [commandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}
func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}
func panicErr(err error, messageToWrap string){
	goBlockchain.PanicErr(err, messageToWrap)
}