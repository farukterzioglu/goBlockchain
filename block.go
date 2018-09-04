package goBlockchain

import (
	"log"
	"time"
	"bytes"
	"encoding/gob"
	"fmt"
	"crypto/sha256"
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
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range block.Transactions {
		txHashes = append(txHashes, tx.Hash())
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
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