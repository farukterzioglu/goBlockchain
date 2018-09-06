package test
/*
import "testing"
import (
	"github.com/farukterzioglu/goBlockchain"
	"fmt"
)

func TestGenesisBlock(t *testing.T){
	genesisBlock := goBlockchain.NewGenesisBlock()
	if string(genesisBlock.Data) != "Genesis Block" {
		t.Errorf("Genesis block data is wrong : %s", genesisBlock.Data)
	}
}

func TestNewBlockChain(t *testing.T){
	newBlockChain, err := goBlockchain.NewBlockchain()
	defer newBlockChain.Dispose()

	if err != nil {
		t.Errorf("Couldn't create new blockchain")
	}
	if newBlockChain == nil {
		t.Errorf("Couldn't create new blockchain")
	}
}

func TestIterator(t *testing.T){
	//Arrange
	newBlockChain, _:= goBlockchain.NewBlockchain()
	defer newBlockChain.Dispose()

	newBlockChain.AddBlock("second block data")
	iter := newBlockChain.Iterator()

	//Act
	var i = 0
	for  {
		next := iter.Next()
		if len(next.Hash) <= 0{
			t.Errorf("Length of hash is not valid")
		}
		fmt.Printf("%d. Lenght of hash : %d\n",i ,len(next.Hash))

		i++
		if len(next.PrevBlockHash) == 0 {
			break
		}
	}

	//Assert
	if i < 2 {
		t.Errorf("Didn't iterate all blocks")
	}
}
//TODO : Test adding block

//TODO : Test Generating hash with PoW
*/