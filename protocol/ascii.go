package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type ASCII struct {
	CommonHeader
}

func (A ASCII) BuildHealthRequest() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (A ASCII) BuildReadRequest(name string, offset, num int64) ([]byte, error) {
	deviceCode := "D*"

	offsetBuff := new(bytes.Buffer)
	_ = binary.Write(offsetBuff, binary.BigEndian, offset)
	offsetHex := fmt.Sprintf("%X", offsetBuff.Bytes()[offsetBuff.Len()-3:])
	// 读取字节数
	pointsBuff := new(bytes.Buffer)
	_ = binary.Write(pointsBuff, binary.BigEndian, num)
	points := fmt.Sprintf("%X", pointsBuff.Bytes()[pointsBuff.Len()-2:])

	// data length 字节数 两个字符占用一个字节 所以除以2
	requestCharLen := len(MONITORING_TIMER+ASCII_READ_COMMAND+READ_SUB_COMMAND+deviceCode+offsetHex+points) / 2
	dataLenBuff := new(bytes.Buffer)
	_ = binary.Write(dataLenBuff, binary.BigEndian, int64(requestCharLen))
	dataLen := fmt.Sprintf("%X", dataLenBuff.Bytes()[dataLenBuff.Len()-2:]) // 2byte固定

	//公共部分+数据长度+超时时间+指令+子指令+数据区+起始地址+读取的字节数
	str := A.Common() + dataLen + MONITORING_TIMER + ASCII_READ_COMMAND + READ_SUB_COMMAND + deviceCode + offsetHex + points
	fmt.Println(str)
	return hex.DecodeString(str)
}

func (A ASCII) BuildWriteRequest(name string, offset, num int64, writeData []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
