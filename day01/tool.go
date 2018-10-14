package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

func Uint64ToByte(i uint64)[]byte{
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, i)
	if err != nil{
		log.Println("change type uint64 to byte: ", err)
		return []byte{}
	}
	return buf.Bytes()
}
