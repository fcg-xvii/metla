package metla

import (
	_ "fmt"
	"strconv"
)

func checkValInt(src []byte) bool {
	for _, v := range src {
		if !checkNumber(v) {
			return false
		}
	}
	return true
}

func checkValFloat(src []byte) bool {
	if !checkNumber(src[0]) {
		return false
	}
	for _, v := range src[1:] {
		if !checkNumber(v) && v != '.' {
			return false
		}
	}
	return true
}

//////////////////////////////////////////////////////////

func newValInt(source []byte) (res *valInt, length int, err error) {
	count := 0
	for i, v := range source {
		if !checkNumber(v) {
			count = i
			break
		}
	}
	res = new(valInt)
	res.val, _ = strconv.ParseInt(string(source[:count]), 10, 64)
	return
	/*var val int64
	if val, err = strconv.ParseInt(string(source), 10, 64); err == nil {
		res = &valInt{val}
	}
	return*/
}

type valInt struct {
	val int64
}

func (s *valInt) Val() interface{} {
	return s.val
}

func (s *valInt) Type() valueType {
	return valTypeInt
}

func (s *valInt) Data() (res []byte, err error) {
	return []byte(strconv.FormatInt(s.val, 10)), nil
}

//////////////////////////////////////////////////////////

func newValFloat(source []byte) (res *valInt, length int, err error) {
	return
}

type valFloat struct {
	val float64
}

func (s *valFloat) Val() interface{} {
	return s.val
}

func (s *valFloat) Type() valueType {
	return valTypeFloat
}

func (s *valFloat) Data() ([]byte, error) {
	return []byte(strconv.FormatFloat(s.val, 'F', -1, 64)), nil
}
