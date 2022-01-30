package socketgo

import (
	"fmt"
	"sync"
)

// PacketHandler 事件处理句柄,用于解析相应的封包
type PacketHandler func(ISession, interface{})

type IDispatcher interface {
	AddHandler(id uint32, handler PacketHandler)
	DelHandler(id uint32)
	GetHandler(id uint32) PacketHandler
	HandleProc(ISession, interface{})
}

type Dispatcher struct {
	rwlock     sync.RWMutex             // 写互斥避免并发状态下相互干扰
	handlerMap map[uint32]PacketHandler // 事件过程回调句柄, key是事件ID
}

// NewDispatcher 事件分发器
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlerMap: make(map[uint32]PacketHandler),
	}
}

// AddHandler 添加新的事件处理器
func (p *Dispatcher) AddHandler(id uint32, handler PacketHandler) {
	p.rwlock.Lock()
	defer p.rwlock.Unlock()
	p.handlerMap[id] = handler
}

// DelHandler 卸载新的事件处理器
func (p *Dispatcher) DelHandler(id uint32) {
	p.rwlock.Lock()
	defer p.rwlock.Unlock()
	delete(p.handlerMap, id)
}

func (p *Dispatcher) GetHandler(id uint32) PacketHandler {
	p.rwlock.Lock()
	defer p.rwlock.Unlock()

	handler, ok := p.handlerMap[id]
	if ok {
		return handler
	}

	return nil
}

// HandleProc 事件处理过程
func (p *Dispatcher) HandleProc(session ISession, packet interface{}) {
	p.rwlock.RLock()
	defer p.rwlock.RUnlock()

	// 子类中实现事件分发功能
	fmt.Println("未实现时间处理过程, 请再子类中实现")
}
