package test

import (
	"github.com/farukterzioglu/goBlockchain"
	"testing"
)

func TestInputHash(t *testing.T) {
	wallet := goBlockchain.NewWallet()
	input := goBlockchain.TXInput{Txid:nil, Signature:nil, Vout: 0, PubKey: wallet.PublicKey}

	pubKeyHash := goBlockchain.HashPubKey(wallet.PublicKey)
	if !input.UsesKey(pubKeyHash){
		t.Errorf("Public key doesn't match with address")
	}
}