package goBlockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

type Block struct {
	Timestamp int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash []byte
	Nonce int
}

func NewBlock(transactions []*Transaction, prevBlockHash []byte) (block *Block) {
	block = &Block {
		time.Now().Unix(),
		transactions,
		prevBlockHash, []byte{},
		0,
	}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return
}

// Serialize serializes the block
func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}
func (block *Block) HashTransactions() []byte{
	var transactions [][]byte

	for _, tx := range block.Transactions {
		transactions = append(transactions, tx.Serialize())
	}
	mTree := NewMerkleTree(transactions)
	return mTree.RootNode.Data
}


func DeserializeBlock(d []byte) (block *Block, err error) {
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err = decoder.Decode(&block)
	if err != nil {
		err = fmt.Errorf("'DeserializeBlock' failed: %v", err)
		return
	}
	return
}
func NewGenesisBlock(coinbase *Transaction) *Block{
	return NewBlock([]*Transaction{coinbase}, []byte{})
}