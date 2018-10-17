package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"fmt"
)

var Reward = 12.5

type Tx struct {
	TxId []byte
	Inputs []*InPut
	Outputs []*OutPut
}

type InPut struct {
	TxId []byte
	VoutIdx int64
	ScriptSig string
}

type OutPut struct {
	Count float64
	ScriptPb string
}

func (obj *Tx)SetTxId() {

	obj.TxId = obj.ToHash()
}

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

