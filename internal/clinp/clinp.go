package clinp

import (
	"github.com/BaldiSlayer/rofl-lab2/internal/mat"
)

// InputProcessor занимается обработкой команд поступающих из stdin
// сделан, чтобы любой мог легко подменить реализацию на свою
type InputProcessor interface {
	ProcessCommands(teacher mat.MAT, commandsChan chan string)
}
