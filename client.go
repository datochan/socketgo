package socketgo

import (
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Client struct {
	conn       net.Conn
	protocol   IPacketProtocol
	stopedChan <-chan os.Signal
}

func NewClient(protocol IPacketProtocol) *Client {
	stopSignal := make(chan os.Signal) // 接收系统中断信号
	var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}
	signal.Notify(stopSignal, shutdownSignals...)

	return &Client{
		stopedChan: stopSignal,
		protocol:   protocol,
	}
}

func (c *Client) Conn(network, address string, readBufferSize, writeBufferSize int) error {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil
	}

	c.conn = conn

	tcpConn := conn.(*net.TCPConn)
	_ = tcpConn.SetNoDelay(true)
	_ = tcpConn.SetReadBuffer(readBufferSize)
	_ = tcpConn.SetWriteBuffer(writeBufferSize)

	return nil
}

func (c *Client) Send(packet interface{}) error {
	select {
	case <-c.stopedChan:
		return ErrSignalStopped
	default:
		content := c.protocol.BuildPacket(packet)
		err := c.protocol.SendPacket(c.conn, content)
		if err != nil {
			return ErrWritePacketFailed
		}
	}

	return nil
}

func (c *Client) Recv() (interface{}, error) {
	packet, err := c.protocol.ReadPacket(c.conn)
	if err != nil {
		return nil, ErrReadPacketFailed
	}

	return packet, nil
}

// AsyncClient 异步通讯客户端
type AsyncClient struct {
	Client
	sendChanSize int
	session      ISession
	dispatcher   IDispatcher
}

func NewAsyncClient(protocol IPacketProtocol, dispatcher IDispatcher, bufferSize int) *AsyncClient {
	stopSignal := make(chan os.Signal) // 接收系统中断信号
	var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}
	signal.Notify(stopSignal, shutdownSignals...)

	return &AsyncClient{
		Client: Client{
			stopedChan: stopSignal,
			protocol:   protocol,
		},
		sendChanSize: bufferSize,
		dispatcher:   dispatcher,
	}
}

// GetDispatcher 获取事件分发器
func (c *AsyncClient) GetDispatcher() IDispatcher {
	return c.dispatcher
}

// Close 关闭连接
func (c *AsyncClient) Close() {
	if nil != c.session {
		_ = c.session.Close()
	}
}

// Conn 与服务器建立TCP连接
// :Param network: 网络类型: "tcp"、"udp"
// :Param address: 链接的地址: "IP:PORT"
func (c *AsyncClient) Conn(network, address string,
	readBufferSize, writeBufferSize int,
	callbackSend FnCallbackSended, callbackClosed FnCallbackClosed) error {

	err := c.Client.Conn(network, address, readBufferSize, writeBufferSize)
	if err != nil {
		return nil
	}

	c.session = NewSession(c.conn, c.protocol, c.dispatcher.HandleProc, c.sendChanSize)

	if callbackSend != nil {
		c.session.SetSendCallback(callbackSend)
	}

	if callbackClosed != nil {
		c.session.SetCloseCallback(callbackClosed)
	}

	c.session.Start()

	return nil
}

func (c *AsyncClient) Send(packet interface{}) error {
	return c.session.Send(packet)
}

// Recv 异步通讯的客户端只能注册接收消息的句柄，不能直接收取封包内容
func (c *AsyncClient) Recv() (interface{}, error) {
	panic("异步通讯客户端不允许直接读取封包内容")
	return nil, ErrReadPacketFailed
}
