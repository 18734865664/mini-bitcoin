package main

import (
	"github.com/boltdb/bolt"
	"log"
)

/*
    创建一个迭代器用于循环操作
*/

type BCIter struct {
	ChainName string
	NowBlock []byte
}


// 迭代去方法
func (obj *BCIter)Next() *Block{
	if len(obj.NowBlock) == 0{
		return nil
	}
	db, err := bolt.Open("block.db", 0600, nil)
	if err != nil{
		log.Println("BCIter Next open db error ", err)
		return nil
	}
	defer db.Close()
	objByte := []byte{}
	db.Update(func(tx *bolt.Tx) error {
		blk := tx.Bucket([]byte(obj.ChainName))
		if blk == nil{
			log.Println("BCIter open block wrong")
			return nil
		}
		objByte = blk.Get([]byte(obj.NowBlock))
		return nil
	})
	blk := new(Block)
	blk = ToObj(objByte)
	obj.NowBlock = blk.PreHash
	return blk
}
