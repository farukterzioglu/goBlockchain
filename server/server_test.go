package server

import "testing"
import (
	_ "github.com/farukterzioglu/goBlockchain"
)

func TestBlockNotNull(t *testing.T){
	StartServer("3001", "1VXWCbQXzrbthh9UVhKyeqnxrvKMcNtkN")
}