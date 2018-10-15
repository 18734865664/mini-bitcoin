package main

import (
	"crypto/sha256"
	"time"
	"bytes"
	"math/big"
	"strings"
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
}
// get a block obj point
func NewBlock(prehash []byte,data string, dif,nonce uint64)*Block{
	var blk Block
    blk.PreHash = prehash
	blk.Data = []byte(data)
	blk.MerkelRoot = []byte{}
	blk.Version = uint64(1)
	blk.Timestamp = uint64(time.Now().Unix())
	blk.Difficulty = uint64(dif)
	blk.Nonce = nonce
	// blk.Hash = blk.SetHash1()
	blk.Hash = blk.NewPoW().Try()
	return &blk
}

// bitcoin  hash  is SHA256
func (obj *Block)SetHash()[]byte{
	h := sha256.New()
	h.Write(obj.PreHash)
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

// 生成创世块
func GenesisBlock()*Block{
	str := "first block"
	return NewBlock([]byte{}, str, 3, 1)
}

// new pow obj
func (obj *Block)NewPoW()*PoW{
     // sha256 转换成16进制是64位， 所以这里的字符串需要64位
     //
     //
    dif := int(obj.Difficulty)
    tgtmpstr := strings.Join([]string{
		strings.Repeat("0", dif),
		"1",
		strings.Repeat("0", 64 - dif - 1),
	},"")
	pow := PoW{
		Blk:obj,
	}


    // 将目标hash 字符串转换成了big.Int类型，指明了进制为16进制
    var temp big.Int
	pow.target , _ = temp.SetString(tgtmpstr, 16)
	return &pow
}
