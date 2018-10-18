package main

import (
	"crypto/rsa"
	"crypto/rand"
	"fmt"
	"log"
	"crypto/x509"
	"encoding/pem"
	"os"
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
	"crypto"
)

const KeysBit = 1024

func GetRsaPrivateKey()*rsa.PrivateKey {
	pkey, err  := rsa.GenerateKey(rand.Reader, KeysBit)
	if err != nil{
		log.Println("<keys.go 15line>create pravateKey wrong: ", err)
		return  nil
	}
	return pkey
}

func RsaSign(pkey *rsa.PrivateKey, hashed []byte)[]byte {
	hash := crypto.SHA256
	signInfo, err := rsa.SignPKCS1v15(rand.Reader, pkey,hash,hashed)
	if err != nil{
		fmt.Println("keys.go 31 ", err )
		return nil
	}
	return signInfo
}

func RsaVerify(pub *rsa.PublicKey,hashed []byte, sig []byte)  bool{
	hash := crypto.SHA256
	err := rsa.VerifyPKCS1v15(pub, hash, hashed, sig)
	if err != nil{
		return false
	}
	return true
}

func PriToFile(pkey *rsa.PrivateKey, fName string)  {
	pkeyHash := x509.MarshalPKCS1PrivateKey(pkey)
	block := pem.Block{
		Type: "pravate key",
		Headers:nil,
		Bytes:pkeyHash,
	}
	f,err  := os.Create(fName)
	if err != nil{
		fmt.Println("keys 30: ", err)
		return
	}
	err = pem.Encode(f, &block)
	if err != nil{
		fmt.Println("pem encode wrong ", err )
		return
	}
}

func GetEccPrivateKey()*ecdsa.PrivateKey {
	curve := elliptic.P256()
	pkey, err  := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil{
		log.Println("<keys.go GetEccPrivateKey 50>: ",err)
		return  nil
	}
	return pkey
}

func EccSign(pkey *ecdsa.PrivateKey, hash []byte)[]byte {
	r, s , err := ecdsa.Sign(rand.Reader, pkey, hash)
	if err != nil{
		fmt.Println("ecc sign wrong: ", err)
		return []byte{}
	}
	tmp := []byte{}
	tmp = append(tmp, r.Bytes()...)
	tmp = append(tmp, s.Bytes()...)
	return tmp
}

func EccVerify(pub *ecdsa.PublicKey, hash []byte, signInfo []byte) bool {
	rByte := signInfo[:len(signInfo)/2]
	sByte := signInfo[len(signInfo)/2:]
	var r, s big.Int
	r.SetBytes(rByte)
	s.SetBytes(sByte)
	istrue := ecdsa.Verify(pub, hash, &r, &s)
	return istrue
}


