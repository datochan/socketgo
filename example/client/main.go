package main

import (
	"fmt"
	socket "github.com/datochan/socketgo"
	"github.com/datochan/socketgo/example/proto"
	goproto "google.golang.org/protobuf/proto"
	"time"
)

type ExampleAsyncClient struct {
	*socket.AsyncClient
	Ctx          map[string]interface{} // 用于在命令行和处理线程之间传递上下文信息
	FinishedChan chan interface{}       // 用于异步通讯转同步使用，recv线程处理完信息后将结果存放于Ctx，并通过FinishedChan通知发送消息继续下一步

}

func NewExampleAsyncClient() *ExampleAsyncClient {
	protocol := NewExampleProtocolImpl()

	return &ExampleAsyncClient{
		AsyncClient: socket.NewAsyncClient(protocol, NewExampleDispatcher(), 10),
		Ctx:         make(map[string]interface{}),
	}
}

func (c *ExampleAsyncClient) Hearbeat() *proto.RequestNode {
	hearbeat := &proto.HearBeat{HearBeatId: 1}
	byCnt, _ := goproto.Marshal(hearbeat)

	return proto.NewRequestNode(0x02, byCnt)
}

func (c *ExampleAsyncClient) OnHearbeat(session socket.ISession, packet interface{}) {
	hearbeat := proto.HearBeat{}
	respNode := packet.(proto.ResponseNode)

	_ = goproto.Unmarshal(respNode.Data.([]byte), &hearbeat)

	fmt.Printf("收到心跳应答封包, 内容为: %d\n", hearbeat.HearBeatId)
}

func (c *ExampleAsyncClient) OnBroadcast(session socket.ISession, packet interface{}) {
	broadcast := proto.Broadcast{}
	respNode := packet.(proto.ResponseNode)

	_ = goproto.Unmarshal(respNode.Data.([]byte), &broadcast)

	fmt.Printf("收到定时广播, 内容为: %s\n", broadcast.Content)
}

func main() {
	client := NewExampleAsyncClient()
	_ = client.Conn("tcp", "127.0.0.1:7190", 64*1024, 64*1024,
		nil, nil)

	heartbeat := client.Hearbeat()
	client.GetDispatcher().AddHandler(heartbeat.Flag, client.OnHearbeat)
	client.GetDispatcher().AddHandler(0x03, client.OnBroadcast)

	for {
		time.Sleep(3 * time.Second)
		fmt.Printf("发送心跳封包....\n")
		err := client.Send(heartbeat)
		if err != nil {
			fmt.Printf("发送出错, Err: %+v", err)
			return
		}
	}
}
