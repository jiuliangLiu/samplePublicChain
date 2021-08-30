/**
 * @Author:LJL
 * @Date: 2021/4/17 20:07
 */
package core

import (
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/boltdb/bolt"
)

const dbName = "blockchain.db"

const blockTableName = "blocks"

type Blockchain struct {
	Tip []byte //最新的区块哈希
	DB  *bolt.DB
}

func (blc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blc.Tip, blc.DB}
}

//打印输出所有区块信息
func (blc *Blockchain) PrintChain() {
	blockchainIterator := blc.Iterator()
	for {
		block := blockchainIterator.Next()
		fmt.Println("Height:", block.Height)
		fmt.Printf("PrevBlockHash:%x\n", block.PrevBlockHash)
		fmt.Printf("Data:%s\n", string(block.Data))
		fmt.Printf("Timestamp:%d\n", block.Timestamp)
		fmt.Printf("Hash:%x\n", block.Hash)
		fmt.Printf("Nonce:%d\n", block.Nonce)
		fmt.Println("---------------------------")
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
}

func (blc *Blockchain) AddBlockToBlockchain(data string) {

	err := blc.DB.Update(func(tx *bolt.Tx) error {
		//获取表
		b := tx.Bucket([]byte(blockTableName))
		//创建新区块
		if b != nil {
			//根据区块哈希获取最新区块
			blockBytes := b.Get(blc.Tip)
			//反序列化后获得最新区块
			block := DeserializeBLock(blockBytes)
			//构建新区块
			newBlock := NewBlock(block.Height+1, block.Hash, data)
			//将新区块序列化后存入数据库
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			//更新数据库中"l"对应的hash
			err = b.Put([]byte("l"), newBlock.Hash)
			if err != nil {
				log.Panic(err)
			}

			//更新blockchain的Tip
			blc.Tip = newBlock.Hash
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

func CreateBlockchainWithGenesisBlock(data string) *Blockchain {

	if DBExists() {
		fmt.Println("创世区块已经存在")
		os.Exit(1)
	}

	fmt.Println("正在创建创世区块")

	//打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	//用于存储区块哈希
	var blockHash []byte
	//插入数据
	err = db.Update(func(tx *bolt.Tx) error {

		b, err := tx.CreateBucket([]byte(blockTableName))
		if err != nil {
			log.Panic(err)
		}

		if b != nil {
			//生成创世块
			genesisBlock := CreateGenesisBlock(data)
			//将创世区块存储到表中
			err := b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			//存储最新的区块哈希
			err = b.Put([]byte("l"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			blockHash = genesisBlock.Hash
		}
		return nil
	})

	//返回区块链对象
	return nil
}

func BlockchainObject() *Blockchain {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))

		if b != nil {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	return &Blockchain{tip, db}
}
