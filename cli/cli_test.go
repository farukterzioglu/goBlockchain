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
		t.Errorf("Couldn't deseralized block data")
	}

	if toBalance != 5 {
		t.Errorf("Couldn't deseralized block data")
	}
}

