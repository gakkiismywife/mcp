package protocol

type DataTypeHandler interface {
	BuildHealthRequest() ([]byte, error)
	BuildReadRequest(name string, offset, num int64) ([]byte, error)
	BuildWriteRequest(name string, offset, num int64, writeData []byte) ([]byte, error)

	Common() string
}

const (
	SubHeader  = "5000" //子头部 固定
	NetworkNum = "00"   //网络编号 固定
	PLCNum     = "FF"   //plc编号 固定
	IONum      = "FF03" //io编号 固定
	ModuleNum  = "00"   //模块编号 固定
)
