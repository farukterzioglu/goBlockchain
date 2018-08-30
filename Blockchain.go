package goBlockchain

import (
	"encoding/hex"
	"github.com/boltdb/bolt"
	"log"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type Blockchain struct{
	tip []byte
	db  *bolt.DB
}

func NewBlockchain(address string) (*Blockchain, error) {
	var tip []byte

	//Open db
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
			genesis := NewGenesisBlock(cbtx)

			b, err := tx.CreateBucket([]byte(blocksBucket))
			res, err := genesis.Serialize()
			if err!= nil {
				return err
			}

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

func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction{
	var unspentTXOs []Transaction
	spentTXOs := make(map[string][]int)

	bci := bc.Iterator()

	for {
		block := bci.Next()

		//All transactions from a block
		for _, tx := range block.Transactions{
			txId := hex.EncodeToString(tx.ID)

		Outputs:
			//outputs of transaction
			for outIdx, output := range tx.Vout{
				//If transaction has any spent output
				if spentTXOs[txId] != nil {
					for _, spentout := range spentTXOs[txId]{
						//Check if the spent ont is this
						if spentout == outIdx{
							continue Outputs
						}
					}
				}

				//output belongs to 'address'
				if  output.CanBeUnlockedWith(address){
					unspentTXOs = append(unspentTXOs, *tx)
				}
			}

			//Coin base transactions don't have inputs
			if tx.IsCoinbase() == false {
				//Inputs of transaction
				for _, input := range tx.Vin {
					//Check if the input tx is belong to address's
					if input.CanUnlockOutputWith(address) {
						inputTxId := hex.EncodeToString(input.Txid)
						//Add to spent TXs of address
						spentTXOs[inputTxId] = append(spentTXOs[inputTxId], input.Vout)
					}
				}
			}
		}

		//CHeck if it is the genesis block
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXOs
}
func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	unsTransaactions := bc.FindUnspentTransactions(address)

	for _, tx := range unsTransaactions{
		for _, output := range tx.Vout{
			if output.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, output)
			}
		}
	}
	return UTXOs
}
func (bc *Blockchain) MineBlock(transactions []*Transaction) error {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
		return err
	}

	newBlock := NewBlock(transactions, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		res, err := newBlock.Serialize()
		err = b.Put(newBlock.Hash, res)
		err = b.Put([]byte("l"), newBlock.Hash)

		if err!= nil {
			log.Panic(err)
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
func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.tip, bc.db}
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
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