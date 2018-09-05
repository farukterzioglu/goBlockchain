package goBlockchain

import (
	"encoding/hex"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

const utxoBucket = "chainstate"

type UTXOSet struct {
	Blockchain *Blockchain
}

func (set UTXOSet) ReIndex() error {
	db := set.Blockchain.db
	bucketName := []byte(utxoBucket)
	
	err := db.Update(func(tx *bolt.Tx) (err error) {
		err = tx.DeleteBucket(bucketName)
		_, err = tx.CreateBucket(bucketName)
		return
	})

	if err != nil {
		return errors.Wrap(err, "db.Update at ReIndex")
	}

	UTXO := set.Blockchain.FindUTXO()

	err = db.Update(func(tx *bolt.Tx) (err error ){
		b := tx.Bucket(bucketName)

		for txID, outs := range UTXO {
			key, _ := hex.DecodeString(txID)
			err = b.Put(key, outs.Serialize())
		}
		return
	})
	return nil
}
func (set UTXOSet) FindSpendableOutputs(pubKeyHash []byte,  amount int) (int, map[string][]int , error){
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	db := set.Blockchain.db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := DeserializeOutputs(v)

			for outIdx, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
					accumulated += out.Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}
		return nil
	})

	if err != nil {
		errors.Wrap(err, "db.View at FindSpendableOutputs")
	}

	return accumulated, unspentOutputs, nil
}
