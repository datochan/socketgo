package main

import (
	"fmt"
	socket "github.com/datochan/socketgo"
	"github.com/datochan/socketgo/example/proto"
	goproto "google.golang.org/protobuf/proto"
	"time"
)

type ExampleServer struct {
	*socket.Server
}

func NewExampleServer() (*ExampleServer, error) {
	protocol := NewExampleProtocolImpl()
	server, err := socket.NewServer("tcp", "127.0.0.1:7190", protocol, NewExampleDispatcher())

	if err != nil {
		return nil, err
	}

	return &ExampleServer{server}, nil
}

func (c *ExampleServer) Broadcast() *proto.ResponseNode {
	broadcast := &proto.Broadcast{Content: "这是一条每30秒重复一次的广播信息..."}
	byCnt, _ := goproto.Marshal(broadcast)

	return proto.NewResponseNode(0x03, byCnt)
}

func (c *ExampleServer) Hearbeat() *proto.ResponseNode {
	hearbeat := &proto.HearBeat{HearBeatId: 2}
	byCnt, _ := goproto.Marshal(hearbeat)

	return proto.NewResponseNode(0x02, byCnt)
}

func (c *ExampleServer) OnHearbeat(session socket.ISession, packet interface{}) {
	var hearbeat = &proto.HearBeat{}
	respNode := packet.(proto.RequestNode)

	_ = goproto.Unmarshal(respNode.Data.([]byte), hearbeat)

	fmt.Printf("收到心跳封包, 内容为: %d, 正在应答\n", hearbeat.HearBeatId)
	_ = session.Send(c.Hearbeat())
}

func main() {
	server, err := NewExampleServer()
	if err != nil {
		fmt.Printf("NewExampleServer error %+v\n", err)
		return
	}

	heartbeat := server.Hearbeat()
	server.GetDispatcher().AddHandler(heartbeat.ReqFlag, server.OnHearbeat)

	go func() {
		for {
			time.Sleep(30 * time.Second)
			for i, s := range server.SessionMng {
				fmt.Printf("向第 %d 个客户端广播通知...\n", i+1)
				_ = s.Send(server.Broadcast())
			}
		}
	}()

	server.AcceptLoop()
}
