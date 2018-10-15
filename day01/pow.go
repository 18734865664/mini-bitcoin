package main

import (
	"math/big"
	"encoding/hex"
	"fmt"
)

type PoW struct {
	target []byte
	Blk *Block
}

func (obj *PoW)Try()[]byte{
	nonce := 0
	tempHash := []byte{}
	for {
		obj.Blk.Nonce = uint64(nonce)
		tempHash = obj.Blk.SetHash1()
		tempHex := hex.EncodeToString(tempHash)
		var tempBig big.Int
		tempBigString, _:= tempBig.SetString(tempHex, 16)
		fmt.Println(tempBigString)
		temptarget,_:= tempBig.SetString(hex.EncodeToString(obj.target), 16)
		fmt.Println("==============================")
		fmt.Println(temptarget)
		result := temptarget.Cmp(tempBigString)
		if result ==  1 {
			break
		}
		nonce ++
		fmt.Println(nonce)
	}
	return tempHash
}

func (obj *PoW)Check() bool {
	tempHash := obj.Blk.SetHash1()
	var tempInt big.Int
	tempHashInt, _ := tempInt.SetString(hex.EncodeToString(tempHash), 16)
	temptarget,_:= tempInt.SetString(hex.EncodeToString(obj.target), 16)
	if temptarget.Cmp(tempHashInt) == 1{
		return true
	}
	return false
}