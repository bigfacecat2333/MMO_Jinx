package core

import "sync"

// WorldManager 世界管理模块

type WorldManager struct {
	// AOI地图
	AoiMgr *AOIManager

	// 当前在线的玩家
	Players map[int32]*Player

	// 保护Player的读写锁
	playersLock sync.RWMutex
}

// WorldMgrObj 全局唯一的世界管理模块对象
var WorldMgrObj *WorldManager

// init 初始化世界管理模块,全局唯一
func init() {
	WorldMgrObj = &WorldManager{
		// 初始化AOI地图
		AoiMgr: NewAOIManager(AoiMinX, AoiMaxX, AoiMinY, AoiMaxY, AoiCntX, AoiCntY),
		// 初始化在线玩家集合
		Players: make(map[int32]*Player),
	}
}

// AddPlayer 添加一个玩家
func (wm *WorldManager) AddPlayer(player *Player) {
	// 加写锁
	wm.playersLock.Lock()
	defer wm.playersLock.Unlock()

	// 添加到世界管理器中
	wm.Players[int32(player.Pid)] = player

	// 将玩家添加到AOI地图中
	wm.AoiMgr.AddPlayerToGridByPos(int(player.Pid), player.X, player.Z)
}

// RemovePlayerByPid 根据玩家ID删除一个玩家
func (wm *WorldManager) RemovePlayerByPid(pid int32) {
	// 加写锁
	wm.playersLock.Lock()
	defer wm.playersLock.Unlock()

	// 从世界管理器中删除
	delete(wm.Players, pid)

	// 从AOI地图中删除
	wm.AoiMgr.RemovePlayerFromGridByPos(int(pid), wm.Players[pid].X, wm.Players[pid].Z)
}

// GetPlayerByPid 根据玩家ID获取一个玩家
func (wm *WorldManager) GetPlayerByPid(pid int32) *Player {
	// 加读锁
	wm.playersLock.RLock()
	defer wm.playersLock.RUnlock()

	return wm.Players[pid]
}

// GetPlayersByGid 根据格子GID获取全部玩家
func (wm *WorldManager) GetPlayersByGid(gid int) (players []*Player) {
	// 加读锁
	wm.playersLock.RLock()
	defer wm.playersLock.RUnlock()

	// 获取全部玩家ID
	pids := wm.AoiMgr.GetPIDsByGID(gid)

	// 根据玩家ID获取全部玩家
	for _, pid := range pids {
		players = append(players, wm.Players[int32(pid)])
	}

	return
}

// GetAllPlayers 获取全部在线玩家
func (wm *WorldManager) GetAllPlayers() []*Player {
	// 加读锁
	wm.playersLock.RLock()
	defer wm.playersLock.RUnlock()

	players := make([]*Player, 0)

	for _, player := range wm.Players {
		players = append(players, player)
	}

	return players
}
