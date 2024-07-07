package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type Binary struct {
	CommonHeader
}

const (
	SUB_HEADER = "5000" // 3Eフレームでは固定

	BINARY_HEALTH_CHECK_COMMAND = "1906" // binary mode expression. if ascii mode then 0619
	ASCII_HEALTH_CHECK_COMMAND  = "1906" // binary mode expression. if ascii mode then 0619
	HEALTH_CHECK_SUBCOMMAND     = "0000"

	BINARY_READ_COMMAND  = "0104" // binary mode expression. if ascii mode then 0401
	ASCII_READ_COMMAND   = "0401" // binary mode expression. if ascii mode then 0401
	READ_SUB_COMMAND     = "0000"
	BIT_READ_SUB_COMMAND = "0100"

	BINARY_WRITE_COMMAND = "0114" // binary mode expression. if ascii mode then 1401
	ASCII_WRITE_COMMAND  = "1401" // binary mode expression. if ascii mode then 1401
	WRITE_SUB_COMMAND    = "0000"

	MONITORING_TIMER = "1000" // 3[sec]
)

// deviceCodes is device name and hex value map
var deviceCodes = map[string]string{
	"X": "9C",
	"Y": "9D",
	"M": "90",
	"L": "92",
	"F": "93",
	"V": "94",
	"B": "A0",
	"W": "B4",
	"D": "A8",
}

func (b Binary) BuildHealthRequest() ([]byte, error) {
	return nil, nil
}

func (b Binary) BuildReadRequest(name string, offset, num int64) ([]byte, error) {
	// get device symbol hex layout
	deviceCode := deviceCodes[name]

	// offset convert to little endian layout
	// MELSECコミュニケーションプロトコル リファレンス(p67) MELSEC-Q/L: 3[byte], MELSEC iQ-R: 4[byte]
	offsetBuff := new(bytes.Buffer)
	_ = binary.Write(offsetBuff, binary.LittleEndian, offset)
	offsetHex := fmt.Sprintf("%X", offsetBuff.Bytes()[0:3]) // 仮にQシリーズとするので3byte trim

	// read points
	pointsBuff := new(bytes.Buffer)
	_ = binary.Write(pointsBuff, binary.LittleEndian, num)
	points := fmt.Sprintf("%X", pointsBuff.Bytes()[0:2]) // 2byte固定

	// data length 字节数 两个字符占用一个字节 所以除以2
	requestCharLen := len(MONITORING_TIMER+BINARY_READ_COMMAND+READ_SUB_COMMAND+deviceCode+offsetHex+points) / 2 // 1byte=2char
	dataLenBuff := new(bytes.Buffer)
	_ = binary.Write(dataLenBuff, binary.LittleEndian, int64(requestCharLen))
	dataLen := fmt.Sprintf("%X", dataLenBuff.Bytes()[0:2]) // 2byte固定

	//公共部分+数据长度+超时时间+指令+子指令+数据区+起始地址+读取的字节数
	str := b.Common() + dataLen + MONITORING_TIMER + BINARY_READ_COMMAND + READ_SUB_COMMAND + offsetHex + deviceCode + points
	fmt.Println(str)
	return hex.DecodeString(str)
}

func (b Binary) BuildWriteRequest(name string, offset, num int64, writeData []byte) ([]byte, error) {
	// get device symbol hex layout
	deviceCode := deviceCodes[name]

	// offset convert to little endian layout
	// MELSECコミュニケーションプロトコル リファレンス(p67) MELSEC-Q/L: 3[byte], MELSEC iQ-R: 4[byte]
	offsetBuff := new(bytes.Buffer)
	_ = binary.Write(offsetBuff, binary.LittleEndian, offset)
	//数据区用3字节标识 小端表示法
	offsetHex := fmt.Sprintf("%X", offsetBuff.Bytes()[0:3]) // 仮にQシリーズとするので3byte trim

	// convert write data to little endian word
	writeBuff := new(bytes.Buffer)
	_ = binary.Write(writeBuff, binary.LittleEndian, writeData)
	writeHex := fmt.Sprintf("%X", writeBuff.Bytes()[0:2*num]) // 2 byte per 1 device point

	// write points
	pointsBuff := new(bytes.Buffer)
	_ = binary.Write(pointsBuff, binary.LittleEndian, num)
	//写入字节数固定2个字节表示
	points := fmt.Sprintf("%X", pointsBuff.Bytes()[0:2]) // 2byte固定

	// data length
	requestCharLen := len(MONITORING_TIMER+BINARY_WRITE_COMMAND+WRITE_SUB_COMMAND+deviceCode+offsetHex+points+writeHex) / 2 // 1byte=2char
	dataLenBuff := new(bytes.Buffer)
	_ = binary.Write(dataLenBuff, binary.LittleEndian, int64(requestCharLen))
	dataLen := fmt.Sprintf("%X", dataLenBuff.Bytes()[0:2]) // 2byte固定
	//return SUB_HEADER +
	//	h.networkNum +
	//	h.pcNum +
	//	h.unitIONum +
	//	h.unitStationNum +
	//	dataLen +
	//	MONITORING_TIMER +
	//	WRITE_COMMAND +
	//	WRITE_SUB_COMMAND +
	//	offsetHex +
	//	deviceCode +
	//	points +
	//	writeHex
	str := b.Common() + dataLen + MONITORING_TIMER + BINARY_WRITE_COMMAND + WRITE_SUB_COMMAND + offsetHex + deviceCode + points + writeHex
	fmt.Println(str)
	return hex.DecodeString(str)
}
