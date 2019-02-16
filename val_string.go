package metla

import (
	"bytes"
	"fmt"
)

func checkValString(src []byte) bool {
	fmt.Println("SRC >>> ", src)
	return src[0] == '\'' || src[0] == '"'
}

func NewValString(source []byte) (res *valString, length int, err error) {
	charID := source[0]
	if index := bytes.IndexByte(source[1:], charID); index == -1 {
		err = fmt.Errorf("Unclosed string value [%c]", charID)
	} else {
		res = &valString{string(source[1 : index+1])}
		length = index + 2
	}
	return
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
