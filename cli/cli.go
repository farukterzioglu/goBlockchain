package main

import (
	"flag"
	"fmt"
	"github.com/farukterzioglu/goBlockchain"
	"github.com/pkg/errors"
	"log"
	"os"
	"strconv"
)

type CLI struct {}

func (cli *CLI) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)


	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")

	var err error
	printFunc := func(){
		err = printChainCmd.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
	}
	createBlockchainFunc := func(){
		err = createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	}
	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain" :
		printFunc()
	case "p" :
		printFunc()
	case "createblockchain":
		createBlockchainFunc()
	case "c":
		createBlockchainFunc()
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//Privates
func (cli *CLI) createBlockchain(address string) {
	bc, err := goBlockchain.NewBlockchain(address)
	if err != nil {
		panic(err)
	}
	bc.Dispose()
	fmt.Println("Done!")
}
func (cli *CLI) printChain() {
	bc, err := goBlockchain.NewBlockchain("")
	defer bc.Dispose()

	if err != nil {
		panic(err)
	}

	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := goBlockchain.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}
func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  createblockchain (c) -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  printchain (p) - print all the blocks of the blockchain")
}
func (cli *CLI) getBalance(address string){
	bc, err := goBlockchain.NewBlockchain(address)
	defer bc.Dispose()

	if err != nil {
		panic(errors.Wrap(err, "Couldn't create blockchain"))
	}

	balance := 0

	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s : %d\n", address, balance)
}