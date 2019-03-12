package metla

import (
	"fmt"
	"io"
	"reflect"
)

// Общий интерфейс объекта результирующих данных
type token interface {
	execObject(sto *storage, tpl *template, parent execObject) (execObject, error)
	IsExecutable() bool
	String() string
	fatalError(string) error
}

// Интерфейс токена, который может содержать блок дочерних токенов
type tokenParent interface {
	appendChild(token)
}

// Интерфейс контейнера значения
type value interface {
	token
	StaticVal() interface{}
	Bool() bool
	Float() float64
	Int() int64
}

type execObject interface {
	Data(io.Writer) error // Запись результирующих данных в выходной поток
	Type() reflect.Kind
	Val() (interface{}, error)
	Vals() ([]interface{}, error)
	ValSingle() bool
	IsNil() bool
	String() string
	positionWarning(string) error
	receiveEvent(name string, params []interface{}) bool
}

func checkKindInt(t reflect.Kind) bool {
	return t == reflect.Int64 || t == reflect.Int32 || t == reflect.Int16 || t == reflect.Int8 || t == reflect.Int
}

func checkKindFloat(t reflect.Kind) bool {
	return t == reflect.Float32 || t == reflect.Float64
}

func checkIfaceInt(i interface{}) (res int64, check bool) {
	val := reflect.ValueOf(i)
	if check = checkKindInt(val.Kind()); check {
		res = val.Int()
	}
	return
}

func checkIfaceFloat(i interface{}) (res float64, check bool) {
	val := reflect.ValueOf(i)
	if check := checkKindFloat(val.Kind()); check {
		res = val.Float()
	}
	return
}

func checkIfaceNumber(i interface{}) (check, integer bool) {
	val := reflect.ValueOf(i)
	if check = checkKindInt(val.Kind()); check {
		integer = true
	} else {
		check = checkKindFloat(val.Kind())
	}
	return
}

type rawInfoRecord struct {
	tplName   string
	line, pos int
}

func (s *rawInfoRecord) fatalError(text string) error {
	return fmt.Errorf("Fatal error [%v %v:%v]: %v", s.tplName, s.line, s.pos, text)
}

func (s *rawInfoRecord) positionWarning(text string) error {
	return fmt.Errorf("Warning [%v %v:%v]: %v", s.tplName, s.line, s.pos, text)
}

type eventExec struct {
	parent execObject
}

func (s *eventExec) sendEvent(name string, params []interface{}) {
	if s.parent != nil {
		//s.parent.re
	}
}

func (s *eventExec) receiveEvent(name string, params []interface{}) bool {
	if s.parent != nil {
		return s.parent.receiveEvent(name, params)
	}
	return false
}
