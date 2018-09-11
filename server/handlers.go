package server

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/farukterzioglu/goBlockchain"
	"io/ioutil"
	"net"
)

var blocksInTransit [][]byte
var mempool = make(map[string]goBlockchain.Transaction)

//Handlers
func handleConnection(conn net.Conn, bc *goBlockchain.Blockchain){
	request, err := ioutil.ReadAll(conn)
	panicErr(err, "Reading connection details failed")

	command := bytesToCommand(request[:commandLength])
	fmt.Printf("Received command : %s \n", command)

	switch command {
	case "version" :
		handleVersion(request, bc)
	default :
		fmt.Printf("Unknown command")
	}

	conn.Close()
}

func handleVersion(request []byte, blockchain *goBlockchain.Blockchain) {
	var buff bytes.Buffer
	var payload version

	buff.Write(request[commandLength:])
	decoder := gob.NewDecoder(&buff)
	err := decoder.Decode(payload)
	panicErr(err, "Couldn't decode payload")

	currentBestHeight := blockchain.GetBestHeight()
	foreignerBestHeight := payload.BestHeight

	if currentBestHeight < foreignerBestHeight {
		sendGetBlocks(payload.AddrFrom)
	} else if currentBestHeight > foreignerBestHeight {
		sendVersion(payload.AddrFrom, blockchain)
	}

}

func handleGetBlocks(request []byte, bc *goBlockchain.Blockchain) {
	var buff bytes.Buffer
	var payload getblocks

	buff.Write(request[commandLength:])
	decoder := gob.NewDecoder(&buff)
	err := decoder.Decode(payload)
	panicErr(err, "Couldn't decode payload")

	blocks := bc.GetBlockHashes()
	sendInv(payload.AddrFrom, "block", blocks)
}

func handleInv(request []byte, blockchain *goBlockchain.Blockchain){
	var buff bytes.Buffer
	var payload inv

	buff.Write(request[commandLength:])
	decoder := gob.NewDecoder(&buff)
	err := decoder.Decode(payload)
	panicErr(err, "Couldn't decode payload")

	fmt.Printf("Recevied inventory with %d %s\n", len(payload.Items), payload.Type)

	if payload.Type == "block" {
		blocksInTransit = payload.Items

		blockHash := payload.Items[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		newInTransit := [][]byte{}
		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}

	if payload.Type == "tx" {
		txID := payload.Items[0]

		if mempool[hex.EncodeToString(txID)].ID == nil {
			sendGetData(payload.AddrFrom, "tx", txID)
		}
	}
}

func handleGetData(request []byte, bc *goBlockchain.Blockchain) {
	var buff bytes.Buffer
	var payload getdata

	buff.Write(request[commandLength:])
	decoder := gob.NewDecoder(&buff)
	err := decoder.Decode(payload)
	panicErr(err, "Couldn't decode payload")

	if payload.Type == "block" {
		block, err := bc.GetBlock([]byte(payload.ID))
		panicErr(err, "bc.GetBlock failed.")

		sendBlock(payload.AddrFrom, &block)
	}

	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := mempool[txID]

		SendTx(payload.AddrFrom, &tx)
	}
}

func handleBlock(request []byte, bc *goBlockchain.Blockchain) {
	var buff bytes.Buffer
	var payload block

	buff.Write(request[commandLength:])
	decoder := gob.NewDecoder(&buff)
	err := decoder.Decode(payload)
	panicErr(err, "Couldn't decode payload")

	blockData := payload.Block
	block := goBlockchain.DeserializeBlock(blockData)

	fmt.Println("Recevied a new block!")
	//TODO : Verify block
	bc.AddBlock(block)

	fmt.Printf("Added block %x\n", block.Hash)

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := goBlockchain.UTXOSet{bc}
		//TODO : Use Update instead of reindex
		UTXOSet.Reindex()
	}
}

func handleTx(request []byte, bc *goBlockchain.Blockchain) {
	var buff bytes.Buffer
	var payload tx

	buff.Write(request[commandLength:])
	decoder := gob.NewDecoder(&buff)
	err := decoder.Decode(payload)
	panicErr(err, "Couldn't decode payload")

	txData := payload.Transaction
	tx := goBlockchain.DeserializeTransaction(txData)
	//TODO : verify transactions
	mempool[hex.EncodeToString(tx.ID)] = tx

	//Central node doesn't do mining, forward to other nodes
	if nodeAddress == KnownNodes[0] {
		for _, node := range KnownNodes {
			if node != nodeAddress && node != payload.AddFrom {
				sendInv(node, "tx", [][]byte{tx.ID})
			}
		}
	} else {
		//Start mining if there are transactions more than 2
		if len(mempool) >= 2 && len(miningAddress) > 0 {
		MineTransactions:
			var txs []*goBlockchain.Transaction

			for id := range mempool {
				tx := mempool[id]
				if bc.VerifyTransaction(&tx) {
					txs = append(txs, &tx)
				}
			}

			if len(txs) == 0 {
				fmt.Println("All transactions are invalid! Waiting for new ones...")
				return
			}

			cbTx := goBlockchain.NewCoinbaseTX(miningAddress, "")
			txs = append(txs, cbTx)

			newBlock := bc.MineBlock(txs)
			UTXOSet := goBlockchain.UTXOSet{bc}
			//TODO : Use Update
			UTXOSet.Reindex()

			fmt.Println("New block is mined!")

			for _, tx := range txs {
				txID := hex.EncodeToString(tx.ID)
				delete(mempool, txID)
			}

			for _, node := range KnownNodes {
				if node != nodeAddress {
					sendInv(node, "block", [][]byte{newBlock.Hash})
				}
			}

			if len(mempool) > 0 {
				goto MineTransactions
			}
		}
	}
}
