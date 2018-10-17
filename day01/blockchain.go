package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"bytes"
	"encoding/gob"
	"time"
	"crypto/sha256"
	"strings"
	"strconv"
)


type BlockChain struct {
	ChainName string
	// Blocks []*Block
}

// 生成创世块
func (obj *BlockChain)GenesisBlock(addr string){
	coinbase := CoinBaseTx(addr, "first block")
	err := NewBlock([]byte{}, []*Tx{coinbase}, 1, 1, obj.ChainName)
	if err == nil{
		log.Println("create first block wrong ")
	}
}

func NewBlockChain(ChainName string)*BlockChain{
	istrue := CheckRepeat(ChainName)
	if istrue{
		return &BlockChain{
			ChainName:ChainName,
		}
	} else {
		return nil
	}
}

// 遍历输出链中的内容
func (obj *BlockChain)GetAllBlockHash(){
	fmt.Printf("show %s blockchain blocks\n", obj.ChainName)
	last := obj.GetLastBlockHash()
	if len(last) == 0{
		return
	}
	blk := obj.GetBlock(last)
	/*
	blk.ShowBlock()
	for {
		if len(blk.PreHash) == 0{
			break
		}
		blk = obj.GetBlock(blk.PreHash)
		blk.ShowBlock()
	}
	*/
	// user iterator
	iter := obj.NewBCIter()
	iter.NowBlock = last
	for {
		blk = iter.Next()
		if blk == nil{
			break
		}
		blk.ShowBlock()
	}
	fmt.Println("show all block info seccuss!!! ")
}

func (obj *BlockChain)GetLastBlockHash()[]byte {
	db, err := bolt.Open("block.db", 0600, nil)
	defer db.Close()
	if err != nil{
		log.Println("open blt db wrong: ", err)
	}
	lastHashtmp := []byte{}
	err = db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(obj.ChainName))
		if bkt == nil{
			log.Println("this Chain is not exists ", obj.ChainName)
			return err
		}
		lastHashtmp = bkt.Get([]byte("lastHash"))
		return nil
	})
	lastHash := make([]byte, len(lastHashtmp), len(lastHashtmp))
	copy(lastHash, lastHashtmp)
	if err != nil {
		return []byte{}
	} else {
		return lastHash
	}
}


// add block
func (obj *BlockChain)AddBlock(txs []*Tx, dif, nonce uint64){
    // 获取上一个区块的hash
	preHash := obj.GetLastBlockHash()
	blk := NewBlock(preHash, txs, dif, nonce, obj.ChainName)
	pw := blk.NewPoW()
	fmt.Printf("target: %x\nfound Hash: %x\n", pw.target, blk.Hash)
	/*
	// 判断区块中唯一
    // 比特币系统中不需要，因为比特币系统中并不保存自身的hash
	for _, v := range obj.Blocks{
		if hex.EncodeToString(v.Hash) == hex.EncodeToString(blk.Hash){
			istrue = false
			return istrue
		}
	}
	*/
}

func (obj *BlockChain)GetBlock(hs []byte) *Block {
	db, err := bolt.Open("block.db", 0600, nil)
	defer db.Close()
	if err != nil{
		log.Println("open bolt db wrong: ", err)
	}
	BlkObjtmp := []byte{}
	db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(obj.ChainName))
		if bkt == nil{
			log.Println("this Chain is not exists ", obj.ChainName)
			return err
		}
		BlkObjtmp = bkt.Get(hs)
		return nil
	})
	BlkObj := make([]byte, len(BlkObjtmp), len(BlkObjtmp))
	copy(BlkObj, BlkObjtmp)
	blk := ToObj(BlkObj)
	return blk
}

func ToObj(objByte []byte) *Block {
	blk := new(Block)
	var buf bytes.Buffer
	_, err := buf.Write(objByte)
	if err != nil{
		log.Println("write objinfo to buffer wrong", err)
	}
	decoder := gob.NewDecoder(&buf)
	err = decoder.Decode(blk)
	if err != nil{
		log.Println("unmarshal to obj wrong ", err)
		return nil
	}
	return blk
}

func CheckRepeat(cName string) bool {
	db, err := bolt.Open("block.db", 0600, nil)
	defer db.Close()
	if err != nil{
		log.Println("open bolt db wrong: ", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(cName))
		if bkt != nil{
			log.Println("this Chain is not exists ", cName)
			return err
		}
		return nil
	})
	if err != nil{
		return false
	} else {
		return true
	}
}

func (obj *BlockChain)GetAllUTXO(adx string) (map[string][]int64, map[string][]float64){
	Iutxs := []string{}
	OutRemark :=  map[string][]float64{}
	OutxsMap :=  map[string][]int64{}

	last := obj.GetLastBlockHash()
	iter := obj.NewBCIter()
	iter.NowBlock = last
	blk := obj.GetBlock(last)
	for {
		blk = iter.Next()
		if blk == nil{
			break
		}
		if adx == "" {
			return nil, nil
		} else {
			for _, tx := range blk.Txs{
				for _, itx := range tx.Inputs{
					if itx.ScriptSig == adx{
						Iutxs = append(Iutxs, string(itx.TxId))
					}
				}
				for idx, otx := range tx.Outputs{
					if otx.ScriptPb == adx{
						OutxsMap[string(tx.TxId)] = append(OutxsMap[string(tx.TxId)], int64(idx))
						OutRemark[string(tx.TxId)] = append(OutRemark[string(tx.TxId)], otx.Count)
					}
				}
			}
		}
	}
	for _, itx := range Iutxs {
		if OutxsMap[itx] != nil{
			delete(OutxsMap, itx)
			delete(OutRemark, itx)
		}
	}
	return OutxsMap, OutRemark
}

func (obj *BlockChain)GetCountCoin(adx string, outCount map[string][]float64)float64{
	count := 0.0
	for _, otx := range outCount{
		for _, c := range otx{
			count += c
		}
	}
	return count
}

func (obj *BlockChain)ShowCountCoin(adx string) {
	_, outCount := obj.GetAllUTXO(adx)
	count := obj.GetCountCoin(adx, outCount)
	fmt.Printf("BlockChain:\t%s\naddress:\t%s\ncount:\t\t%f\n", obj.ChainName, adx, count)
}

func (obj *BlockChain)FindProperUtxo(adx string, request float64)([]*InPut,float64){
	ipt := []*InPut{}
	utxos, outCount := obj.GetAllUTXO(adx)
	count := obj.GetCountCoin(adx, outCount)
	counttmp := 0.0
	if request > count {
		log.Println("request is larger than count")
		return []*InPut{}, 0.0
	}
	_, utxoRemarks := obj.GetAllUTXO(adx)
	for id, cs := range utxoRemarks{
		if counttmp >=  request{
			break
		}
		for idx, c := range cs{
			counttmp += c
			ipt = append(ipt, &InPut{
				TxId:[]byte(id),
				VoutIdx: utxos[id][idx],
				ScriptSig:adx,
			})
		}
	}
	return ipt, counttmp - request
}

func (obj *BlockChain)GetChangeOutPut(adx string, looseChange float64) *OutPut{
	opt := &OutPut{}
	// addr := obj.GetNewAddress()
	opt.Count = looseChange
	opt.ScriptPb = string(adx)
	return opt
}

func (obj *BlockChain)GetNewAddress()[]byte {
	addr := sha256.Sum256([]byte(time.Now().String()))
	return addr[:]
}

func Get()  {

}

func (obj *BlockChain)NewBCIter()*BCIter {
	bct := BCIter{
		ChainName:obj.ChainName,
	}
	return &bct
}


func (obj *BlockChain)NewTx(adx, target string)*Tx {
	tx := &Tx{}
	opts:= []*OutPut{}
	request := 0.0
	for _, v :=range  strings.Split(target,","){
		args := strings.Split(v, ":")
		count, _ := strconv.ParseFloat(args[1], 10)
		request += count
		opt := &OutPut{
			Count:count,
			ScriptPb:args[0],
		}
		opts = append(opts, opt)
	}
	ipts, looseChange := obj.FindProperUtxo(adx, request)
	opt := obj.GetChangeOutPut(adx, looseChange)
	opts = append(opts, opt)
	tx.Outputs = opts
	tx.Inputs = ipts
	tx.SetTxId()
	return tx
}

func (obj *BlockChain)CreateCommTrans(addr ,target string) {
	txs := []*Tx{}
	tx := obj.NewTx(addr, target)
	txs = append(txs, tx)
	obj.AddBlock(txs, 0, 0)
}
