package protocol

import (
	"fmt"
)

type ASCII struct {
}

func (A ASCII) BuildHealthRequest() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (A ASCII) BuildReadRequest(name string, offset, num int64) ([]byte, error) {
	deviceCode := "D*"

	//offsetBuff := new(bytes.Buffer)
	//_ = binary.Write(offsetBuff, binary.BigEndian, offset)
	//offsetHex := fmt.Sprintf("%X", offsetBuff.Bytes()[offsetBuff.Len()-3:])
	//地址固定占3字节 所以用字符串表示就是例如D1090 在协议中就是D*001090 每两个字符占一个字节
	//length=3时就是 %06d生成6个字符的字符串
	offsetStr := ConvertToStr(offset, 3)
	// 读取字节数
	//pointsBuff := new(bytes.Buffer)
	//_ = binary.Write(pointsBuff, binary.BigEndian, num)
	//points := fmt.Sprintf("%X", pointsBuff.Bytes()[pointsBuff.Len()-2:])

	pointStr := ConvertToStr(num, 2)

	// data length 字节数 两个字符占用一个字节 所以除以2
	//requestCharLen := len(MONITORING_TIMER+ASCII_READ_COMMAND+READ_SUB_COMMAND+deviceCode+offsetHex+points) / 2
	//dataLenBuff := new(bytes.Buffer)
	//_ = binary.Write(dataLenBuff, binary.BigEndian, int64(requestCharLen))
	//dataLen := fmt.Sprintf("%X", dataLenBuff.Bytes()[dataLenBuff.Len()-2:]) // 2byte固定
	//字符串长度
	strLen := len(MONITORING_TIMER + ASCII_READ_COMMAND + READ_SUB_COMMAND + deviceCode + offsetStr + pointStr)
	//转成4个字符长度的16进制
	dataLen := fmt.Sprintf("%04X", strLen)

	//公共部分+数据长度+超时时间+指令+子指令+数据区+起始地址+读取的字节数
	str := A.Common() + dataLen + MONITORING_TIMER + ASCII_READ_COMMAND + READ_SUB_COMMAND + deviceCode + offsetStr + pointStr
	fmt.Println(str)
	return []byte(str), nil
}

func (A ASCII) BuildWriteRequest(name string, offset, num int64, writeData []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func ConvertToStr(target int64, length int64) string {

	format := fmt.Sprintf("%%0%d", length*2)
	fmt.Println(format)
	return fmt.Sprintf(fmt.Sprintf("%sd", format), target)
}

func (A ASCII) Common() string {
	return SubHeader + NetworkNum + IONum + PLCNum + ModuleNum
}
