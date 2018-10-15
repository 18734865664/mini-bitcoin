package main

import (
	"fmt"
	"encoding/hex"
)

type BlockChain struct {
	Blocks []*Block
}

func NewBlockChain()*BlockChain{
	return &BlockChain{
		Blocks:[]*Block{GenesisBlock()},
	}
}

// 遍历输出链中的内容
func (obj *BlockChain)GetAllBlockHash(){
	for i, v := range obj.Blocks{
		fmt.Println("now height ", i)
		fmt.Printf("pre block hash: %x\n",v.PreHash)
		fmt.Printf("this block hash: %x\n",v.Hash)
	}
}

func (obj *BlockChain)AddBlock(data string, dif, nonce uint64)bool{
    // 获取上一个区块的hash
	preHash := obj.Blocks[len(obj.Blocks) - 1].Hash
	blk := NewBlock(preHash, data, dif, nonce)
	istrue := true
	// 判断区块中唯一
    // 比特币系统中不需要，因为比特币系统中并不保存自身的hash
	for _, v := range obj.Blocks{
		if hex.EncodeToString(v.Hash) == hex.EncodeToString(blk.Hash){
			istrue = false
			return istrue
		}
	}
	obj.Blocks = append(obj.Blocks, blk)
	return istrue
}
