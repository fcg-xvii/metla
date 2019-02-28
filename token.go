package metla

import (
	"io"
	"reflect"
)

/*type execObjectType byte

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
	execObjectFunction
	execObjectMethod
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
		"function",
		"method",
	}
)

func (s execObjectType) String() string {
	if s < 0 || int(s) >= len(execObjectTypeString) {
		return execObjectTypeString[0]
	} else {
		return execObjectTypeString[s]
	}
}*/

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
