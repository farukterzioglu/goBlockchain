package goBlockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"log"
)

const subsidy = 10

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

// IsCoinbase checks whether the transaction is coinbase
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}
// SetID sets ID of a transaction
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

type TXOutput struct {
	Value        int
	ScriptPubKey string
}
type TXInput struct {
	//Origin transaction
	Txid      []byte
	//Order of output
	Vout      int
	ScriptSig string
}

func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx
}

func NewTransaction(from, to string, amount int, bc *Blockchain) (*Transaction, error){
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs, err := bc.FindSpendableOutputs(from, amount)
	if err != nil {
		return nil, errors.Wrap(err, "Error while finding spendable outputs")
	}

	if acc< amount{
		log.Panic("Not enough balance ")
	}

	for txid, outs := range validOutputs{
		txId, err := hex.DecodeString(txid)

		if err != nil {
			return nil, errors.Wrap(err, "DecodeString failed.")
		}

		for _, out := range outs{
			input := TXInput{Txid:txId, Vout:out, ScriptSig:from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TXOutput{Value:amount, ScriptPubKey:to})
	if acc > amount {
		outputs = append(outputs, TXOutput{Value: acc -amount, ScriptPubKey:from})
	}

	tx := Transaction{Vin:inputs, Vout:outputs,ID:nil}
	tx.SetID()

	return &tx, nil
}