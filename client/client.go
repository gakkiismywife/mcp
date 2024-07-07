package client

import (
	"errors"
	"fmt"
	"mcp/protocol"
	"net"
)

type Mc3EClient struct {
	Host     string   //plc地址
	Port     int      //plc端口
	NetWork  string   //网络类型 udp或者tcp
	DataType DataFlag //交互传输数据类型 binary ascii
	Handler  protocol.DataTypeHandler
}

type DataFlag int

const (
	Binary DataFlag = 1
	Ascii  DataFlag = 2
)

func NewClient(host string, port int, netWork string, dataType DataFlag) *Mc3EClient {
	if err := checkNetWork(netWork); err != nil {
		panic(fmt.Errorf("NewClient: %w", err))
	}

	if err := checkDataType(dataType); err != nil {
		panic(fmt.Errorf("NewClient: %w", err))
	}
	return &Mc3EClient{
		Host:     host,
		Port:     port,
		NetWork:  netWork,
		DataType: dataType,
		Handler:  NewDataTypeHandlerByDataFlag(dataType),
	}
}

func checkNetWork(network string) error {
	if network == "tcp" {
		return nil
	}
	if network == "udp" {
		return nil
	}
	return errors.New("not support network")
}

func checkDataType(datatype DataFlag) error {
	if datatype == Binary {
		return nil
	}
	if datatype == Ascii {
		return nil
	}
	return errors.New("not support datatype")
}

func NewDataTypeHandlerByDataFlag(d DataFlag) protocol.DataTypeHandler {
	switch d {
	case Binary:
		return protocol.Binary{}
	case Ascii:
		return protocol.ASCII{}
	}
	return nil
}

func (c *Mc3EClient) Read(name string, addr, num int64) ([]byte, error) {
	//构建请求
	request, err := c.Handler.BuildReadRequest(name, addr, num)
	if err != nil {
		return nil, err
	}

	//建立连接
	conn, err := net.Dial(c.NetWork, fmt.Sprintf("%s:%d", c.Host, c.Port))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	//发送读取请求
	_, err = conn.Write(request)
	if err != nil {
		return nil, err
	}

	//读取响应
	readBuff := make([]byte, 22+2*num) // 22 is response header size. [sub header + network num + unit i/o num + unit station num + response length + response code]
	readLen, err := conn.Read(readBuff)
	if err != nil {
		return nil, err
	}

	return readBuff[:readLen], nil
}

func (c *Mc3EClient) Write(name string, addr, num int64, data []byte) ([]byte, error) {
	//构建写入请求
	request, err := c.Handler.BuildWriteRequest(name, addr, num, data)
	if err != nil {
		return nil, err
	}

	//建立连接
	conn, err := net.Dial(c.NetWork, c.Addr())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	//发送请求
	if _, err = conn.Write(request); err != nil {
		return nil, err
	}

	//读取响应
	readBuff := make([]byte, 22)
	readLen, err := conn.Read(readBuff)
	if err != nil {
		return nil, err
	}
	return readBuff[:readLen], nil
}

// Addr 拼接地址
func (c *Mc3EClient) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
