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
	Pid int32

	// 链接, 用于和客户端通信(conn是服务器接收后返回的链接，用于和客户端通信)
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
		Pid:  id,
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
		Pid: int32(p.Pid),
	}
	// 发送给客户端
	p.SendMsg(1, data)
}

// BroadCastStartPosition 广播玩家自己的出生地点给周围的玩家
func (p *Player) BroadCastStartPosition() {
	// 组建MsgID:200的proto数据
	protoMsg := &pb.Broadcast{
		Pid: int32(p.Pid),
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

// Talk 向周围的玩家广播自己的聊天内容
func (p *Player) Talk(content string) {
	//1. 组建MsgId200 proto数据
	msg := &pb.Broadcast{
		Pid: p.Pid,
		Tp:  1, //TP 1 代表聊天广播
		Data: &pb.Broadcast_Content{
			Content: content,
		},
	}

	//2. 得到当前世界所有的在线玩家
	players := WorldMgrObj.GetAllPlayers()

	//3. 向所有的玩家发送MsgId:200消息
	for _, player := range players {
		player.SendMsg(200, msg)
	}
}

// SyncSurrounding 同步周围玩家的位置信息给当前玩家
func (p *Player) SyncSurrounding() {
	// 1. 获取当前玩家周围的玩家
	players := p.GetSurroundingPlayers()

	// 2. 将玩家的位置信息发送给周围的玩家（广播），让周围的玩家看到当前玩家
	// 2.1 组建MsgID:200的Player proto数据
	protoMsg := &pb.Broadcast{
		Pid: p.Pid,
		Tp:  2, //TP 2 代表位置广播
		Data: &pb.Broadcast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	// 2.2 分别将消息发送给周围的玩家
	for _, player := range players {
		player.SendMsg(200, protoMsg)
	}

	// 3. 将周围玩家的位置信息发送给当前玩家，让当前玩家看到周围的玩家
	// 3.1 组建MsgID:202的Player proto数据
	// 3.1.1 创建一个Player切片，用于存放周围玩家的位置信息
	playersProtoMsg := make([]*pb.Player, 0, len(players))
	// 3.1.2 遍历周围玩家，将位置信息组建成Player proto数据
	for _, player := range players {
		// 制作Player proto数据
		p := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		// 将Player proto数据添加到切片中
		playersProtoMsg = append(playersProtoMsg, p)
	}

	// 3.1.2 封装SyncPlayers proto数据
	syncPlayersProtoMsg := &pb.SyncPlayers{
		Ps: playersProtoMsg[:],
	}

	// 3.2 将组件好的数据发给客户端
	p.SendMsg(202, syncPlayersProtoMsg)
}

// UpdatePos 更新玩家的位置信息
func (p *Player) UpdatePos(x, y, z, v float32) {
	// 更新玩家的位置信息
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	// 组建MsgID:200的proto数据 位置广播
	protoMsg := &pb.Broadcast{
		Pid: p.Pid,
		Tp:  4, //TP 4 代表位置更新广播
		Data: &pb.Broadcast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	// 获取当前玩家周围的玩家
	players := p.GetSurroundingPlayers()
	for _, player := range players {
		player.SendMsg(200, protoMsg)
	}
}

// GetSurroundingPlayers 获取周围的玩家
func (p *Player) GetSurroundingPlayers() []*Player {
	// 1. 获取当前玩家周围的玩家
	pids := WorldMgrObj.AoiMgr.GetPlayerIDsByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}

	return players
}
