package defaults

import (
	"github.com/BaldiSlayer/rofl-lab2/pkg/models"
)

func GetAlphabet() []byte {
	return []byte{'N', 'S', 'W', 'E'}
}

func GetDirections() []models.Vector {
	return []models.Vector{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
}
