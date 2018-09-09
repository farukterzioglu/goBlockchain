package server

import "github.com/farukterzioglu/goBlockchain"

//Requests
func sendVersion(nodeAddress string, blockchain *goBlockchain.Blockchain) {
	bestHeight := blockchain.GetBestHeight()
	payload := gobEncode{version{Version:nodeVersion, BestHeight: bestHeight, AddrFrom:nodeAddress}}

	request := append(commandToBytes( "version"), payload...)
	sendData(nodeAddress, request)
}


func sendGetBlocks(address string) {

}

func sendGetData(address string, payloadType string, txID []byte) {

}

func sendBlock(address string, block *goBlockchain.Block) {

}

func sendInv(node string, invType string, array [][]byte) {

}