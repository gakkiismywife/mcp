package client

import (
	"errors"
	"fmt"
	"mcp/protocol"
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
