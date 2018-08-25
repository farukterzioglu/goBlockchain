package main

import (
	"fmt"
	"strconv"
	goBlockchain "github.com/farukterzioglu/goBlockchain"
)

func main(){
	bc := goBlockchain.NewBlockchain()
	bc.AddBlock("Second block")
	bc.AddBlock("Third block")

	for _, block := range bc.Blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := goBlockchain.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
