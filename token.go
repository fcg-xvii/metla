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

// Интерфейс токена, который может содержать блок дочерних токенов
type tokenParent interface {
	appendChild(token)
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

func checkKindInt(t reflect.Kind) bool {
	return t == reflect.Int64 || t == reflect.Int32 || t == reflect.Int16 || t == reflect.Int8 || t == reflect.Int
}

func checkIfaceInt(i interface{}) (res int64, check bool) {
	val := reflect.ValueOf(i)
	if check = checkKindInt(val); check {
		res = val.Int()
	}
	return
}
