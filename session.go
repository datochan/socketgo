package socketgo

import (
	"fmt"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// IPacketProtocol 此接口定义了如何组包,解包,发送
// TODO: 如果要正常发送需要实现此接口
type IPacketProtocol interface {
	// ReadPacket 读取封包
	ReadPacket(net.Conn) (interface{}, error)
	// BuildPacket 组包
	BuildPacket(interface{}) []byte
	// SendPacket 发送数据
	SendPacket(net.Conn, []byte) error
}

// FnCallbackClosed 链接本身的处理句柄
type FnCallbackClosed func(net.Conn)

type FnCallbackSended func(net.Conn, interface{})

type ISession interface {
	RawConn() net.Conn
	Start()
	Send(packet interface{}) error
	Close() error
	SetCloseCallback(callback FnCallbackClosed)
	SetSendCallback(callback FnCallbackSended)
}

// Session 异步会话管理
type Session struct {
	lock sync.Mutex

	conn          net.Conn
	protocol      IPacketProtocol
	packetHandler PacketHandler
	closeCallback func(net.Conn)
	sendCallback  func(net.Conn, interface{})

	closed int32 // session是否关闭，-1未开启，0未关闭，1关闭

	sendChan   chan interface{} // 发送管道
	stopedChan chan interface{}
}

func NewSession(conn net.Conn, protocol IPacketProtocol, handler PacketHandler, sendChanSize int) *Session {
	return &Session{
		conn:          conn,
		protocol:      protocol,
		packetHandler: handler,
		closed:        -1,
		stopedChan:    make(chan interface{}),
		sendChan:      make(chan interface{}, sendChanSize),
	}
}

// RawConn return net.Conn
func (s *Session) RawConn() net.Conn {
	return s.conn.(*net.TCPConn)
}

// Close 关闭连接并释放相关资源.
func (s *Session) Close() error {
	if atomic.CompareAndSwapInt32(&s.closed, 0, 1) {
		_ = s.conn.Close()
		close(s.stopedChan)

		if s.closeCallback != nil {
			s.closeCallback(s.conn)
		}
	}

	return nil
}

func (s *Session) SetCloseCallback(callback FnCallbackClosed) {
	s.closeCallback = callback
}

func (s *Session) SetSendCallback(callback FnCallbackSended) {
	s.sendCallback = callback
}

// SetReadDeadline 设置读取的超时时间
// goroutine safe, 如果不需要设置，则不要调用
func (s *Session) SetReadDeadline(delt time.Duration) {
	s.lock.Lock()
	_ = s.conn.SetReadDeadline(time.Now().Add(delt)) // timeout
	defer s.lock.Unlock()
}

func (s *Session) sendLoop() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("panic recover! p: %+v", p)
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			fmt.Printf("%s\n", string(buf))
		}

		_ = s.Close()
	}()

	var err error

	for {
		select {
		case <-s.stopedChan:
			{
				fmt.Println("exit,exit,exit")
				return
			}
		case packet, ok := <-s.sendChan:
			{
				if !ok {
					fmt.Printf("发送循环已退出...\n")
					return
				}

				pkgcnt := s.protocol.BuildPacket(packet)
				err = s.protocol.SendPacket(s.conn, pkgcnt)

				if err != nil {
					fmt.Printf("发送循环已退出, 错误信息为:%v...\n", err)
					return
				}

				if s.sendCallback != nil {
					s.sendCallback(s.conn, packet)
				}
			}
		}
	}
}

func (s *Session) recvLoop() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("panic recover! p: %+v", p)
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			fmt.Printf("%s\n", string(buf))
		}
		_ = s.Close()
	}()

	for {
		select {
		case <-s.stopedChan:
			return
		default:
			{
				recvBuff, err := s.protocol.ReadPacket(s.conn)
				if recvBuff == nil || nil != err {
					fmt.Printf("Read packet error %+v", err)
					return
				}

				s.packetHandler(s, recvBuff) // 任务封包分发
			}
		}
	}
}

// Start 开始会话，循环监听发送与接收
func (s *Session) Start() {
	if atomic.CompareAndSwapInt32(&s.closed, -1, 0) {
		go s.sendLoop()
		go s.recvLoop()
	}
}

// Send 异步发送方法, 仅将 packet 写入 sendChan 中等待sendLoop处理,
// 如果 sendChan 满了, 则return ErrSendChanBlocking。
// 如果 sendChan 已关闭, 则return ErrSessionClosed。
func (s *Session) Send(packet interface{}) error {
	select {
	case s.sendChan <- packet:
	case <-s.stopedChan:
		return ErrSessionClosed
	default:
		return ErrSendChanBlocking
	}
	return nil
}
