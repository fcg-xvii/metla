package metla

import (
	"io"
)

type execObjectType byte

const (
	execObjectUndefined = iota
	execObjectString
	execObjectInteger
	execObjectFloat
	execObjectObject
	execObjectArray
	execObjectByteSlice
	execObjectNil
	execObjectToken
)

var (
	execObjectTypeString = []string{
		"undefined",
		"string",
		"integer",
		"float",
		"object",
		"array",
		"byteSlice",
		"nil",
		"token",
	}
)

func (s execObjectType) String() string {
	if s < 0 || int(s) >= len(execObjectTypeString) {
		return execObjectTypeString[0]
	} else {
		return execObjectTypeString[s]
	}
}

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
	Type() execObjectType
	Val() interface{}
	Vals() []interface{}
	ValSingle() bool
	IsNil() bool
	String() string
}
