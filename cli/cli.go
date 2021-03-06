package main

import (
	"flag"
	"fmt"
	"github.com/farukterzioglu/goBlockchain"
	"github.com/farukterzioglu/goBlockchain/server"
	"github.com/pkg/errors"
	"log"
	"os"
	"strconv"
)

type CLI struct {}

func (cli *CLI) Run() {
	cli.validateArgs()

	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Printf("NODE_ID env. var is not set!")
		os.Exit(1)
	}

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")
	sendMine := sendCmd.Bool("mine", false, "Mine immediately on the same node")
	startNodeMiner := startNodeCmd.String("miner", "", "Enable mining mode and send reward to ADDRESS")

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
	case "reindexutxo":
		err := reindexUTXOCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "startnode":
		err := startNodeCmd.Parse(os.Args[2:])
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

		balance := cli.GetBalance(*getBalanceAddress, nodeID)
		fmt.Printf("Balance of %s : %d\n", *getBalanceAddress, balance)
	}

	if printChainCmd.Parsed() {
		cli.printChain(nodeID)
	}

	if createWalletCmd.Parsed() {
		cli.CreateWallet(nodeID)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.CreateBlockchain(*createBlockchainAddress, nodeID)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.Send(*sendFrom, *sendTo, *sendAmount, nodeID, *sendMine)
	}

	if reindexUTXOCmd.Parsed() {
		cli.reindexUTXO(nodeID)
	}

	if startNodeCmd.Parsed() {
		nodeID := os.Getenv("NODE_ID")
		if nodeID == "" {
			startNodeCmd.Usage()
			os.Exit(1)
		}
		cli.startNode(nodeID, *startNodeMiner)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//Privates
func (cli *CLI) CreateBlockchain(address string, nodeID string) {
	if !goBlockchain.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}

	bc, err := goBlockchain.CreateBlockchain(address, nodeID)
	defer bc.Dispose()

	UTXOSet := goBlockchain.UTXOSet{Blockchain: bc}
	UTXOSet.Reindex()

	if err != nil {
		goBlockchain.PanicErr(err, "UTXOSet.ReIndex failed.")
	}
	fmt.Println("Done!")
}
func (cli *CLI) printChain(nodeID string) {
	bc, err := goBlockchain.NewBlockchain(nodeID)
	defer bc.Dispose()

	if err != nil {
		panic(err)
	}

	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("============ Block %x ============\n", block.Hash)
		fmt.Printf("Height: %d\n", block.Height)
		fmt.Printf("Prev. block: %x\n", block.PrevBlockHash)
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
	fmt.Println("  reindexutxo - Rebuilds the UTXO set")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT -mine - Send AMOUNT of coins from FROM address to TO. Mine on the same node, when -mine is set.")
	fmt.Println("  startnode -miner ADDRESS - Start a node with ID specified in NODE_ID env. var. -miner enables mining")
}
func (cli *CLI) GetBalance(address string, nodeID string) int {
	if !goBlockchain.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}
	bc, _ := goBlockchain.NewBlockchain(nodeID)
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
func (cli *CLI) Send(from, to string, amount int, nodeID string, mineNow bool) (err error){
	if !goBlockchain.ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !goBlockchain.ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc, err := goBlockchain.NewBlockchain(nodeID)
	if err != nil {
		errors.Wrap(err, "NewBlockchain failed")
	}

	UTXOSet := goBlockchain.UTXOSet{Blockchain: bc}
	defer bc.Dispose()

	wallets, err := goBlockchain.NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)
	tx := goBlockchain.NewUTXOTransaction(&wallet, to, amount, &UTXOSet)
	if mineNow {
		cbTx := goBlockchain.NewCoinbaseTX(from, "")
		txs := []*goBlockchain.Transaction{cbTx, tx}
		newBlock := bc.MineBlock(txs)
		UTXOSet.Update(newBlock)
	} else {
		server.SendTx(server.KnownNodes[0], tx)
	}

	fmt.Println("Success!")

	return nil
}
func (cli *CLI) CreateWallet(nodeID string) string{
	wallets, _ := goBlockchain.NewWallets(nodeID)
	address := wallets.CreateWallet()
	wallets.SaveToFile(nodeID)
	fmt.Printf("Your new address: %s\n", address)
	return address
}
func (cli *CLI) reindexUTXO(nodeID string) {
	bc, _ := goBlockchain.NewBlockchain(nodeID)
	UTXOSet := goBlockchain.UTXOSet{bc}
	UTXOSet.Reindex()
	count := UTXOSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}
func (cli *CLI) startNode(nodeID, minerAddress string) {
	fmt.Printf("Starting node %s\n", nodeID)
	if len(minerAddress) > 0 {
		if goBlockchain.ValidateAddress(minerAddress) {
			fmt.Println("Mining is on. Address to receive rewards: ", minerAddress)
		} else {
			log.Panic("Wrong miner address!")
		}
	}
	server.StartServer(nodeID, minerAddress)
}