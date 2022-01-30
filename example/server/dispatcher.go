package main

import (
	"encoding/hex"
	"fmt"
	"sync"

	socket "github.com/datochan/socketgo"
	"github.com/datochan/socketgo/example/proto"
)

type ExampleDispatcher struct {
	*socket.Dispatcher
	lock sync.RWMutex // 写互斥避免并发状态下相互干扰
}

// NewExampleDispatcher 事件分发器
func NewExampleDispatcher() *ExampleDispatcher {
	return &ExampleDispatcher{Dispatcher: socket.NewDispatcher()}
}

func UnknownPkgHandler(session socket.ISession, packet interface{}) {
	requestNode := packet.(proto.RequestNode)

	fmt.Printf("收到未知封包: ReqFlag=%X, len=%d\n; content= %s\n", requestNode.Flag,
		requestNode.BodyLength, hex.EncodeToString(requestNode.Data.([]byte)))
}

// HandleProc 事件处理过程
func (p *ExampleDispatcher) HandleProc(session socket.ISession, packet interface{}) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	requestNode := packet.(proto.RequestNode)
	handlerProc := p.GetHandler(requestNode.Flag)
	if nil == handlerProc {
		UnknownPkgHandler(session, packet)
		return
	}

	handlerProc(session, packet)
}
