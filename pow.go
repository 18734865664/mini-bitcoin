package main

import (
	"math/big"
)

type PoW struct {
	target *big.Int
	Blk *Block
}

func (obj *PoW)Try()[]byte{
	nonce := 0
	tempHash := []byte{}
	for {
		obj.Blk.Nonce = uint64(nonce)
		tempHash = obj.Blk.SetHash1()
		tempBig := new(big.Int)
		tempBig = tempBig.SetBytes(tempHash)
		result := tempBig.Cmp(obj.target)
		// fmt.Println(nonce)
		// fmt.Println(obj.target)
		// fmt.Println(tempBig)
		if result ==  -1 {
			break
		}
		nonce ++
	}
	return tempHash
}

func (obj *PoW)Check() bool {
	tempHash := obj.Blk.SetHash1()
	var tempInt big.Int
	tempInt.SetBytes(tempHash)
	if obj.target.Cmp(&tempInt) == 1{
		return true
	}
	return false
}