package mat

import (
	"github.com/BaldiSlayer/rofl-lab2/internal/eqtable"
	"github.com/BaldiSlayer/rofl-lab2/pkg/models"
)

type MAT interface {
	// Generate - генерация лабиринта
	Generate() error
	// Include - запрос на включение
	Include(query string) (bool, error)
	// Equal - запрос на эквивалентность
	Equal(eqTable eqtable.EqTable) (models.EqualResponse, error)
	// Visualize - визуализация лабиринта
	Visualize() ([]string, error)
}
