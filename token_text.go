package metla

import (
	"io"
	"reflect"
	"strconv"
)

type tokenText struct {
	*rawInfoRecord
	src []byte
}

func (s *tokenText) execObject(*storage, *template) (execObject, error) {
	return s, nil
}

func (s *tokenText) Data(w io.Writer) (err error) {
	_, err = w.Write(s.src)
	return
}

func (s *tokenText) String() string {
	return "[text :: { []byte :: len:" + strconv.Itoa(len(s.src)) + " }]"
}

func (s *tokenText) IsExecutable() bool           { return false }
func (s *tokenText) IsNil() bool                  { return false }
func (s *tokenText) Type() reflect.Kind           { return reflect.Slice }
func (s *tokenText) Val() (interface{}, error)    { return s.src, nil }
func (s *tokenText) Vals() ([]interface{}, error) { return []interface{}{s.src}, nil }
func (s *tokenText) ValSingle() bool              { return true }
