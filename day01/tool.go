package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

// 进行类型转换
func Uint64ToByte(i uint64)[]byte{
	var buf bytes.Buffer
    // 指定了存储的大小端, 现在网络上传输数据通常使用大端模式
    // 小端存储，低位内存存放地位数据
	err := binary.Write(&buf, binary.BigEndian, i)
	if err != nil{
		log.Println("change type uint64 to byte: ", err)
		return []byte{}
	}
	return buf.Bytes()
}
