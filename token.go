package metla

import (
	"io"
)

// Общий интерфейс объекта результирующих данных
type token interface {
	Data(io.Writer, *storage) error // Запись результирующих данных в выходной поток
	IsExecutable() bool
	String() string
}

// Интерфейс контейнера значения
type value interface {
	token
	Val() interface{}
}
