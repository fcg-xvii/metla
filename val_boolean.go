package metla

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

/*func init() {
	creators = append(creators, &valueCreator{
		checker:     checkValBoolean,
		constructor: newValBoolean,
	}, &valueCreator{
		checker:     checkValNil,
		constructor: newValNil,
	})
}*/

func checkValBoolean(src []byte) bool {
	return bytes.Index(src, []byte("true")) == 0 || bytes.Index(src, []byte("false")) == 0
}

func newValBoolean(p *parser) (res token, err error) {
	if name, check := p.ReadName(); !check {
		err = p.positionError("Boolean value expected")
	} else if string(name) == "true" {
		res = &valBoolean{p.infoRecordFromMark(), true}
	} else if string(name) == "false" {
		res = &valBoolean{p.infoRecordFromMark(), false}
	} else {
		err = p.positionError(fmt.Sprintf("Unexpected boolean name '%v', 'true' or 'false' expected", string(name)))
	}
	return
}

type valBoolean struct {
	*rawInfoRecord
	val bool
}

func (s *valBoolean) execObject(*storage) (executor, error) {
	return s, nil
}

func (s *valBoolean) posInfo() *rawInfoRecord                             { return s.rawInfoRecord }
func (s *valBoolean) IsExecutable() bool                                  { return false }
func (s *valBoolean) Type() reflect.Kind                                  { return reflect.Bool }
func (s *valBoolean) IsStatic() bool                                      { return true }
func (s *valBoolean) IsNumber() bool                                      { return false }
func (s *valBoolean) StaticVal() interface{}                              { return s.val }
func (s *valBoolean) Bool() bool                                          { return s.val }
func (s *valBoolean) IsNil() bool                                         { return false }
func (s *valBoolean) Val() (interface{}, error)                           { return s.val, nil }
func (s *valBoolean) Vals() ([]interface{}, error)                        { return []interface{}{s.val}, nil }
func (s *valBoolean) Data(w io.Writer) (err error)                        { _, err = w.Write([]byte(s.String())); return }
func (s *valBoolean) ValSingle() bool                                     { return true }
func (s *valBoolean) receiveEvent(name string, params []interface{}) bool { return false }

func (s *valBoolean) String() (res string) {
	res = "false"
	if s.val {
		res = "true"
	}
	res = "[bool { " + res + " }]"
	return
}

/////////////////////////////////////////////////////////////////////////////////////////////

func checkValNil(src []byte) bool {
	return bytes.Index(src, []byte("nil")) == 0
}

func newValNil(p *parser) (res token, err error) {
	if name, check := p.ReadName(); !check {
		err = p.positionError("Nil value expected")
	} else if string(name) == "nil" {
		res = &valNil{p.infoRecordFromMark()}
	} else {
		err = p.positionError(fmt.Sprintf("Unexpected name '%v', 'nil' expected", string(name)))
	}
	return
}

type valNil struct {
	*rawInfoRecord
}

func (s *valNil) execObject(*storage) (executor, error) {
	return s, nil
}

func (s *valNil) IsExecutable() bool                                  { return false }
func (s *valNil) Type() reflect.Kind                                  { return reflect.Invalid }
func (s *valNil) StaticVal() interface{}                              { return nil }
func (s *valNil) IsNil() bool                                         { return true }
func (s *valNil) Bool() bool                                          { return false }
func (s *valNil) Float() float64                                      { return 0 }
func (s *valNil) Int() int64                                          { return 0 }
func (s *valNil) IsNumber() bool                                      { return false }
func (s *valNil) IsStatic() bool                                      { return true }
func (s *valNil) Val() (interface{}, error)                           { return nil, nil }
func (s *valNil) Vals() ([]interface{}, error)                        { return nil, nil }
func (s *valNil) ValSingle() bool                                     { return true }
func (s *valNil) Data(w io.Writer) (err error)                        { _, err = w.Write([]byte("nil")); return }
func (s *valNil) String() string                                      { return "nil" }
func (s *valNil) receiveEvent(name string, params []interface{}) bool { return false }

func (s *valNil) v() value {
	return s
}
