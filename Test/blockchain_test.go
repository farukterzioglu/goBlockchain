package Test

import "testing"
import (
	"github.com/farukterzioglu/goBlockchain"
	_ "fmt"
)

func TestGenesisBlock(t *testing.T){
	genesisBlock := goBlockchain.NewGenesisBlock()
	if string(genesisBlock.Data) != "Genesis Block" {
		t.Errorf("Genesis block data is wrong : %s", genesisBlock.Data)
	}
}

func TestNewBlockChain(t *testing.T){
	newBlockChain, err := goBlockchain.NewBlockchain()
	if err != nil {
		t.Errorf("Couldn't create new blockchain")
	}
	if newBlockChain == nil {
		t.Errorf("Couldn't create new blockchain")
	}
}

//TODO : Test adding block

//TODO : Test Generating hash with PoW

