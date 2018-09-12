package server

type command struct {
	AddrFrom string
}

type addr struct {
	AddrFrom string
	AddrList []string
}

type version struct {
	Version int
	BestHeight int
	AddrFrom string
}

type getblocks struct {
	AddrFrom string
}

type inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

type getdata struct {
	AddrFrom string
	Type     string
	ID       []byte
}

type block struct {
	AddrFrom string
	Block    []byte
}

type tx struct {
	AddrFrom    string
	Transaction []byte
}

