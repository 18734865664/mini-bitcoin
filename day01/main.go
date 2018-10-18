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
			acc := GetNewAccount()
			if blkc != nil {
				blkc.GenesisBlock(acc.Addr)
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
			txs := []Tx{}
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
			if *miner == ""{
				*miner = *adr
			}
			blkc := BlockChain{*cname}
			tx := blkc.CreateCommTrans(*adr, *target)
			coinBasetx := CoinBaseTx(*miner, *data)
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
			blkc := BlockChain{*cname}
			addr := blkc.GetNewAddress()
			fmt.Printf("newAddress: %s\n", addr)
		}else{
			ShowUsage()
		}
	case "listWallet":
			if *cname != ""{
			blkc := BlockChain{*cname}
			blkc.ShowWallet()
		}else{
			ShowUsage()
		}
	default:
		ShowUsage()
	}
}

