package server

import (
	"fmt"
	"github.com/farukterzioglu/goBlockchain"
	"net"
)

const protocol = "tcp"
const nodeVersion = 1
const commandLength = 12

var nodeAddress string
var KnownNodes = []string { "localhost:3000"}
var miningAddress string

// StartServer starts a node
func StartServer(nodeId, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeId)
	miningAddress = minerAddress

	ln, err := net.Listen(protocol, nodeAddress)
	panicErr(err, "Couldn't start server " + nodeId)
	defer ln.Close()

	bc, err := goBlockchain.NewBlockchain(nodeId)
	panicErr(err,"NewBlockchain failed")

	if nodeAddress != KnownNodes[0] {
		sendVersion(KnownNodes[0], bc)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panicErr(err,"ln.Accept failed")
		}
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