package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"bytes"
	"encoding/gob"
)


type BlockChain struct {
	ChainName string
	// Blocks []*Block
}

// 生成创世块
func (obj *BlockChain)GenesisBlock(){
	str := "first block"
	err := NewBlock([]byte{}, str, 3, 1, obj.ChainName)
	if err != nil{
		log.Println("create first block wrong ")
	}
}

func NewBlockChain(ChainName string)*BlockChain{
	istrue := CheckRepeat(ChainName)
	if istrue{
		return &BlockChain{
			ChainName:ChainName,
		}
	} else {
		return nil
	}
}

// 遍历输出链中的内容
func (obj *BlockChain)GetAllBlockHash(){
	fmt.Printf("show %s blockchain blocks\n", obj.ChainName)
	last := GetLastBlockHash(obj.ChainName)
	blk := obj.GetBlock(last)
	blk.ShowBlock()
	for {
		if len(blk.PreHash) == 0{
			break
		}
		blk = obj.GetBlock(blk.PreHash)
		blk.ShowBlock()
	}
	fmt.Println("show all block info seccuss!!! ")
}

func GetLastBlockHash(ChainName string)[]byte {
	db, err := bolt.Open("block.db", 0600, nil)
	defer db.Close()
	if err != nil{
		log.Println("open blt db wrong: ", err)
	}
	lastHashtmp := []byte{}
	db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(ChainName))
		if bkt == nil{
			log.Println("this Chain is not exists ", ChainName)
			return err
		}
		lastHashtmp = bkt.Get([]byte("lastHash"))
		return nil
	})
	lastHash := make([]byte, len(lastHashtmp), len(lastHashtmp))
	copy(lastHash, lastHashtmp)
	return lastHash
}


// add block
func (obj *BlockChain)AddBlock(data string, dif, nonce uint64){
    // 获取上一个区块的hash
	preHash := GetLastBlockHash(obj.ChainName)
	blk := NewBlock(preHash, data, dif, nonce, obj.ChainName)
	pw := blk.NewPoW()
	fmt.Printf("target: %x\nfound Hash: %x\n", pw.target, blk.Hash)
	/*
	// 判断区块中唯一
    // 比特币系统中不需要，因为比特币系统中并不保存自身的hash
	for _, v := range obj.Blocks{
		if hex.EncodeToString(v.Hash) == hex.EncodeToString(blk.Hash){
			istrue = false
			return istrue
		}
	}
	*/
}

func (obj *BlockChain)GetBlock(hs []byte) *Block {
	db, err := bolt.Open("block.db", 0600, nil)
	defer db.Close()
	if err != nil{
		log.Println("open bolt db wrong: ", err)
	}
	BlkObjtmp := []byte{}
	db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(obj.ChainName))
		if bkt == nil{
			log.Println("this Chain is not exists ", obj.ChainName)
			return err
		}
		BlkObjtmp = bkt.Get(hs)
		return nil
	})
	BlkObj := make([]byte, len(BlkObjtmp), len(BlkObjtmp))
	copy(BlkObj, BlkObjtmp)
	blk := ToObj(BlkObj)
	return blk
}

func ToObj(objByte []byte) *Block {
	blk := new(Block)
	var buf bytes.Buffer
	_, err := buf.Write(objByte)
	if err != nil{
		log.Println("write objinfo to buffer wrong", err)
	}
	decoder := gob.NewDecoder(&buf)
	err = decoder.Decode(blk)
	if err != nil{
		log.Println("unmarshal to obj wrong ", err)
		return nil
	}
	return blk
}

func CheckRepeat(cName string) bool {
	db, err := bolt.Open("block.db", 0600, nil)
	defer db.Close()
	if err != nil{
		log.Println("open bolt db wrong: ", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(cName))
		if bkt != nil{
			log.Println("this Chain is not exists ", cName)
			return err
		}
		return nil
	})
	if err != nil{
		return false
	} else {
		return true
	}
}

