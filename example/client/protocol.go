package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/Re-volution/sizestruct"
	"github.com/datochan/socketgo/example/proto"
	"net"
)

// CombineBytes 拼接byte数组
func CombineBytes(pBytes ...[]byte) []byte {
	// 将要拼接的数字凑成一个二维数组,通过join进行拼接
	bLen := len(pBytes)
	s := make([][]byte, bLen)
	for index := 0; index < bLen; index++ {
		s[index] = pBytes[index]
	}
	sep := []byte("")

	return bytes.Join(s, sep)
}

// ExampleProtocolImpl 只做最简单实现: IPacketProtocol 接口
type ExampleProtocolImpl struct {
	PacketBuffer *proto.LittleEndianStreamImpl // 封包接收的缓冲区
}

func NewExampleProtocolImpl() *ExampleProtocolImpl {
	return &ExampleProtocolImpl{
		PacketBuffer: proto.NewLittleEndianStream(make([]byte, 1024*5, 1024*5)),
	}
}

// ReadPacket 直接返回所有封包内容，不做任何处理
func (pool *ExampleProtocolImpl) ReadPacket(s net.Conn) (interface{}, error) {
	var newBuffer bytes.Buffer
	var header proto.ResponseHeader

	headerSize := sizestruct.SizeOf(proto.ResponseHeader{})

	tmpBuffer := make([]byte, 1024*5, 1024*5)
	recvLen, err := s.Read(tmpBuffer[:])
	if recvLen <= 0 || err != nil {
		return 0, err
	}

	err = pool.PacketBuffer.WriteBuff(tmpBuffer[:recvLen])
	_ = pool.PacketBuffer.CleanBuff()

	// 预读并解析包头信息
	tmpHeader, _ := pool.PacketBuffer.PeekBuff(headerSize)

	newBuffer.Write(tmpHeader)
	err = binary.Read(&newBuffer, binary.LittleEndian, &header)

	if nil != err {
		return nil, err
	}

	// buffer中已经有完整的封包
	_, _ = pool.PacketBuffer.ReadBuff(headerSize) // 跳过包头

	// 开始读取封包体
	pkgBody, _ := pool.PacketBuffer.ReadBuff(int(header.BodyLength))

	// 清理掉已经处理的数据
	_ = pool.PacketBuffer.CleanBuff()

	//fmt.Printf("--> 收到封包: %s\n", hex.EncodeToString(pkgBody))
	return proto.ResponseNode{ResponseHeader: header, Data: pkgBody}, nil
}

// BuildPacket 组包方法
func (pool *ExampleProtocolImpl) BuildPacket(pkgNode interface{}) []byte {
	requestNode := pkgNode.(proto.IRequestNode)
	byteHeader := requestNode.GenerateHeader()

	return CombineBytes(byteHeader, requestNode.GetRawData().([]byte))
}

func (pool *ExampleProtocolImpl) SendPacket(conn net.Conn, buff []byte) error {
	_, err := conn.Write(buff)

	if err != nil {
		fmt.Printf("发生错误了, 错误信息: %+v", err)
	}
	return err
}
