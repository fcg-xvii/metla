package metla

import (
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
	for i, v := range source {
		if !checkNumber(v) || i == len(source)-1 {
			length = i + 1
			break
		}
	}
	res = new(valInt)
	res.val, _ = strconv.ParseInt(string(source[:length]), 10, 64)
	return
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

func (s *valInt) String() string { return "[int :: {" + strconv.FormatInt(s.val, 10) + "}]" }

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
