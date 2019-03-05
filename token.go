package metla

import (
	"io"
	"reflect"
)

// Общий интерфейс объекта результирующих данных
type token interface {
	execObject(*storage, *template) (execObject, error)
	IsExecutable() bool
	String() string
}

// Интерфейс контейнера значения
type value interface {
	token
	Val() interface{}
}

type execObject interface {
	Data(io.Writer) error // Запись результирующих данных в выходной поток
	Type() reflect.Kind
	Val() (interface{}, error)
	Vals() ([]interface{}, error)
	ValSingle() bool
	IsNil() bool
	String() string
}
