package main

import (
	"crypto/sha256"
	"time"
	"bytes"
	"strings"
	"github.com/boltdb/bolt"
	"log"
	"encoding/gob"
	"fmt"
	"math/big"
	"strconv"
)

// defind block struct
type Block struct {
	PreHash []byte

// 正常比特币区块中没有区块的HASH，这里是为了方便操作作弊
	Hash []byte
	Data []byte

	Version uint64

// 默克尔根
	MerkelRoot []byte
	Timestamp uint64

// 挖矿难度系数
	Difficulty uint64

// 挖矿随机数
	Nonce uint64

// blockName
	ChainName string
}
// get a block obj point
func NewBlock(prehash []byte,data string, dif,nonce uint64, cName string)*Block{
	var blk Block
    blk.PreHash = prehash
	blk.Data = []byte(data)
	blk.MerkelRoot = []byte{}
	blk.Version = uint64(1)
	blk.Timestamp = uint64(time.Now().Unix())
	blk.Difficulty = uint64(dif)
	blk.Nonce = nonce
	blk.ChainName = cName
	// blk.Hash = blk.SetHash1()
	blk.Hash = blk.NewPoW().Try()
	err := blk.AddBlockToDb()
	if err != nil{
		fmt.Println("new block", err)
	}
	return &blk
}

// bitcoin  hash  is SHA256
func (obj *Block)SetHash()[]byte{
	h := sha256.New()
	h.Write(obj.PreHash)
	h.Write([]byte(obj.ChainName))
	h.Write([]byte(obj.Data))
	h.Write(Uint64ToByte(obj.Difficulty))
	h.Write(Uint64ToByte(obj.Timestamp))
	h.Write(Uint64ToByte(obj.Nonce))
	h.Write(Uint64ToByte(obj.Version))
	h.Write(obj.MerkelRoot)
	return h.Sum(nil)
}

func (obj *Block)SetHash1()[]byte {
    // 使用bytes.Join 函数
	blkInfo := [][]byte{
		obj.PreHash,
		[]byte(obj.ChainName),
		[]byte(obj.Data),
		Uint64ToByte(obj.Difficulty),
		Uint64ToByte(obj.Timestamp),
		Uint64ToByte(obj.Nonce),
		Uint64ToByte(obj.Version),
		obj.MerkelRoot,
	}
	tmp := bytes.Join(blkInfo, []byte{})
	sh := sha256.Sum256(tmp)
	return  sh[:]
}

/*
// 生成创世块
func GenesisBlock()*Block{
	str := "first block"
	return NewBlock([]byte{}, str, 3, 1)
}
*/

// new pow obj
func (obj *Block)NewPoW()*PoW{
     // sha256 转换成16进制是64位， 所以这里的字符串需要64位
    dif := int(obj.Difficulty)
    tgtbyte := bytes.Join([][]byte{
    	bytes.Repeat([]byte{0}, dif),
    	[]byte{1},
    	bytes.Repeat([]byte{0}, 32 - dif - 1),
	},[]byte{})
	pow := PoW{
		Blk:obj,
	}

    // 将目标hash 字符串转换成了big.Int类型，指明了进制为16进制
	// pow.target , _ = temp.SetString(tgtmpstr, 16)
	tempInt := new(big.Int)
	tempInt = tempInt.SetBytes(tgtbyte)
	pow.target = tempInt
	return &pow
}

func (obj *Block)AddBlockToDb()error  {
	db, err  := bolt.Open("block.db", 0600, nil)
	if err != nil{
		log.Println("block.go 108: ", err)
		return err
	}
	defer db.Close()
	//lastHashTmp := []byte{}
	err = db.Update(func(tx *bolt.Tx) error {
		blk := tx.Bucket([]byte(obj.ChainName))
		if blk == nil{
			blk, err = tx.CreateBucket([]byte(obj.ChainName))
			if err != nil{
				return err
			}
		}
		err = blk.Put([]byte(obj.Hash), obj.ToHash())
		err = blk.Put([]byte("lastHash"), obj.Hash)
		//lastHashTmp = blk.Get([]byte(obj.ChainName))
		return nil
	})
	//lastHash := make([]byte, len(lastHashTmp), len(lastHashTmp))
	//copy(lastHash, lastHashTmp)
	//obj1 := Dec(lastHash)
	//fmt.Println("obj unmarshal", obj1)
	return err
}

func (obj *Block)ToHash()[]byte  {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil{
		log.Println("ToHash ", err)
		return nil
	}
	return buf.Bytes()
}

func Dec(blockInfo []byte) *Block {
	blk := new(Block)
	var buf bytes.Buffer
	dec := gob.NewDecoder(&buf)
	buf.Write(blockInfo)
	dec.Decode(blk)
	return blk
}

func (obj *Block)ShowBlock()  {
	str := "blk Hash: %x\nblk Nonce: %s\nblk PreHash: %x\n"
	fmt.Printf(str, obj.Hash, strconv.Itoa(int(obj.Nonce)), obj.PreHash)
	fmt.Println(strings.Repeat("+", 20))
}
