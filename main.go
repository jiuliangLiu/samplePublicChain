/**
 * @Author:LJL
 * @Date: 2021/4/17 19:47
 */
package main

import (
	"blockchain/publicChain/core"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

func CreateBlockchain(data string) {
	blc := core.CreateBlockchainWithGenesisBlock(data)
	blc.AddBlockToBlockchain("first")
	blc.AddBlockToBlockchain("second")
	blc.PrintChain()
}

func testBoltDb(blc *core.Block) {
	//创建数据库
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//插入或更新数据
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		if b == nil {
			fmt.Println("不存在blocks bucket，正在创建")
			b, err = tx.CreateBucket([]byte("blocks"))
			if err != nil {
				log.Panic("blocks table create faild")
			}
		}

		//err = b.Put([]byte("l"), blc.Serialize())
		err = b.Put([]byte("ll"), []byte("test string"))
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	//查询数据
	err = db.View(func(tx *bolt.Tx) error {
		//获取表对象
		b := tx.Bucket([]byte("blocks"))

		//读取数据
		if b != nil {
			data := b.Get([]byte("l"))
			fmt.Printf("l: %v\n", data)
			blc := core.DeserializeBLock(data)
			fmt.Println("l blc: ", blc)
			data = b.Get([]byte("ll"))
			fmt.Printf("ll: %s\n", data)
		}
		return nil
	})

	//查询失败
	if err != nil {
		log.Panic(err)
	}

}

func main() {
	//CreateBlockchain()
	//blc := core.NewBlock(1, []byte{0}, "test block")
	//fmt.Println("blc =", blc)
	//serialResult := blc.Serialize()
	//fmt.Println("serialResult:", serialResult)
	//newBlc := core.DeserializeBLock(serialResult)
	//fmt.Println("newBlc:", string(newBlc.Data))
	//testBoltDb(blc)

	cli := core.CLI{}
	cli.Run()
}
