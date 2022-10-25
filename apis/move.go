package apis

import (
	"Jinx/jinterface"
	"Jinx/jnet"
	"MMO/core"
	"MMO/pb"
	"fmt"
	"google.golang.org/protobuf/proto"
)

type MoveApi struct {
	jnet.BaseRouter
}

func (*MoveApi) Handle(request jinterface.IRequest) {
	//1. 将客户端传来的proto协议解码
	protoMsg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), protoMsg)
	if err != nil {
		fmt.Println("Move Unmarshal error ", err)
		return
	}
	//2. 得知当前的消息是从哪个玩家传递来的,从连接属性pid中获取
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid error", err)
		request.GetConnection().Stop()
		return
	}
	fmt.Println("Move pid = ", pid, " x = ", protoMsg.X, " y = ", protoMsg.Y, " z = ", protoMsg.Z, " V = ", protoMsg.V)
	//3. 给其他玩家广播当前玩家的位置信息
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	// 根据得到的玩家对象,将位置信息同步给周边玩家 (广播)
	player.UpdatePos(protoMsg.X, protoMsg.Y, protoMsg.Z, protoMsg.V)
}
