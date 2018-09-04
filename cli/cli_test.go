package main

import (
	"testing"
)

func TestSendTRansaction(t *testing.T){
	from := "account1"
	to := "account2"

	cli := CLI{}
	cli.Send(from, to, 5)
	fromBalance := cli.GetBalance(from)
	toBalance := cli.GetBalance(to)

	if fromBalance != 5 {
		t.Errorf("Sender's balance isn't correct")
	}

	if toBalance != 5 {
		t.Errorf("Receiver's balance isn't correct")
	}
}


func TestSendMultipleTransaction(t *testing.T){
	cli := CLI{}
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
	}
}

