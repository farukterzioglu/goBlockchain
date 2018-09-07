package test
/*
import "testing"
import (
	"github.com/farukterzioglu/goBlockchain"
	_ "fmt"
)

func TestBlockNotNull(t *testing.T){
	block := goBlockchain.NewBlock("test", []byte{})
	//fmt.Printf("Hash :%s\n", block.Hash)

	if len(block.Hash) == 0 {
		t.Errorf("Hash is wrong : %s", block.Hash)
	}
}

func TestBlockIsEqual(t *testing.T){
	var blockTests = []struct {
		 input string
		 expected []byte
	}{
		//TODO : Get hash value and replace 'test' value below
		{"test", []byte{}},
	}

	for _,tt := range blockTests{
		block := goBlockchain.NewBlock(tt.input, []byte{})

		actual := block.Hash
		//TODO : Check values
		if len(actual) != len(tt.expected) {
			//TODO : Activate this
			//t.Errorf("Hash is wrong : %s", block.Hash)
		}
	}
}


func TestSerializeDeserializeBlock(t *testing.T){
	block := goBlockchain.NewBlock("test", []byte{})
	res, err := block.Serialize()

	if err != nil{
		t.Errorf("Couldn't seralized block")
	}

	newblock,err := goBlockchain.DeserializeBlock(res)
	if err != nil {
		t.Errorf("Couldn't deseralized block data")
	}

	if string(newblock.Data) != "test" {
		t.Errorf("Data values doesn't match")
	}
}
*/