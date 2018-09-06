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
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

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
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}

		balance := cli.GetBalance(*getBalanceAddress)
		fmt.Printf("Balance of %s : %d\n", *getBalanceAddress, balance)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.CreateBlockchain(*createBlockchainAddress)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.Send(*sendFrom, *sendTo, *sendAmount)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//Privates
func (cli *CLI) CreateBlockchain(address string) {
	if !goBlockchain.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}

	bc, err := goBlockchain.CreateBlockchain(address)
	defer bc.Dispose()

	UTXOSet := goBlockchain.UTXOSet{bc}
	UTXOSet.Reindex()

	if err != nil {
		goBlockchain.PanicErr(err, "UTXOSet.ReIndex failed.")
	}
	fmt.Println("Done!")
}
func (cli *CLI) printChain() {
	bc, err := goBlockchain.NewBlockchain()
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
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  createwallet - Generates a new key-pair and saves it into the wallet file")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  listaddresses - Lists all addresses from the wallet file")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
}
func (cli *CLI) GetBalance(address string) int {
	if !goBlockchain.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc, _ := goBlockchain.NewBlockchain()
	UTXOSet := goBlockchain.UTXOSet{bc}
	defer bc.Dispose()

	balance := 0
	pubKeyHash := goBlockchain.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	return balance
}
func (cli *CLI) Send(from, to string, amount int) (err error){
	if !goBlockchain.ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !goBlockchain.ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc, err := goBlockchain.NewBlockchain()
	if err != nil {
		errors.Wrap(err, "NewBlockchain failed")
	}

	UTXOSet := goBlockchain.UTXOSet{bc}
	defer bc.Dispose()

	tx := goBlockchain.NewUTXOTransaction(from, to, amount, &UTXOSet)
	cbTx := goBlockchain.NewCoinbaseTX(from, "")
	txs := []*goBlockchain.Transaction{cbTx, tx}

	newBlock := bc.MineBlock(txs)
	UTXOSet.Update(newBlock)
	fmt.Println("Success!")

	return nil
}
func (cli *CLI) createWallet() {
	wallets, _ := goBlockchain.NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()
	fmt.Printf("Your new address: %s\n", address)
}