package main

import (
	"github.com/farukterzioglu/goBlockchain"
)

func main() {
	bc, _ := goBlockchain.NewBlockchain()
	defer bc.Dispose()

	cli := CLI{bc}
	cli.Run()
}