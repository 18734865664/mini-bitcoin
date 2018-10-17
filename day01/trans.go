package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"fmt"
)

// 挖矿奖励
var Reward = 12.5

// 交易结构体定义
type Tx struct {
	TxId []byte     // 交易ID, 这里就是对交易的Hash, 比特币会复杂一些
	Inputs []*InPut    // 输入集合
	Outputs []*OutPut   // 输出集合
}

type InPut struct {
	TxId []byte
	VoutIdx int64     // 可用的余额在输入交易中的索引
	ScriptSig string    // 锁定tag
}

type OutPut struct {
	Count float64        // 输出的额度
	ScriptPb string      // 解锁tag
}

func (obj *Tx)SetTxId() {
	obj.TxId = obj.ToHash()
}

// 将交易进行序列化, 后获取hash值
func (obj *Tx)ToHash()[]byte  {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil{
		log.Fatal("tx obj trans to hash wrong: ", err)
	}
	objHash := sha256.Sum256(buf.Bytes())
	return objHash[:]
}

// 挖矿的交易
func CoinBaseTx(address string, data string)*Tx  {
	if data == ""{
		data = fmt.Sprintf("reward %s %f\n", address, Reward)
	}
	ipt := InPut{nil, -1, data}
	opt := OutPut{Reward, address}
	tx := Tx{[]byte{}, []*InPut{&ipt},[]*OutPut{&opt} }
	tx.SetTxId()
	return &tx
}

// 普通交易创建
func NewTx(ipts []*InPut, change *OutPut, target string)*Tx {
	tx := &Tx{}
	tx.Inputs = ipts
	return tx
}

