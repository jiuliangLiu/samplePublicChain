/**
 * @Author:LJL
 * @Date: 2021/4/19 10:03
 */
package core

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockchainIterator struct {
	CurrentHash []byte
	DB          *bolt.DB
}

//迭代方法
func (blockchainIterator *BlockchainIterator) Next() *Block {
	var block *Block

	err := blockchainIterator.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			//根据区块哈希获取当前区块字节数组
			currentBlockBytes := b.Get(blockchainIterator.CurrentHash)
			//反序列化为区块
			block = DeserializeBLock(currentBlockBytes)
			//迭代器中的区块哈希改为上一个区块哈希
			blockchainIterator.CurrentHash = block.PrevBlockHash
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	return block
}
