package core

import (
	"fmt"
	"sync"
)

/*
	AOI地图中的格子类型
*/

type Grid struct {
	//格子ID
	gID int

	//当前格子左边边界坐标
	minX int

	//当前格子右边边界坐标
	maxX int

	//当前格子上边边界坐标
	minY int

	//当前格子下边边界坐标
	maxY int

	//当前格子内玩家或物体的ID集合
	playerIDs map[int]bool

	//锁
	pIDLock sync.RWMutex
}

// NewGrid 初始化格子
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		gID:       gID,
		minX:      minX,
		maxX:      maxX,
		minY:      minY,
		maxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

// Add 添加玩家
func (g *Grid) Add(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playerID] = true
}

// Remove 移除玩家
func (g *Grid) Remove(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, playerID)
}

// GetPlayerIDs 获取当前格子内玩家ID集合
func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for id, _ := range g.playerIDs {
		playerIDs = append(playerIDs, id)
	}
	return
}

// override String() 打印格子信息
func (g *Grid) String() string {
	return fmt.Sprintf("Grid id: %d, minX: %d, maxX: %d, minY: %d, maxY: %d, playerIDs: %v",
		g.gID, g.minX, g.maxX, g.minY, g.maxY, g.playerIDs)
}
