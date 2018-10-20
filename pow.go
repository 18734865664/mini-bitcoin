package main

import (
	"math/big"
	"encoding/hex"
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
		tempHex := hex.EncodeToString(tempHash)
		var tempBig big.Int
		tempBigString, _:= tempBig.SetString(tempHex, 16)
		result := obj.target.Cmp(tempBigString)
		if result >  0 {
			break
		}
		nonce ++
	}
	return tempHash
}

func Check()  {

}