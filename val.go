package metla

import (
	"errors"
	"fmt"
	"reflect"
)

func convert(val, defaultVal interface{}) (res interface{}, check bool) {
	res = val
	lVal, rVal := reflect.ValueOf(val), reflect.ValueOf(defaultVal)
	check = lVal.Kind() == rVal.Kind()
	if rVal.Kind() != reflect.Invalid && lVal.Kind() != rVal.Kind() {
		if lVal.Type().ConvertibleTo(rVal.Type()) {
			res, check = lVal.Convert(rVal.Type()).Interface(), true
		} else {
			res = defaultVal
		}
	}
	return
}

type valueCheckMethod func([]byte) bool
type valueConstructor func(p *parser) (res token, err error)

type valueCreator struct {
	checker     valueCheckMethod
	constructor valueConstructor
}

var (
	creators               []*valueCreator
	errValueUnexpectedType = errors.New("Value init error :: Unexpected value type")
)

func initCodeVal(p *parser) (val interface{}, err error) {
	p.PassSpaces()
	switch p.Char() {
	case '+', '-', '*', '/', '(', '!', '>', '<', '%', '&', '|':
		val, err = newValArifmetic(p)
	case '"', '\'':
		val, err = newValString(p)
	case ',':
		val, err = newValSet(p)
	case '=':
		if p.NextChar() != '=' {
			val, err = newValSet(p)
		} else {
			val, err = newValArifmetic(p)
		}
	case '{':
		val, err = newValObject(p)
	case '[':
		val, err = newValArray(p)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		val, err = newValNumber(p)
	default:
		if name, check := p.ReadName(); !check {
			err = p.positionError(fmt.Sprintf("Unexpected symbol '%c'", p.Char()))
		} else {
			if keyword, check := getKeywordConstructor(string(name)); check {
				val, err = keyword(p)
			} else {
				switch p.Char() {
				case '(':
					if fStatic, check := functions[string(name)]; check {
						val, err = newStaticFunction(fStatic, p)
					} else {
						val, err = newValFunction(string(name), p)
					}
				case '[':
					val = &valVariable{p.infoRecordFromMark(), string(name)}
					p.stack.Push(val)
					val, err = newValIndex(p)
				case '.':
					val = &valVariable{p.infoRecordFromMark(), string(name)}
					p.stack.Push(val)
					val, err = newValField(p)
				default:
					val = &valVariable{p.infoRecordFromMark(), string(name)}
					p.stack.Push(val)
				}
			}
		}
	}
	return
}
