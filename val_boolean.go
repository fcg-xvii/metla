package metla

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

func init() {
	creators = append(creators, &valueCreator{
		checker:     checkValBoolean,
		constructor: newValBoolean,
	})
}

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

func (s *valBoolean) execObject(*storage, *template, execObject) (execObject, error) {
	return s, nil
}

func (s *valBoolean) IsExecutable() bool                                  { return false }
func (s *valBoolean) Type() reflect.Kind                                  { return reflect.Bool }
func (s *valBoolean) StaticVal() interface{}                              { return s.val }
func (s *valBoolean) Bool() bool                                          { return s.val }
func (s *valBoolean) Float() float64                                      { return float64(s.Int()) }
func (s *valBoolean) IsNil() bool                                         { return false }
func (s *valBoolean) Val() (interface{}, error)                           { return s.val, nil }
func (s *valBoolean) Vals() ([]interface{}, error)                        { return []interface{}{s.val}, nil }
func (s *valBoolean) Data(w io.Writer) (err error)                        { _, err = w.Write([]byte(s.String())); return }
func (s *valBoolean) ValSingle() bool                                     { return true }
func (s *valBoolean) receiveEvent(name string, params []interface{}) bool { return false }

func (s *valBoolean) Int() (res int64) {
	if s.val {
		res = 1
	}
	return
}

func (s *valBoolean) String() (res string) {
	res = "false"
	if s.val {
		res = "true"
	}
	return
}

/////////////////////////////////////////////////////////////////////////////////////////////

type valNil struct {
	*rawInfoRecord
}

func (s *valNil) execObject(*storage, *template, execObject) (execObject, error) {
	return s, nil
}

func (s *valNil) IsExecutable() bool                                  { return false }
func (s *valNil) Type() reflect.Kind                                  { return reflect.Bool }
func (s *valNil) StaticVal() interface{}                              { return nil }
func (s *valNil) IsNil() bool                                         { return true }
func (s *valNil) Val() (interface{}, error)                           { return nil, nil }
func (s *valNil) Vals() ([]interface{}, error)                        { return nil, nil }
func (s *valNil) ValSingle() bool                                     { return true }
func (s *valNil) Data(w io.Writer) (err error)                        { _, err = w.Write([]byte("nil")); return }
func (s *valNil) String() string                                      { return "nil" }
func (s *valNil) receiveEvent(name string, params []interface{}) bool { return false }