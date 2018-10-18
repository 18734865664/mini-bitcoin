package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"bytes"
	"encoding/gob"
	"strings"
	"strconv"
	"crypto/ecdsa"
)


type BlockChain struct {
	ChainName string
	// Blocks []*Block
}

// 生成创世块
func (obj *BlockChain)GenesisBlock(addr string){
	coinbase := CoinBaseTx(addr, "first block")
	err := NewBlock([]byte{}, []Tx{coinbase}, 1, 1, obj.ChainName)
	if err == nil{
		log.Println("create first block wrong ")
	}
}

// 新创建一个链（独立的链）
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

// 获取指定区块链的最后一个区块hash
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
	// 这里需要注意，如果不是make出来的变量，会导致最后的结果位数不正确
	lastHash := make([]byte, len(lastHashtmp), len(lastHashtmp))
	// 一定要使用copy进行深拷贝，浅拷贝，会由于db句柄关闭，导致底层数组被回收，引起内存错误
	copy(lastHash, lastHashtmp)
	if err != nil {
		return []byte{}
	} else {
		return lastHash
	}
}


// add block
func (obj *BlockChain)AddBlock(txs []Tx, dif, nonce uint64){
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

// 根据区块的hash 获取区块对象
func (obj *BlockChain)GetBlock(hs []byte) *Block {
	db, err := bolt.Open("block.db", 0600, nil)
	defer db.Close()
	if err != nil{
		log.Println("open bolt db wrong: ", err)
	}
	BlkObjtmp := []byte{}
	// 更新bolt数据库
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

// 将数据库中读取出的内容转变未block 对象
func ToObj(objByte []byte) *Block {
    // 需要一个结构体对象接收解码后的对象
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

// 检查区块链名是否重复
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

// 遍历账本，返回utxo的count 和 在交易中的idx  map，以交易ID作为key
func (obj *BlockChain)GetAllUTXO(addr string) (map[string][]int64, map[string][]float64){
	Iutxs := map[string][]int64{}
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
		if addr == "" {
			return nil, nil
		} else {
			for _, tx := range blk.Txs{
				// 匹配交易输入集合，将其中输出地址为addr的加入ntxo已使用集合
				for _, itx := range tx.Inputs{
					if bytes.Equal(itx.GetPubKeyHash(), GetPubKeyHash(addr)){
						if Iutxs[string(itx.TxId)] != nil{
							Iutxs[string(itx.TxId)] = append(Iutxs[string(itx.TxId)], itx.VoutIdx)
						} else {
							Iutxs[string(itx.TxId)] = []int64{itx.VoutIdx}
						}
					}
				}
				// 匹配交易输出集合，将其中输出地址为addr的加入总的ntxo集合
				OUTPUT:
				for idx, otx := range tx.Outputs{
					if bytes.Equal(otx.PubKeyHash,GetPubKeyHash(addr)) {
						if Iutxs[string(tx.TxId)] == nil{
							for _, v := range Iutxs[string(tx.TxId)]{
								if v == int64(idx){
									// 如果判断出输出已经被消费，则跳过后续操作
									continue OUTPUT
								}
							}
							OutxsMap[string(tx.TxId)] = append(OutxsMap[string(tx.TxId)], int64(idx))
							OutRemark[string(tx.TxId)] = append(OutRemark[string(tx.TxId)], otx.Count)
						}
					}
				}
			}
		}
	}
	/*  这种方式多一次遍历，会更麻烦
		之前对utxo理解有误,交易的单个输出是utxo的单位，而不是交易
    // 遍历已使用ntxo集合，将其从总集合总去掉
	for _, itx := range Iutxs {
		if OutxsMap[itx] != nil{
			delete(OutxsMap, itx)
			delete(OutRemark, itx)
		}
	}
	*/
	return OutxsMap, OutRemark
}

// 返回账户总余额
func (obj *BlockChain)GetCountCoin(addr string, outCount map[string][]float64)float64{
	count := 0.0
	for _, otx := range outCount{
		for _, c := range otx{
			count += c
		}
	}
	return count
}

// 输出账户信息
func (obj *BlockChain)ShowCountCoin(addr string) {
	_, outCount := obj.GetAllUTXO(addr)
	count := obj.GetCountCoin(addr, outCount)
	fmt.Printf("BlockChain:\t%s\naddress:\t%s\ncount:\t\t%f\n", obj.ChainName, addr, count)
}

// 创建交易时，返回合适的Input集合
func (obj *BlockChain)FindProperUtxo(addr string,PriKey *ecdsa.PrivateKey, request float64)([]InPut,float64){
	ipt := []InPut{}
	utxos, outCount := obj.GetAllUTXO(addr)
	count := obj.GetCountCoin(addr, outCount)
	PubKey := PriKey.PublicKey
	counttmp := 0.0
// 使用的逻辑很简单，总的Input coin 总量大于等于 需求的coin
	if request > count {
		log.Println("request is larger than count")
		return []InPut{}, 0.0
	}
	_, utxoRemarks := obj.GetAllUTXO(addr)
	for id, cs := range utxoRemarks{
		if counttmp >=  request{
			break
		}
		for idx, c := range cs{
			counttmp += c
			input := InPut{
				TxId:[]byte(id),
				VoutIdx: utxos[id][idx],
				//ScriptSig:adx,
				PubKey:PubKey,
			}
			input.SignInfo = input.SignInfo
			ipt = append(ipt, input)
		}
	}
	return ipt, counttmp - request
}

// 根据余额创建找零交易
func (obj *BlockChain)GetChangeOutPut(addr string, looseChange float64) OutPut{
	opt := OutPut{}
	// addr := obj.GetNewAddress()
	opt.Count = looseChange
	opt.PubKeyHash = GetPubKeyHash(addr)
	return opt
}


// 生成地址值
func (obj *BlockChain)GetNewAddress()string {
	acc := GetNewAccount()
	w := Wallet{}
	w.WName = obj.ChainName + ".wal"
	accs := w.GetAllAccount()
	// 将新生成的地址写入钱包
	accs[acc.Addr] = acc
	w.Addrs = accs
	w.SaveToFile()
	return acc.Addr
}

// 返回区块链对象的迭代器, 有一个Next函数，用于返回父区块对象
func (obj *BlockChain)NewBCIter()*BCIter {
	bct := BCIter{
		ChainName:obj.ChainName,
	}
	return &bct
}

// 创建建议对象
func (obj *BlockChain)NewTx(addr, target string)Tx {
	tx := Tx{}
	opts:= []OutPut{}
	pKey := GetPriFromFile(addr)
	if pKey  == nil{
		log.Fatal("<NewTx> bloclchain.go")
	}
	request := 0.0
	for _, v :=range  strings.Split(target,","){
		args := strings.Split(v, ":")
		count, _ := strconv.ParseFloat(args[1], 10)
		request += count
		opt := OutPut{
			Count:count,
			//ScriptPb:args[0],
			PubKeyHash: GetPubKeyHash(target),
		}
		opts = append(opts, opt)
	}
    // 寻求输入集合，返回零钱值
	ipts, looseChange := obj.FindProperUtxo(addr,pKey, request)
	opt := obj.GetChangeOutPut(addr, looseChange)
	opts = append(opts, opt)
	tx.Outputs = opts
	tx.Inputs = ipts
	tx.SetTxId()
	return tx
}

// 创建一个普通交易
func (obj *BlockChain)CreateCommTrans(addr ,target string) Tx{
	// txs := []*Tx{}
	tx := obj.NewTx(addr, target)
	return tx
	// txs = append(txs, tx)
	// obj.AddBlock(txs, 0, 0)
}

// 显示整个钱包中的地址信息
func (obj *BlockChain)ShowWallet()  {
	w := &Wallet{}
	w.WName = obj.ChainName + ".wal"
	accs := w.GetAllAccount()
	fmt.Printf("show %s wallet address: \n\n", obj.ChainName)
	for k, _ := range accs {
		fmt.Printf("\taddress: %s\n", k)
	}
}

func (obj *BlockChain)CheckInputPub(ipt *InPut)bool  {
	addr := GetAddress(ipt.PubKey)
	idxs, _ := obj.GetAllUTXO(addr)
	if idxs[addr] != nil{
		for _, v := range idxs[addr]{
			if ipt.VoutIdx == v{
				isSign := ipt.Verify()
				if isSign {
					return true
				}
				fmt.Printf("this Input is tampered: %x -- %d", ipt.TxId, ipt.VoutIdx)
				return false
			}
		}
	}
	fmt.Println("checkInputPub wrong: %x : %d", ipt.TxId, ipt.VoutIdx)
	return false
}
