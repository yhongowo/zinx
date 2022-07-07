package core

import "fmt"

type AOIManager struct {
	MinX  int           //Area's left border pos
	MaxX  int           //Area's right border pos
	CntsX int           //X(horizontal) direction grid counts
	MinY  int           //Area's top border pos
	MaxY  int           //Area's bottom border pos
	CntsY int           //Y(vertical) direction grid counts
	grids map[int]*Grid //Current area's grids, K= Grid.GID, V= Grid
}

func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grids: make(map[int]*Grid),
	}
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			gID := y*cntsX + x
			aoiMgr.grids[gID] = NewGrid(gID,
				aoiMgr.MinX+x*aoiMgr.gridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.MinY+y*aoiMgr.gridLength(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridLength())
		}
	}
	return aoiMgr
}

func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManagr:\nminX:%d, maxX:%d, cntsX:%d, minY:%d, maxY:%d, cntsY:%d\n Grids in AOI Manager:\n",
		m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}

	return s
}

func (m *AOIManager) GetSurroundGridsByGID(gID int) (grids []*Grid) {
	//判断gID是否存在
	if _, ok := m.grids[gID]; !ok {
		return
	}

	//将当前gID添加到九宫格中
	grids = append(grids, m.grids[gID])

	// 根据gID, 得到格子所在的坐标
	x, y := gID%m.CntsX, gID/m.CntsX

	// 新建一个临时存储周围格子的数组
	surroundGID := make([]int, 0)

	// 新建8个方向向量: 左上: (-1, -1), 左中: (-1, 0), 左下: (-1,1), 中上: (0,-1), 中下: (0,1), 右上:(1, -1)
	// 右中: (1, 0), 右下: (1, 1), 分别将这8个方向的方向向量按顺序写入x, y的分量数组
	dx := []int{-1, -1, -1, 0, 0, 1, 1, 1}
	dy := []int{-1, 0, 1, -1, 1, -1, 0, 1}

	// 根据8个方向向量, 得到周围点的相对坐标, 挑选出没有越界的坐标, 将坐标转换为gID
	for i := 0; i < 8; i++ {
		newX := x + dx[i]
		newY := y + dy[i]

		if newX >= 0 && newX < m.CntsX && newY >= 0 && newY < m.CntsY {
			surroundGID = append(surroundGID, newY*m.CntsX+newX)
		}
	}

	// 根据没有越界的gID, 得到格子信息
	for _, gID := range surroundGID {
		grids = append(grids, m.grids[gID])
	}

	return
}

func (m *AOIManager) GetGIDByPos(x, y float32) int {
	gx := (int(x) - m.MinX) / m.gridWidth()
	gy := (int(y) - m.MinY) / m.gridLength()

	return gy*m.CntsX + gx
}

//通过横纵坐标得到周边九宫格内的全部PlayerIDs
func (m AOIManager) GetPIDsByPos(x, y float32) (playerIDs []int) {

	gID := m.GetGIDByPos(x, y)

	girds := m.GetSurroundGridsByGID(gID)
	for _, v := range girds {
		playerIDs = append(playerIDs, v.GetPlayerIDs()...)
	}
	return
}

func (m *AOIManager) RemovePIDFromGrid(pID, gID int) {
	m.grids[gID].Remove(pID)
}

//添加一个PlayerID到一个格子中
func (m *AOIManager) AddPIDToGrid(pID, gID int) {
	m.grids[gID].Add(pID)
	
}

//通过横纵坐标添加一个Player到一个格子中
func (m *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGIDByPos(x, y)
	grid := m.grids[gID]
	grid.Add(pID)
}

//通过横纵坐标把一个Player从对应的格子中删除
func (m *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := m.GetGIDByPos(x, y)
	grid := m.grids[gID]
	grid.Remove(pID)
}
