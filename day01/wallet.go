package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
)

func init()  {
	//gob.Register()
}

type Account struct {
	PriKey *ecdsa.PrivateKey
	PubKey *ecdsa.PublicKey
	Addr []byte
}

func GetNewAccount()*Account  {
	// publickey rip160 hash
    // 获取私钥信息
	pkey := GetEccPrivateKey()
	pkeyHash := []byte{}
    // 查看公钥源码，公钥中有三部分，其中椭圆曲线相当于常量，所以我们只传输X,Y两个点
	pkeyHash = append(pkeyHash, pkey.X.Bytes()...)//
	pkeyHash = append(pkeyHash, pkey.Y.Bytes()...)//
	rip160 := ripemd160.New()
	rip160.Write(pkeyHash)
	ripKey := rip160.Sum(nil)

	// first part
	vByte := []byte{1}
	pkHash := ripKey[:20]     // 公钥 rip160密文取前20位
	fPart := []byte{}
	fPart = append(fPart, vByte...)
	fPart = append(fPart, pkHash...)

	// second part  sum256加密公钥rip160密文
	sPart := sha256.Sum256(fPart)
	// AllPart  拼接整体
	AllPart := []byte{}
	AllPart  = append(AllPart, fPart...)
	AllPart  = append(AllPart, sPart[:]...)


	act := Account{
		PriKey: pkey,
		PubKey: &pkey.PublicKey,
		Addr: AllPart,
	}
	// addr := sha256.Sum256([]byte(time.Now().String()))
	// 使用base58编码，TODO
	return &act
}

