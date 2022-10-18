# MMO_Jinx

## AOI算法和创建
![AOI](/img/AOI.png)
![AOI](/img/AOI_alg.png)
![AOI](/img/AOI管理.png)


## protobuf
![protobuf](/img/protobuf.png)

## 游戏业务 proto3
![proto3](/img/proto3.png)
### MsgID:1
SyncPid：

- 同步玩家本次登录的ID(用来标识玩家), 玩家登陆之后，由Server端主动生成玩家ID发送给客户端
- 发起者： Server
- Pid: 玩家ID

### MsgID:2
Talk:

- 同步玩家本次登录的ID(用来标识玩家), 玩家登陆之后，由Server端主动生成玩家ID发送给客户端
- 发起者： Client
- Content: 聊天信息
### MsgID:3
MovePackege:

- 移动的坐标数据
- 发起者： Client
- P: Position类型，地图的左边点

### MsgID:200
BroadCast:

- 广播消息
- 发起者： Server
- Tp: 1 世界聊天, 2 坐标, 3 动作, 4 移动之后坐标信息更新
- Pid: 玩家ID
```protobuf
message BroadCast{
	int32 Pid=1;
	int32 Tp=2;
	oneof Data {
        string Content=3;
        Position P=4;
		int32 ActionData=5;
    }
}
```

### MsgID:201
SyncPid：
- 广播消息 掉线/aoi消失在视野
- 发起者： Server
- Pid: 玩家ID
```protobuf
message SyncPid{
	int32 Pid=1;
}
```
### MsgID:202
- 同步周围的人位置信息(包括自己)
- 发起者： Server
- ps: Player 集合,需要同步的玩家
```protobuf
message SyncPlayers{
	repeated Player ps=1;
}

message Player{
	int32 Pid=1;
	Position P=2;
}
```

![proto3](/img/arc.png)