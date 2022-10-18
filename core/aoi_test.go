package core

import (
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	aoiMgr := NewAOIManager(0, 250, 0, 250, 5, 5)
	for _, grid := range aoiMgr.grids {
		t.Logf("grid id: %d, minX: %d, maxX: %d, minY: %d, maxY: %d", grid.gID, grid.minX, grid.maxX, grid.minY, grid.maxY)
	}
}

func TestAOIManager_GetSurroundGridsByGID(t *testing.T) {
	aoiMgr := NewAOIManager(0, 250, 0, 250, 5, 5)
	for gid, _ := range aoiMgr.grids {
		grids := aoiMgr.GetSurroundGridsByGID(gid)
		t.Logf("gid: %d, len: %d", gid, len(grids))
		gids := make([]int, 0, len(grids))
		for _, grid := range grids {
			gids = append(gids, grid.gID)
		}
		t.Logf("grids: %v", gids)
	}
}
