package metla

import (
	"fmt"
	"io"
	"reflect"
)

type positionInformer interface {
	fatalError(string) error
	posInfo() *rawInfoRecord
}

// Общий интерфейс объекта результирующих данных
type token interface {
	positionInformer
	execObject(sto *storage, tpl *template, parent executor) (executor, error)
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
	Type() reflect.Kind
	IsStatic() bool
	StaticVal() interface{}
	IsNumber() bool
	IsNil() bool
}

type valueNumber interface {
	value
	IsInteger() bool
	Float() float64
	Int() int64
	Add(float64)
}

type valueBoolean interface {
	value
	Bool() bool
}

type executor interface {
	String() string
	Data(io.Writer) error // Запись результирующих данных в выходной поток
	//Data(io.Writer) error // Запись результирующих данных в выходной поток
	//Type() reflect.Kind
	//Val() (interface{}, error)
	//Vals() ([]interface{}, error)
	//ValSingle() bool
	//IsNil() bool
	//String() string
	//positionWarning(string) error
	//receiveEvent(name string, params []interface{}) bool
}

/*type valueExecutor interface {
	executor
	Type() reflect.Kind

}*/

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

func canInt(val float64) bool {
	return val == float64(int64(val))
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

func (s *rawInfoRecord) posInfo() *rawInfoRecord {
	return s
}

type eventExec struct {
	parent executor
}

func (s *eventExec) sendEvent(name string, params []interface{}) {
	if s.parent != nil {
		//s.parent.re
	}
}

/*func (s *eventExec) receiveEvent(name string, params []interface{}) bool {
	if s.parent != nil {
		return s.parent.receiveEvent(name, params)
	}
	return false
}*/
