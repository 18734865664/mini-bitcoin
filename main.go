package main

import (
	"flag"
	"fmt"
	"os"
)

var adr = flag.String("a", "", "usage: main  createBlockChain --a <address>")
var miner = flag.String("m", "", "usage: main  createBlockChain --m <address>")
var cname = flag.String("c", "", "usage: main  createBlockChain --c <blockChainName>")
var data = flag.String("d", "", "usage: main  createBlockChain --d <message data>")
var target = flag.String("t", "", "usage: main  createBlockChain --t <targetAddr1:count1,targetAddr2:count2...>")

func main() {

	//chaTx := make(chan *Tx, 3)
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
		if *cname != ""{
			blkc := NewBlockChain("test1")
			addr := blkc.GetNewAddress()
			if blkc != nil {
				blkc.GenesisBlock(addr)
				fmt.Println("create Blockchain success")
			} else {
				fmt.Println("this blockchain is exists")
			}
		} else {
			ShowUsage()
		}

	case "addBlock":
		if *cname != "" && *adr != "" {
			if !CheckAddress(*adr){
				ShowUsage()
				return
			}
			blkc := &BlockChain{ChainName: *cname}
			// TODO
			txs := []Tx{}
			blkc.AddBlock(txs, 2, 0)
		} else {
			ShowUsage()
		}
	case "showBlockChainBlock":
		if *cname != "" {
			blkc := GetNewBlockChainObj(*cname)
			blkc.GetAllBlockHash()
		} else {
			ShowUsage()
		}
	case "getCountInfo":
		fmt.Println("cname: ", *cname)
		fmt.Println("adr: ", *adr)
		if *cname != "" && *adr != ""{
			if !CheckAddress(*adr){
				fmt.Println("checkaddress error")
				ShowUsage()
				return
			}
			blkc := GetNewBlockChainObj(*cname)
			blkc.ShowCountCoin(*adr)
		} else {
			ShowUsage()
		}
	case "transfer":
		if *cname != "" && *adr != "" && *target != ""{
			var addr string
			if *miner == ""{
				addr = *adr
			} else {
				addr = *miner
			}
			if !CheckAddress(*adr){
				ShowUsage()
				return
			}
			blkc := GetNewBlockChainObj(*cname)
			tx := blkc.CreateCommTrans(*adr, *target)
			coinBasetx := CoinBaseTx(addr, *data)
			txs := []Tx{}
			txs = append(txs, coinBasetx)
			if len(tx.Inputs) != 0{
				txs = append(txs, tx)
			}
			blkc.AddBlock(txs, 1, 0)
			//chaTx <- tx
			//chaTx <- coinBasetx
		} else {
			ShowUsage()
		}
	case "getAddress":
		if *cname != ""{
			blkc := GetNewBlockChainObj(*cname)
			addr := blkc.GetNewAddress()
			fmt.Printf("newAddress: %s\n", addr)
		}else{
			ShowUsage()
		}
	case "listWallet":
		if *cname != ""{
			blkc := GetNewBlockChainObj(*cname)
			blkc.ShowWallet()
		}else{
			ShowUsage()
		}
	default:
		ShowUsage()
	}
}

