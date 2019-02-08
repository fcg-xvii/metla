package metla

import (
	"errors"
)

func checkValString(src []byte) bool {
	return src[0] == '\'' || src[0] == "'"
}

func NewValString(source []byte) (res *valString, length int, err error) {
	charID := source[0]
	for...
	
}

type valString struct {
	val string
}

func (s *valString) Val() interface{} {
	return s.val
}

func (s *valString) Type() valueType {
	return valTypeString
}

func (s *valString) Data() (res []byte, err error) {
	return []byte(s.val), nil
}
