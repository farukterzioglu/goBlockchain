package main

import (
	"testing"
)

//TODO : Implement deleting existing blockchain
//Delete existing blockchain before use
func TestSendTransaction(t *testing.T){
	//Arrange
	node := "3000"

	cli := CLI{}
	from := cli.CreateWallet(node)
	to := cli.CreateWallet(node)
	cli.CreateBlockchain(from, node)

	//Act
	cli.Send(from, to, 5, node, true)

	//Assert
	fromBalance := cli.GetBalance(from, node)
	toBalance := cli.GetBalance(to, node)

	if fromBalance != 5 {
		t.Errorf("Sender's balance isn't correct")
	}

	if toBalance != 5 {
		t.Errorf("Receiver's balance isn't correct")
	}
}


func TestSendMultipleTransaction(t *testing.T){
	//TODO : Fix this
/*	cli := CLI{}
	cli.CreateBlockchain("acc1")
	cli.Send("acc1", "acc2", 5)
	cli.Send("acc2", "acc3", 2)
	cli.Send("acc1", "acc3", 4)

	balance := cli.GetBalance("acc3")
	if balance != 6 {
		t.Errorf("Balance isn't correct")
	}

	balance = cli.GetBalance("acc1")
	if balance != 1 {
		t.Errorf("Balance isn't correct")
	}

	balance = cli.GetBalance("acc2")
	if balance != 3 {
		t.Errorf("Balance isn't correct")
	}*/
}

