package test

import (
	cli "github.com/farukterzioglu/goBlockchain/cli"
	"testing"
)
func TestSendTRansaction(t *testing.T){
	from := "account1"
	to := "account2"


	cli := cli.CLI{}
	cli.Send(from, to, 5)

	fromBalance := cli.
}

