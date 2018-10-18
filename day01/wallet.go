package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"os"
	"fmt"
	"io"
	"encoding/gob"
	"bytes"
	"github.com/btcsuite/btcutil/base58"
	"crypto/elliptic"
	"log"
)

func init()  {
	// 非基本类型，需要显示的注册才能进行解码
	gob.Register(elliptic.P256())
}

type account struct {
	PriKey ecdsa.PrivateKey
	PubKey ecdsa.PublicKey
	Addr string
}

type Wallet struct {
	Addrs map[string]*account
	WName string
}

func GetAddress(pkey ecdsa.PublicKey)string  {
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
	AllPart  = append(AllPart, sPart[:][:4]...)

	// addr := sha256.Sum256([]byte(time.Now().String()))
	// 使用base58编码，TODO
	b58 := base58.Encode(AllPart)
	return b58

}


func GetNewAccount()*account  {
	// publickey rip160 hash
    // 获取私钥信息
	pkey := GetEccPrivateKey()
	b58 := GetAddress(pkey.PublicKey)

	act := account{
		PriKey: *pkey,
		PubKey: pkey.PublicKey,
		Addr: b58,
	}
	PriToFile(pkey, act.Addr)
	return &act
}

// 将钱包数据写入文件
func (obj *Wallet)SaveToFile()  {
	f, err := os.OpenFile(obj.WName, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil{
		fmt.Println("wallet.go 66: ", err)
		return
	}
   // 执行序列化程序
	wHash := obj.ToHash()
	_, err = f.Write(wHash)
	if err !=nil{
		fmt.Println("wallet.go line 76: ", err )
	}
}

// 获取钱包中的账户信息map
func (obj *Wallet)GetAllAccount()map[string]*account  {
	f, err := os.OpenFile(obj.WName, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil{
		fmt.Println("wallet.go 83: ", err)
		return map[string]*account{}
	}
	data := []byte{}
	buf := make([]byte, 4 * 1024)
	for {
		n, err := f.Read(buf)
		data = append(data, buf[:n]...)
		if err == io.EOF{
			break
		} else if err != nil {
			fmt.Println("wallet.go 94: ", err)
			break
		}
	}
	// 将读取到的字节码解码
	w := ToWalletObj(data)
	// 如果钱包中没有内容，w.Addrs是nil,返回后会报错，所以进行一下处理
	if w.Addrs == nil{
		return map[string]*account{}
	}
	return w.Addrs
}

func (obj *account)ToHash()[]byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil{
		fmt.Println("wallet line 104: ",err)
		return nil
	}
	return buf.Bytes()
}

// 将对向序列化为字节码
func (obj *Wallet)ToHash()[]byte  {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil{
		fmt.Println("wallet line 114: ", err)
		return nil
	}
	return buf.Bytes()
}

func ToWalletObj(data []byte)*Wallet  {
	var buf bytes.Buffer
    // 这里需要执行new方法 为对象分配内存，否则解码时将无法成功解码
	w := new(Wallet)
	buf.Write(data)
	dec := gob.NewDecoder(&buf)
	err := dec.Decode(w)
	if err != nil{
		log.Println("wallet.go line 131: ", err)
	}
	return w
}


