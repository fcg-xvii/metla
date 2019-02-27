package metla

import (
	"io"
	"strconv"
)

type tokenText struct {
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

func (s *tokenText) IsExecutable() bool   { return false }
func (s *tokenText) IsNil() bool          { return false }
func (s *tokenText) Type() execObjectType { return execObjectToken }
func (s *tokenText) Val() interface{}     { return s.src }
func (s *tokenText) Vals() []interface{}  { return []interface{}{s.src} }
func (s *tokenText) ValSingle() bool      { return true }
