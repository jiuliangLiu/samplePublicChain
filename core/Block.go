/**
 * @Author:LJL
 * @Date: 2021/4/17 19:53
 *	定义区块结构
 */
package core

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Height        int64  //区块高度
	PrevBlockHash []byte //上一个区块哈希
	Data          []byte //交易数据
	Timestamp     int64  //时间戳
	Hash          []byte //区块哈希
	Nonce         int64  //PoW nonce
}

func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

func DeserializeBLock(blockBytes []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil{
		log.Panic(err)
	}
	return &block
}

//生成创世区块
func CreateGenesisBlock(data string) *Block {
	return NewBlock(1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, data)
}

//创建新区块
func NewBlock(height int64, prevBlockHash []byte, data string) *Block {
	// 构造区块
	block := &Block{height, prevBlockHash, []byte(data), time.Now().Unix(), nil, 0}
	//调用工作量证明方法，获取Hash和Nonce
	pow := NewPoW(block)

	hash, nonce := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}
