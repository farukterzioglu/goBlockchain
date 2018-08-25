package goBlockchain

import (
	"github.com/boltdb/bolt"
	_ "fmt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

type Blockchain struct{
	tip []byte
	db  *bolt.DB
}

func (bc *Blockchain) AddBlock(data string) error {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		return err
	}

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		res, err := newBlock.Serialize()
		err = b.Put(newBlock.Hash, res)
		err = b.Put([]byte("l"), newBlock.Hash)

		if err!= nil {
			return err
		}

		bc.tip = newBlock.Hash

		return nil
	})

	return nil
}
func (bc *Blockchain) Dispose() {
	bc.db.Close()
}
func NewGenesisBlock() *Block{
	return NewBlock("Genesis Block", []byte{})
}
func NewBlockchain() (*Blockchain, error) {
	var tip []byte

	//Open db
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			res, err := genesis.Serialize()
			err = b.Put(genesis.Hash, res)
			err = b.Put([]byte("l"), genesis.Hash)

			if err!= nil {
				return err
			}

			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	bc := Blockchain{tip, db}

	return &bc, nil
}

//BlockchainIterator
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}
func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.tip, bc.db}
}
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		if encodedBlock == nil {
			return nil
		}
		blockDes, errAlt := DeserializeBlock(encodedBlock)
		if errAlt != nil {
			return errAlt
		}
		block = blockDes
		return nil
	})

	if err != nil {
		//TODO : Handle error
		return nil
	}

	if block == nil {
		return nil
	}

	i.currentHash = block.PrevBlockHash

	return block
}