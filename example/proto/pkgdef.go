package proto

import (
	"bytes"
	"encoding/binary"
	"github.com/Re-volution/sizestruct"
)

// RequestHeader 通讯封包包头结构
type RequestHeader struct {
	Flag       uint32 // 固定为0x0C, 扩展市场固定为01
	BodyLength uint32 // H: =封包长度+2字节(事件标识占用封包的2字节长度)
}

// IRequestNode 请求封包必须要实现生成封包头的接口
type IRequestNode interface {
	GetRawData() interface{} // 获取原始数据
	GenerateHeader() []byte  // 生成封包头(对于行情、交易、用户验证等封包头的要求不一致需要单独生成)
}

// RequestNode 请求封包的基本包装
type RequestNode struct {
	RequestHeader
	Data interface{} // 原始请求封包体的原始数据
}

func NewRequestNode(flag uint32, data interface{}) *RequestNode {
	return &RequestNode{RequestHeader: RequestHeader{Flag: flag}, Data: data}
}

func (r *RequestNode) GetRawData() interface{} {
	return r.Data
}

// GenerateHeader 空实现，具体参考实现类的方法
func (r *RequestNode) GenerateHeader() []byte {
	newBuffer := new(bytes.Buffer)

	length := uint32(sizestruct.SizeOf(r.Data))

	_ = binary.Write(newBuffer, binary.LittleEndian, RequestHeader{
		Flag:       r.Flag,
		BodyLength: length,
	})

	return newBuffer.Bytes()
}

// IResponseNode 请求封包必须要实现生成封包头的接口
type IResponseNode interface {
	GetRawData() interface{} // 获取原始数据
	GenerateHeader() []byte  // 生成封包头(对于行情、交易、用户验证等封包头的要求不一致需要单独生成)
}

// ResponseHeader 应答封包的包头结构
type ResponseHeader struct {
	RespFlag   uint32 // 封包标识
	ReqFlag    uint32 // 请求封包头
	BodyLength uint32 // 分别是封包长度
}

// ResponseNode 响应封包的基本包装
type ResponseNode struct {
	ResponseHeader             // 收到的封包头信息
	Data           interface{} // 原始相应封包体的原始数据
}

func NewResponseNode(flag uint32, data interface{}) *ResponseNode {
	return &ResponseNode{ResponseHeader: ResponseHeader{RespFlag: 0x1001, ReqFlag: flag, BodyLength: 0}, Data: data}
}

// GenerateHeader 简单实现
func (r *ResponseNode) GenerateHeader() []byte {
	newBuffer := new(bytes.Buffer)

	length := uint32(sizestruct.SizeOf(r.Data))

	_ = binary.Write(newBuffer, binary.LittleEndian, ResponseHeader{
		RespFlag:   0x1001,
		ReqFlag:    r.ReqFlag,
		BodyLength: length,
	})

	return newBuffer.Bytes()
}

func (r *ResponseNode) GetRawData() interface{} {
	return r.Data
}
