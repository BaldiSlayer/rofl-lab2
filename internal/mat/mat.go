package mat

import "github.com/BaldiSlayer/rofl-lab2/internal/models"

type MAT interface {
	// Generate - генерация лабиринта
	Generate() error
	// Include - запрос на включение
	Include(query string) (bool, error)
	// Equal - запрос на эквивалентность
	Equal() (models.EqualResponse, error)
	// Print - визуализация лабиринта
	Print() ([]string, error)
}
