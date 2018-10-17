package main

import (
	"flag"
	"fmt"
	"os"
)

var adr = flag.String("address", "", "usage: main  createBlockChain --address <address>")
var cname = flag.String("cname", "", "usage: main  createBlockChain --address <address>")
var data = flag.String("data", "", "usage: main  createBlockChain --address <address>")

func main() {
	//blkc := BlockChain{"test1"}
	//blkc.GetAllBlockHash()
	args := os.Args
	if len(args) > 3 {
		flag.CommandLine.Parse(os.Args[2:])
	} else if len(args) < 2{
		ShowUsage()
		return
	}
	switch args[1] {
	case "createBlockChain":
		if *cname != "" && *adr != ""{
			blkc := NewBlockChain("test1")
			if blkc != nil {
				blkc.GenesisBlock(*adr)
				fmt.Println("create Blockchain success")
			} else {
				fmt.Println("this blockchain is exists")
			}
		} else {
			ShowUsage()
		}

	case "addBlock":
		if *data != "" && *cname != "" && *adr != "" {
			blkc := BlockChain{*cname}
			// TODO
			txs := []*Tx{}
			blkc.AddBlock(txs, 2, 0)
		} else {
			ShowUsage()
		}
	case "showBlockChain":
		if *cname != "" {
			blkc := BlockChain{*cname}
			blkc.GetAllBlockHash()
		} else {
			ShowUsage()
		}

	default:
		ShowUsage()
	}
}

