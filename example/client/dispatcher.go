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
	respNode := packet.(proto.ResponseNode)

	fmt.Printf("收到未知封包: RespFlag=%X, ReqFlag=%X, len=%d\n; content= %s\n",
		respNode.RespFlag, respNode.ReqFlag, respNode.BodyLength, hex.EncodeToString(respNode.Data.([]byte)))
}

// HandleProc 事件处理过程
func (p *ExampleDispatcher) HandleProc(session socket.ISession, packet interface{}) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	respNode := packet.(proto.ResponseNode)
	handlerProc := p.GetHandler(respNode.ReqFlag)
	if nil == handlerProc {
		UnknownPkgHandler(session, packet)
		return
	}

	handlerProc(session, packet)
}
