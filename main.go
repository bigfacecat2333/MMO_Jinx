package main

import (
	"Jinx/jinterface"
	"Jinx/jnet"
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

	fmt.Println("======> player pid = ", player.PID, " is online <======")
}

func main() {
	// 创建jinx的句柄
	s := jnet.NewServer()
	fmt.Println(s)

	// 链接创建和销毁的hook函数
	s.SetOnConnStart(OoConnectionAdd)

	// 注册路由
	// 启动服务
	s.Serve()
}
