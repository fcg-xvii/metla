package metla

import (
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/fcg-xvii/lineman"
)

func init() {
	creators = append(creators, &valueCreator{
		checker:     checkValVariable,
		constructor: newValVariable,
	})
}

func checkValVariable(src []byte) bool {
	if lineman.CheckFirsNameChar(src) > 0 {
		for i := 1; i < len(src); i++ {
			if lineman.CheckBodyNameChar(src[i:]) == 0 {
				res := src[i] != '(' && src[i] != '['
				return res
			}
		}
		return true
	}
	return false
}

// Конструктор строки.
func newValVariable(p *parser) (res token, err error) {
	if name, check := p.ReadName(); !check {
		err = errors.New("Variable parse error :: Unexpected name")
	} else {
		res = &valVariable{string(name)}
	}
	return
}

type valVariable struct {
	name string
}

func (s *valVariable) Val() interface{} {
	return s.name
}

func (s *valVariable) String() string     { return "[variable :: { " + s.name + " }]" }
func (s *valVariable) IsExecutable() bool { return false }

func (s *valVariable) execObject(sto *storage, tpl *template) (execObject, error) {
	if v, check := sto.findVariable(s.name); !check {
		return nil, fmt.Errorf("Variable not found [%v]", s.name)
	} else {
		return &valVariableExec{v}, nil
	}
}

///////////////////////////////////////////////////////////////////

type valVariableExec struct {
	v *variable
}

func (s *valVariableExec) Data(w io.Writer) (err error) {
	_, err = w.Write([]byte(fmt.Sprint(s.v.value)))
	return
}

func (s *valVariableExec) String() string {
	return "[variable { name: " + s.v.key + ", value: " + fmt.Sprint(s.v.value) + " }]"
}

func (s *valVariableExec) IsNil() bool        { return s.v.IsNil() }
func (s *valVariableExec) Type() reflect.Kind { return s.v.Kind() }
func (s *valVariableExec) ValSingle() bool    { return true }

func (s *valVariableExec) Val() (interface{}, error) {
	return s.v.value, nil
}

func (s *valVariableExec) Vals() ([]interface{}, error) {
	return []interface{}{s.v.value}, nil
}
