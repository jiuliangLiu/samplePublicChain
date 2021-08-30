/**
 * @Author:LJL
 * @Date: 2021/4/17 21:04
 */
package common

import (
	"bytes"
	"encoding/binary"
	"log"
)

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err!=nil{
		log.Panic(err)
	}
	return buff.Bytes()
}