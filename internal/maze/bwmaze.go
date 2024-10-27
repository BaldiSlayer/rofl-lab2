package maze

import "strings"

type Maze struct {
	// use bitset for better memory utilisation
	grid [][]bool
}

func New(grid [][]bool) *Maze {
	return &Maze{
		grid: grid,
	}
}

func (m *Maze) Print() []string {
	ans := make([]string, len(m.grid))

	for y := 0; y < len(m.grid); y++ {
		var sb strings.Builder

		for x := 0; x < len(m.grid[0]); x++ {
			if m.grid[y][x] {
				sb.WriteByte(' ')
			} else {
				sb.WriteByte('#')
			}
		}

		ans[y] = sb.String()
	}

	return ans
}

// IsEmpty проверяет пуста ли клетка (то есть в ней нет стены)
func (m *Maze) IsEmpty(i, j int) bool {
	return m.grid[i][j]
}

// IsOut проверяет, находится ли клетка вне лабиринта
func (m *Maze) IsOut(i, j int) bool {
	return i < 0 || i >= len(m.grid) || j < 0 || j >= len(m.grid[0])
}
