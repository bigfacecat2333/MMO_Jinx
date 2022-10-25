package main

import (
	"Jinx/jinterface"
	"Jinx/jnet"
	"MMO/apis"
	"MMO/core"
	"fmt"
)

// OoConnectionAdd 当前客户端建立链接后的hook函数
func OoConnectionAdd(conn jinterface.IConnection) {
	// 创建一个玩家
	player := core.NewPlayer(conn)

	// 给客户端发送MsgId:1的消息（同步）
	player.SyncPid()

	// 给客户端发送MsgId:200的消息（广播）,同步当前玩家的位置信息
	player.BroadCastStartPosition()

	// 将当前玩家添加到世界管理器中
	core.WorldMgrObj.AddPlayer(player)

	// 将当前玩家的链接绑定到conn属性中
	conn.SetProperty("pid", player.Pid)

	// 同步周边玩家的信息, 告知当前玩家，当前世界有哪些玩家
	player.SyncSurrounding()

	fmt.Println("======> player pid = ", player.Pid, " is online <======")
}

func main() {
	// 创建jinx的句柄
	s := jnet.NewServer()
	fmt.Println(s)

	// 链接创建和销毁的hook函数
	s.SetOnConnStart(OoConnectionAdd)

	// 注册路由
	s.AddRouter(2, &apis.WorldChatApi{})

	// 启动服务
	s.Serve()
}
