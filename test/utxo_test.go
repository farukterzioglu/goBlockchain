package test

import (
	"encoding/hex"
	"github.com/farukterzioglu/goBlockchain"
	"log"
	"testing"
)

//Not an actual unit test. Testing transaction concepts
func TestFindUTXO(t *testing.T) {
	//Create new transaction
	wallets, _ := goBlockchain.NewWallets()
	fromAddress := wallets.CreateWallet()
	wallets.SaveToFile()
	fromWallet := wallets.GetWallet(fromAddress)

	//Coinbase transaction
	txin := goBlockchain.TXInput{Txid: []byte{}, Vout: -1, PubKey: []byte("Coinbase")}
	output := *goBlockchain.NewTXOutput(10, fromAddress)
	coinBaseOutputs  := []goBlockchain.TXOutput{ output }
	aTransaction := goBlockchain.Transaction{ ID: nil, Vout: coinBaseOutputs , Vin: []goBlockchain.TXInput{txin}}
	aTransaction.ID = aTransaction.Hash()

	//mapping : txId - output array
	utxoBucket := make(map[string][]goBlockchain.TXOutput)
	txId := hex.EncodeToString(aTransaction.ID)
	utxoBucket[txId] = aTransaction.Vout

	//New transaction
	pubkeyHash := goBlockchain.HashPubKey(fromWallet.PublicKey)

	////receiver
	toAddress := wallets.CreateWallet()
	wallets.SaveToFile()
	var amount = 5

	var unspentOutputs = make(map[string][]int)
	var accumulated  = 0
	////loop all transactions
	for txID, transactionOutputFromBucket := range utxoBucket {
		var outputsFromUTXOSet []goBlockchain.TXOutput = transactionOutputFromBucket
		//loop unspent outputs of transaction
		for outIdx, oneOutput := range outputsFromUTXOSet{
			//Check if it belongs to sender
			if oneOutput.IsLockedWithKey(pubkeyHash){
				if accumulated < amount {
					accumulated += oneOutput.Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}
	}
	if accumulated < amount {
		log.Panic("ERROR: Not enough funds")
	}

	//Create inputs
	var inputs []goBlockchain.TXInput
	for txid, outs := range unspentOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := goBlockchain.TXInput{Txid: txID, Vout: out, PubKey: fromWallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	//Create outputs of new transaction
	var outputs []goBlockchain.TXOutput
	outputToReceiver := *goBlockchain.NewTXOutput(amount, toAddress)
	outputs = append(outputs, outputToReceiver)
	if accumulated > amount {
		outputs = append(outputs, *goBlockchain.NewTXOutput(accumulated-amount, fromAddress)) // a change
	}

	tx := goBlockchain.Transaction{Vin: inputs, Vout: outputs}
	tx.ID = tx.Hash()
	//TODO : Sign the transaction

	//TODO : Assert if change is correct
	//TODO : Assert if receiver got the transaction
	//TODO : Assert if the utxo is updated
}