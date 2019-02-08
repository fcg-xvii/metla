package metla

import (
	"errors"
	_ "fmt"
)

type valueType byte
type checkValMethod func([]byte) bool

const (
	valTypeInt valueType = iota
	valTypeFloat
	valTypeString
	valTypeBool
	valTypeNil
	valTypePointer
)

type valChecker struct {
	valType valueType
	method  checkValMethod
}

var (
	valCheckers = []*valChecker{
		&valChecker{valTypeInt, checkValInt},
		&valChecker{valTypeFloat, checkValFloat},
	}
)

func getStartTypes(first []byte) (res []*valChecker) {
	for _, checker := range valCheckers {
		if checker.method(first) {
			res = append(res, checker)
		}
	}
	return res
}

type val interface {
	Val() interface{}
	Type() valueType
	Data() ([]byte, error)
}

func initVal(source []byte) (res val, length int, err error) {
	if len(source) == 0 {
		return nil, errors.New("Value parse error :: source slice is empty")
	}
	types := getStartTypes(source[:1])
	count := len(source)
	if count > 64 {
		count = 64
	}
	for i := 1; i < count; i++ {
		offset := 0
		for i, v := range types {
			if offset > 0 {
				types[i-offset] = v
			}
			if !v.method(source[:i]) {
				offset++
			}
		}
		if offset > 0 {
			types = types[:offset]
		}
		if len(types) == 0 {
			return nil, errors.New("Unexpected value type...")
		}
	}
	switch types[0].valType {
	case valTypeInt:
		return newValInt(source)
	default:
		return
	}
}
