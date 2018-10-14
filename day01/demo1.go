package main

import (
	"crypto/sha256"
	"fmt"
	"encoding/hex"
)

// defind block struct
type Block struct {
	PreHash []byte
	Hash []byte
	Data []byte
}

type BlockChain struct {
	Blocks []*Block
}

func main() {
	blkC := NewBlockChain()
	str := "首先是它的树的结构,默克尔树常见的结构是二叉"
	blkC.AddBlock(str)
	blkC.GetAllBlockHash()
}

// get a block obj point
func NewBlock(prehash []byte,data string)*Block{
	var blk Block
    blk.PreHash = prehash
	blk.Data = []byte(data)
	blk.Hash = blk.SetHash()
	return &blk
}

// bitcoin  hash  is SHA256
func (obj *Block)SetHash()[]byte{
	h := sha256.New()
	h.Write(obj.PreHash)
	h.Write([]byte(obj.Data))
	return h.Sum(nil)
}

func NewBlockChain()*BlockChain{
	return &BlockChain{
		Blocks:[]*Block{GenesisBlock()},
	}
}

func GenesisBlock()*Block{
	str := "first block"
	return NewBlock([]byte{}, str)
}

// 遍历输出链中的内容
func (obj *BlockChain)GetAllBlockHash(){
	for i, v := range obj.Blocks{
		fmt.Println("now height ", i)
		fmt.Printf("pre block hash: %x\n",v.PreHash)
		fmt.Printf("this block hash: %x\n",v.Hash)
	}
}

func (obj *BlockChain)AddBlock(data string)bool{
    // 获取上一个区块的hash
	preHash := obj.Blocks[len(obj.Blocks) - 1].Hash
	blk := NewBlock(preHash, data)
	istrue := true
	for _, v := range obj.Blocks{
		if hex.EncodeToString(v.Hash) == hex.EncodeToString(blk.Hash){
			istrue = false
			return istrue
		}
	}
	obj.Blocks = append(obj.Blocks, blk)
	return istrue
}
