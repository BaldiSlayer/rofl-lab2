package defaults

import (
	"github.com/BaldiSlayer/rofl-lab2/pkg/models"
)

func GetStartState() models.Cell {
	return models.Cell{X: 0, Y: 0}
}

func GetAlphabet() []byte {
	return []byte{'N', 'S', 'W', 'E'}
}

func GetDirections() []models.Vector {
	return []models.Vector{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
}

const (
	EpsilonSymbol = "e"
)
