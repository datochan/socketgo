package socketgo

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	once       *sync.Once
	listener   net.Listener
	dispatcher IDispatcher
	stopedChan chan struct{}
	protocol   IPacketProtocol
	SessionMng []ISession
}

// NewServer 新建服务器
func NewServer(network, address string, protocol IPacketProtocol, dispatcher IDispatcher) (*Server, error) {
	listener, err := net.Listen(network, address)
	if err != nil {
		fmt.Println("listen failed !!!")
		return nil, ErrListenFailed
	}

	return &Server{
		dispatcher: dispatcher,
		listener:   listener,
		once:       &sync.Once{},
		protocol:   protocol,
		stopedChan: make(chan struct{}),
	}, nil
}

// GetDispatcher 获取事件分发器
func (s *Server) GetDispatcher() IDispatcher {
	return s.dispatcher
}

func (s *Server) Close() {
	s.once.Do(func() {
		if s.listener != nil {
			_ = s.listener.Close()
		}

		close(s.stopedChan)
	})
}

func (s *Server) acceptLoop() error {
	tcpConn, err := s.listener.Accept()
	if err != nil {
		fmt.Println("Accept error!")
		return ErrAcceptFailed
	}

	session := NewSession(tcpConn, s.protocol, s.dispatcher.HandleProc, 0)

	fmt.Println("A client connected :" + tcpConn.RemoteAddr().String())
	s.SessionMng = append(s.SessionMng, session)
	session.Start()

	return nil
}

func (s *Server) AcceptLoop() {
	for {
		if err := s.acceptLoop(); err != nil {
			fmt.Println("server accept failed!!")
			return
		}
	}
}
