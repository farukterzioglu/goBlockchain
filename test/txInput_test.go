package test

import (
	"bytes"
	"crypto/sha256"
	"github.com/farukterzioglu/goBlockchain"
	"golang.org/x/crypto/ripemd160"
	"log"
	"testing"
)

type TXInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	PubKey    []byte
}
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return publicRIPEMD160
}

//TODO : fails
func TestInputHash(t *testing.T) {
	wallet := goBlockchain.NewWallet()
	input := goBlockchain.TXInput{Txid:nil, Signature:nil, Vout: 0, PubKey: wallet.PublicKey}

	pubKeyHash := goBlockchain.HashPubKey(wallet.PublicKey)
	if !input.UsesKey(pubKeyHash){
		t.Errorf("Public key doesn't match with address")
	}
}