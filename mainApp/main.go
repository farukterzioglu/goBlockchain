package main

import (
	"fmt"
	"strconv"
	goBlockchain "github.com/farukterzioglu/goBlockchain"
)

func main(){
	bc, err := goBlockchain.NewBlockchain()
	if err != nil {
		panic(err)
	}
	defer bc.Dispose()

	bc.AddBlock("Second block")
	bc.AddBlock("Third block")


	iter := bc.Iterator()
	var block = iter.Next()
	for  i := 0; block != nil; i++{
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := goBlockchain.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		block = iter.Next()
	}
}
