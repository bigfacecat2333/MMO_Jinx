package core

import (
	"Jinx/jinterface"
	"MMO/pb"
	"fmt"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"sync"
)

// Player 玩家
type Player struct {
	// 玩家ID
	PID int

	// 链接, 用于和客户端通信(不是服务器和客户端的链接!)
	Conn jinterface.IConnection

	// 当前玩家的X坐标
	X float32

	// 当前玩家的Y坐标(这里是高度（client）决定)
	Y float32

	// 当前玩家的Z坐标(平面y坐标)
	Z float32

	// 当前玩家的角度
	V float32
}

// PidGen Player Id 生成器
var PidGen int32 = 1  // 用于生成玩家ID的计数器（数据库中可直接得到）
var IdLock sync.Mutex // 用于生成玩家ID的锁

// NewPlayer 初始化玩家
func NewPlayer(conn jinterface.IConnection) *Player {
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()
	p := &Player{
		Conn: conn,
		PID:  int(id),
		X:    float32(160 + rand.Intn(10)), // 随机生成一个坐标
		Y:    0,
		Z:    float32(140 + rand.Intn(20)), // 随机生成一个坐标
		V:    0,
	}

	return p
}

/*
SendMsg
提供一个发送给客户端消息的方法
主要是将pb的protobuf数据序列化之后再发送，调用jinx的链接的SendMsg方法
*/
func (p *Player) SendMsg(msgID uint32, data proto.Message) {
	// 将proto数据序列化
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg error: ", err)
		return
	}

	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}

	// 将序列化后的数据发送给客户端
	if err := p.Conn.SendMsg(msgID, msg); err != nil {
		fmt.Println("player send msg error: ", err)
		return
	}
}

// SyncPid 同步玩家自己的ID给客户端
func (p *Player) SyncPid() {
	// 组建MsgID:1的proto数据
	data := &pb.SyncPid{
		Pid: int32(p.PID),
	}
	// 发送给客户端
	p.SendMsg(1, data)
}

// BroadCastStartPosition 广播玩家自己的出生地点给周围的玩家
func (p *Player) BroadCastStartPosition() {
	// 组建MsgID:200的proto数据
	protoMsg := &pb.Broadcast{
		Pid: int32(p.PID),
		Tp:  2,
		Data: &pb.Broadcast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	// 将消息发送给周围的玩家
	p.SendMsg(200, protoMsg)
}