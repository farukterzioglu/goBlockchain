package test

import (
	"fmt"
	"github.com/farukterzioglu/goBlockchain"
	"testing"
)

func TestLockingOutput(t *testing.T) {
	wallet := goBlockchain.NewWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())
	pubKeyHash := goBlockchain.HashPubKey(wallet.PublicKey)

	output := goBlockchain.NewTXOutput(10, address)
	if !output.IsLockedWithKey(pubKeyHash){
		t.Errorf("Public key doesn't match with address")
	}
}