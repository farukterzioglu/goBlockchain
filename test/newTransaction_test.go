package test

import (
	"github.com/farukterzioglu/goBlockchain"
	"log"
	"testing"
)

func TestSendTransaction(t *testing.T){
	node := "3000"
	//Create new wallet
	wallets, _ := goBlockchain.NewWallets(node)
	fromAddress := wallets.CreateWallet()
	wallets.SaveToFile(node)

	wallet := wallets.GetWallet(fromAddress)

	toAddress := wallets.CreateWallet()
	wallets.SaveToFile(node)

	//Validate address
	if !goBlockchain.ValidateAddress(fromAddress) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !goBlockchain.ValidateAddress(toAddress) {
		log.Panic("ERROR: Receiver address is not valid")
	}

	var bc *goBlockchain.Blockchain

	//Create or get blockchain
	if goBlockchain.DbExists(node) {
		panic("remove blockchain.db before running test")
	} else {
		bc, _ = goBlockchain.CreateBlockchain(fromAddress, node)
	}

	//Create UTXO set & reindex it
	UTXOSet := goBlockchain.UTXOSet{Blockchain: bc}
	UTXOSet.Reindex()
	bc.Dispose()

	//Get created blockchain
	bc, _ = goBlockchain.NewBlockchain(node)
	defer bc.Dispose()

	//Get UTXO set
	UTXOSet = goBlockchain.UTXOSet{Blockchain: bc}

	//New transaction & reward transaction
	tx := goBlockchain.NewUTXOTransaction(&wallet, toAddress, 6, &UTXOSet)
	cbTx := goBlockchain.NewCoinbaseTX(fromAddress, "")
	txs := []*goBlockchain.Transaction{cbTx, tx}

	newBlock := bc.MineBlock(txs)
	UTXOSet.Update(newBlock)

	//Get balance
	UTXOSet = goBlockchain.UTXOSet{Blockchain: bc}

	balance := 0
	pubKeyHash := goBlockchain.Base58Decode([]byte(fromAddress))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	if balance != 14 {
		panic("balance is not correct")
	}
}