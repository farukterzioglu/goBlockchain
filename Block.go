package goBlockchain

import (
	"time"
	"bytes"
	"encoding/gob"
	"fmt"
)

type Block struct {
	Timestamp int64
	Data []byte
	PrevBlockHash []byte
	Hash []byte
	Nonce int
}

func NewBlock(data string, prevBlockHash []byte) (block *Block) {
	block = &Block {
		time.Now().Unix(),
		[]byte(data),
		prevBlockHash, []byte{},
		0,
	}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return
}

func (block *Block) Serialize() (resultBytes []byte, err error){
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err = encoder.Encode(block)
	if err != nil {
		err = fmt.Errorf("'Serialize' failed: %v", err)
		return
	}

	resultBytes = result.Bytes()
	return
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