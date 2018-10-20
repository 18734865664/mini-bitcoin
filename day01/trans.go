package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"fmt"
	"crypto/ecdsa"
	"golang.org/x/crypto/ripemd160"
	"github.com/btcsuite/btcutil/base58"
	"time"
)

// 挖矿奖励
var Reward = 12.5

// 交易结构体定义
type Tx struct {
	TxId []byte     // 交易ID, 这里就是对交易的Hash, 比特币会复杂一些
	Inputs []InPut    // 输入集合
	Outputs []OutPut   // 输出集合
}

type InPut struct {
	TxId []byte
	VoutIdx int64     // 可用的余额在输入交易中的索引
	//ScriptSig string    // 锁定tag
	PubKey ecdsa.PublicKey
	SignInfo []byte
}

type OutPut struct {
	Count float64        // 输出的额度
	PubKeyHash []byte
	//ScriptPb string      // 解锁tag
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
func CoinBaseTx(address string, data string)Tx  {
	if data == ""{
		data = fmt.Sprintf("reward %s %f\n", address, Reward)
	}
	ipt := InPut{nil, -1, ecdsa.PublicKey{}, []byte(time.Now().Format("2006-01-02 15:04:05"))}
	opt := OutPut{Reward, GetPubKeyHash(address)}
	tx := Tx{[]byte("helloworld"), []InPut{ipt},[]OutPut{opt} }
	tx.SetTxId()
	return tx
}

// 普通交易创建
func NewTx(ipts []InPut, change OutPut, target string)Tx {
	tx := Tx{}
	tx.Inputs = ipts
	return tx
}

func GetPubKeyHash(addr string)[]byte {
	AddrHash := base58.Decode(addr)
	PubHash := AddrHash[1: len(AddrHash) - 4]
	return PubHash
}

func PubKeyToHash(PubKey *ecdsa.PublicKey)[]byte {
	pubKeyHash := []byte{}
	pubKeyHash = append(pubKeyHash, PubKey.X.Bytes()...)
	pubKeyHash = append(pubKeyHash, PubKey.Y.Bytes()...)
	rip := ripemd160.New()
	_, err := rip.Write(pubKeyHash)
	if err !=nil{
		log.Println("trans.go line 80: ", err)
		return nil
	}
	ripHash := rip.Sum(nil)
	pubKeyHash = ripHash[:20]
	return pubKeyHash
}

func (obj *InPut)ToHash()[]byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil{
		log.Println("<ToHash> trans.go: ", err)
		return []byte{}
	}
	return buf.Bytes()
}

type SignInfo struct {
	FromAddr string
	ToAddr string
	Count float64
}

func (obj *SignInfo)ToHash()[]byte  {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil{
		log.Println("<ToHash> trans.go: ", err)
		return []byte{}
	}
	return buf.Bytes()
}

func Sign(pubKey *ecdsa.PublicKey, signdata []byte )[]byte {
	addr := GetAddress(pubKey)
	signInfo := EccSign(GetPriFromFile(addr), signdata)
	return signInfo
}

func Verify(ipt *InPut, signData []byte)bool {
	istrue := EccVerify(&ipt.PubKey, signData, ipt.SignInfo)
	return istrue
}

