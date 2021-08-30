/**
 * @Author:LJL
 * @Date: 2021/4/17 20:29
 */
package core

import (
	"blockchain/publicChain/common"
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//为0的位数
const targetBit = 20

type PoW struct {
	Block  *Block
	target *big.Int
}

func (pow *PoW) prepareData(nonce int64) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevBlockHash,
			pow.Block.Data,
			common.IntToHex(pow.Block.Timestamp),
			common.IntToHex(int64(targetBit)),
			common.IntToHex(nonce),
			common.IntToHex(pow.Block.Height),
		},
		[]byte{},
	)
	return data
}

func (poW *PoW) IsValid() (result bool) {
	var hashInt big.Int
	hashInt.SetBytes(poW.Block.Hash)
	if poW.target.Cmp(&hashInt) == 1 {
		result = true
	} else {
		result = false
	}
	return
}

func (poW *PoW) Run() (hash []byte, nonce int64) {
	nonce = 0
	var hashInt big.Int

	for {
		dataBytes := poW.prepareData(nonce)
		hash32 := sha256.Sum256(dataBytes)
		hash = hash32[:]
		hashInt.SetBytes(hash)
		//如果hashInt小于目标，求解成功，跳出循环
		if poW.target.Cmp(&hashInt) == 1 {
			fmt.Println(" target:", poW.target)
			fmt.Println("hashInt:", &hashInt)
			break
		}

		nonce = nonce + 1
	}
	return
}

//创建pow对象
func NewPoW(block *Block) *PoW {
	target := big.NewInt(1)
	//左移256-targetBit位
	target = target.Lsh(target, 256-targetBit)
	return &PoW{block, target}
}
