package mat

import (
	"github.com/BaldiSlayer/rofl-lab2/pkg/models"
)

type MAT interface {
	// Generate - генерация лабиринта
	Generate() error
	// Include - запрос на включение
	Include(query string) (bool, error)
	// Equal - запрос на эквивалентность
	Equal(prefixes []string, suffixes []string, matrix [][]bool) (models.EqualResponse, error)
	// Print - визуализация лабиринта
	Print() ([]string, error)
}
