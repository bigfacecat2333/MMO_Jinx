package core

import (
	"fmt"
	"math"
)

const (
	AoiMinX int = 85
	AoiMaxX int = 410
	AoiCntX int = 10
	AoiMinY int = 75
	AoiMaxY int = 400
	AoiCntY int = 20
)

/*
	AOI地图(管理)模块
*/

type AOIManager struct {
	// 区域的左边界坐标
	MinX int

	// 区域的右边界坐标
	MaxX int

	// 区域的上边界坐标
	MinY int

	// 区域的下边界坐标
	MaxY int

	// 每个格子在X轴方向的格子数量
	CntX int

	// 每个格子在Y轴方向的格子数量
	CntY int

	// 当前AOI地图中有哪些格子map[格子ID] *格子对象
	grids map[int]*Grid
}

// NewAOIManager 初始化一个AOI地图管理模块
func NewAOIManager(minX, maxX, minY, maxY, cntX, cntY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		MinY:  minY,
		MaxY:  maxY,
		CntX:  cntX,
		CntY:  cntY,
		grids: make(map[int]*Grid),
	}

	// 初始化AOI地图中的格子, 进行编号
	for y := 0; y < cntY; y++ {
		for x := 0; x < cntX; x++ {
			// 计算格子ID
			gID := y*cntX + x

			// 初始化一个格子
			aoiMgr.grids[gID] = NewGrid(gID,
				aoiMgr.MinX+x*aoiMgr.gridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.MinY+y*aoiMgr.gridHeight(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridHeight(),
			)
		}
	}

	return aoiMgr
}

// gridWidth 计算每个格子在X轴方向的宽度
func (am *AOIManager) gridWidth() int {
	return (am.MaxX - am.MinX) / am.CntX
}

// gridHeight 计算每个格子在Y轴方向的高度
func (am *AOIManager) gridHeight() int {
	return (am.MaxY - am.MinY) / am.CntY
}

// 调试打印当前AOI地图中的格子信息
func (am *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager:\n MinX:%d, MaxX:%d, MinY:%d, MaxY:%d, CntX:%d, CntY:%d\n Grids in AOIManager:\n",
		am.MinX, am.MaxX, am.MinY, am.MaxY, am.CntX, am.CntY)

	for _, grid := range am.grids {
		s += fmt.Sprintln(grid)
	}

	return s
}

// GetSurroundGridsByGID 根据格子GID获取周边九宫格内的全部格子信息
func (am *AOIManager) GetSurroundGridsByGID(gID int) (grids []*Grid) {
	// 如果当前格子不存在, 则直接返回
	if _, ok := am.grids[gID]; !ok {
		return
	}

	// 将当前格子加入到九宫格切片中
	grids = append(grids, am.grids[gID])

	// 判断当前格子是否在左边,右边是否有格子
	// 得到当前格子的x轴索引 %
	idx := gID % am.CntX

	// 如果当前格子不在最左边, 则将左边的格子加入到九宫格切片中
	if idx > 0 {
		grids = append(grids, am.grids[gID-1])
	}
	// 如果当前格子不在最右边, 则将右边的格子加入到九宫格切片中
	if idx < am.CntX-1 {
		grids = append(grids, am.grids[gID+1])
	}

	// 将当前的格子都取出来, 然后判断上面和下面是否有格子
	gridsX := make([]int, 0, len(grids))
	for _, grid := range grids {
		gridsX = append(gridsX, grid.gID)
	}

	for _, grid := range gridsX {
		// 得到当前格子的y轴索引 /
		idy := grid / am.CntX
		// 如果当前格子不在最上面, 则将上面的格子加入到九宫格切片中
		if idy > 0 {
			grids = append(grids, am.grids[grid-am.CntX])
		}
		// 如果当前格子不在最下面, 则将下面的格子加入到九宫格切片中
		if idy < am.CntY-1 {
			grids = append(grids, am.grids[grid+am.CntX])
		}
	}

	return
}

// GetGIDByPos 通过坐标得到当前坐标所在的格子ID
func (am *AOIManager) GetGIDByPos(x, y float32) int {
	// 根据格子的宽高计算当前坐标所在的格子的索引
	idx := int(math.Floor(float64((x - float32(am.MinX)) / float32(am.gridWidth()))))
	idy := int(math.Floor(float64((y - float32(am.MinY)) / float32(am.gridHeight()))))

	// 根据格子的索引计算出格子的ID
	gID := idy*am.CntX + idx

	return gID
}

// GetPlayerIDsByPos GetSurroundGridsByPos 根据玩家坐标获取周边九宫格内的全部格子信息
func (am *AOIManager) GetPlayerIDsByPos(x, y float32) (playerIDs []int) {
	// 得到当前玩家所在的格子ID
	gID := am.GetGIDByPos(x, y)

	// 通过gid得到周边九宫格信息
	grids := am.GetSurroundGridsByGID(gID)

	// 将九宫格内的全部玩家信息取出来，放入playerIDs中
	for _, grid := range grids {
		playerIDs = append(playerIDs, grid.GetPlayerIDs()...)
	}

	return
}

// AddPlayerToGrid 添加一个玩家到一个格子中
func (am *AOIManager) AddPlayerToGrid(playerID, gID int) {
	// 将玩家添加到对应的格子中
	am.grids[gID].Add(playerID)
}

// RemovePlayerFromGrid 从一个格子中删除一个玩家
func (am *AOIManager) RemovePlayerFromGrid(playerID, gID int) {
	// 将玩家从对应的格子中删除
	am.grids[gID].Remove(playerID)
}

// AddPlayerToGridByPos 根据玩家坐标添加一个玩家到一个格子中
func (am *AOIManager) AddPlayerToGridByPos(playerID int, x, y float32) {
	// 得到当前玩家所在的格子ID
	gID := am.GetGIDByPos(x, y)

	// 将玩家添加到对应的格子中
	am.AddPlayerToGrid(playerID, gID)
}

// RemovePlayerFromGridByPos 根据玩家坐标从一个格子中删除一个玩家
func (am *AOIManager) RemovePlayerFromGridByPos(playerID int, x, y float32) {
	// 得到当前玩家所在的格子ID
	gID := am.GetGIDByPos(x, y)

	// 将玩家从对应的格子中删除
	am.RemovePlayerFromGrid(playerID, gID)
}

// GetPIDsByGID 根据格子ID得到全部玩家ID
func (am *AOIManager) GetPIDsByGID(gID int) (playerIDs []int) {
	// 将九宫格内的全部玩家信息取出来，放入playerIDs中
	playerIDs = append(playerIDs, am.grids[gID].GetPlayerIDs()...)

	return
}
