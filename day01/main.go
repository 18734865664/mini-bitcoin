package main

import (
	"flag"
	"fmt"
	"os"
)

var adr = flag.String("a", "", "usage: main  createBlockChain --a <address>")
var cname = flag.String("c", "", "usage: main  createBlockChain --c <blockChainName>")
// var data = flag.String("data", "", "usage: main  createBlockChain --address <ad>")
var target = flag.String("t", "", "usage: main  createBlockChain --t <targetAddr1:count1,targetAddr2:count2...>")

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
		if *cname != "" && *adr != "" {
			blkc := BlockChain{*cname}
			// TODO
			txs := []*Tx{}
			blkc.AddBlock(txs, 2, 0)
		} else {
			ShowUsage()
		}
	case "showBlockChainBlock":
		if *cname != "" {
			blkc := BlockChain{*cname}
			blkc.GetAllBlockHash()
		} else {
			ShowUsage()
		}
	case "getCountInfo":
		if *cname != "" && *adr != ""{
			blkc := BlockChain{*cname}
			blkc.ShowCountCoin(*adr)
		} else {
			ShowUsage()
		}
	case "transfer":
		if *cname != "" && *adr != "" && *target != ""{
			blkc := BlockChain{*cname}
			blkc.CreateCommTrans(*adr, *target)
		} else {
			ShowUsage()
		}
	case "getAddress":
		if *cname != ""{
			blkc := BlockChain{*cname}
			addr := blkc.GetNewAddress()
			fmt.Printf("newAddress: %x\n", addr)
		}else{
			ShowUsage()
		}
	default:
		ShowUsage()
	}
}

