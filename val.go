package metla

import (
	"errors"
	"fmt"
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
		&valChecker{valTypeString, checkValString},
	}
)

func getStartTypes(first []byte) (res []*valChecker) {
	fmt.Println(string(first))
	for _, checker := range valCheckers {
		if checker.method(first) {
			res = append(res, checker)
		}
	}
	return res
}

func defineType(p *parser, types []*valChecker) []*valChecker {
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
		fmt.Println("TL>>>", types)
		if len(types) == 0 {
			return types
		}
	}
	return types
}

func initVal(p *parser) (res value, length int, err error) {
	if len(source) == 0 {
		return nil, 0, errors.New("Value parse error :: source slice is empty")
	}
	types := getStartTypes(source[:1])
	//fmt.Println("TPS", types)
	l := len(types)
	switch {
	case l == 0:
		err = errors.New("Unexpected value type...")
		return
	case l > 1:
		types = defineType(p, types)
	}
	//fmt.Println("SOURCE...", string(source), types)
	switch types[0].valType {
	case valTypeInt:
		return newValInt(source)
	case valTypeFloat:
		return newValFloat(source)
	case valTypeString:
		return NewValString(source)
	default:
		return
	}
}
