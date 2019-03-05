/*
 *  Объекты целого числа и числа с плавающей точкой,
 *  а так же их креаторы с методами поверок типа
 *
 */
package metla

import (
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/fcg-xvii/lineman"
)

// Вшиваем креаторы valInt и valFloat в глобальный срез
func init() {
	creators = append(creators, &valueCreator{
		checker:     checkValInt,
		constructor: newValInt,
	}, &valueCreator{
		checker:     checkValFloat,
		constructor: newValFloat,
	})
}

// Проверка соответствия valInt (целого числа)
// Будет совпадать, если хотя бы один символ к ряду является числом и не является точкой
func checkValInt(src []byte) bool {
	if !lineman.CheckNumber(src[0]) {
		return false
	}
	for _, v := range src {
		if !lineman.CheckNumber(v) {
			if v == '.' {
				return false
			} else {
				return true
			}
		}
	}
	return true
}

// Проверка соответствия valFloat (числа с плавающей точкой)
// Будет совпадать, если ряд начинается с числа и разделён (или заканчивается) одой точкой,
// 2-е точки к ряду не будут соответствовать типу (нет чисел с 2-я точками)
func checkValFloat(src []byte) (pointArrived bool) {
	if !lineman.CheckNumber(src[0]) {
		return false
	}
	for _, v := range src[1:] {
		if !lineman.CheckNumber(v) {
			if v == '.' {
				if pointArrived {
					return false
				} else {
					pointArrived = true
				}
			} else {
				break
			}
		}
	}
	return
}

//////////////////////////////////////////////////////////

// Конструктор целого числа. Тут всё просто - находим ряд цифр
func newValInt(p *parser) (token, error) {
	for !p.IsEndDocument() && lineman.CheckNumber(p.Char()) {
		p.IncPos()
	}
	res := new(valInt)
	res.val, _ = strconv.ParseInt(p.MarkValString(0), 10, 64)
	return res, nil
}

type valInt struct {
	val int64
}

func (s *valInt) Val() (interface{}, error) {
	return s.val, nil
}

func (s *valInt) Vals() ([]interface{}, error) {
	return []interface{}{s.val}, nil
}

func (s *valInt) Data(w io.Writer) (err error) {
	_, err = w.Write([]byte(strconv.FormatInt(s.val, 10)))
	return
}

func (s *valInt) String() string     { return "[int :: {" + strconv.FormatInt(s.val, 10) + "}]" }
func (s *valInt) IsExecutable() bool { return false }

func (s *valInt) execObject(sto *storage, tpl *template) (execObject, error) {
	return s, nil
}

func (s *valInt) IsNil() bool        { return false }
func (s *valInt) Type() reflect.Kind { return reflect.Int64 }
func (s *valInt) ValSingle() bool    { return true }

//////////////////////////////////////////////////////////

// Конструктор числа с плавающей точкой - тут так же всё просто, находим ряд чисел и точку
func newValFloat(p *parser) (token, error) {
	if !lineman.CheckNumber(p.Char()) {
		err := fmt.Errorf("Float parse error :: Unexpected float value [%c]", p.Char())
		return nil, err
	}
	p.IncPos()
	for !p.IsEndDocument() && (lineman.CheckNumber(p.Char()) || p.Char() == '.') {
		p.IncPos()
	}
	res := new(valFloat)
	res.val, _ = strconv.ParseFloat(p.MarkValString(0), 64)
	return res, nil
}

type valFloat struct {
	val float64
}

func (s *valFloat) Data(w io.Writer) (err error) {
	_, err = w.Write([]byte(strconv.FormatFloat(s.val, 'F', -1, 64)))
	return
}

func (s *valFloat) String() string {
	return "[float :: {" + strconv.FormatFloat(s.val, 'f', -1, 64) + "}]"
}

func (s *valFloat) IsExecutable() bool { return false }

func (s *valFloat) execObject(sto *storage, tpl *template) (execObject, error) {
	return s, nil
}

func (s *valFloat) Val() (interface{}, error) {
	return s.val, nil
}

func (s *valFloat) Vals() ([]interface{}, error) {
	return []interface{}{s.val}, nil
}

func (s *valFloat) IsNil() bool        { return false }
func (s *valFloat) Type() reflect.Kind { return reflect.Float64 }
func (s *valFloat) ValSingle() bool    { return true }
