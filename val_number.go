/*
 *  Объекты целого числа и числа с плавающей точкой,
 *  а так же их креаторы с методами поверок типа
 *
 */
package metla

import (
	"fmt"
	"strconv"
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
	if !checkNumber(src[0]) {
		return false
	}
	for _, v := range src {
		if !checkNumber(v) {
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
	fmt.Println("CHECK_VAL_FLOAT")
	if !checkNumber(src[0]) {
		return false
	}
	for _, v := range src[1:] {
		if !checkNumber(v) {
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

// Конструктор целого числа. Тут всё просто - находим ряд чисел
func newValInt(p *parser) (token, error) {
	for !p.isEndDocument() && checkNumber(p.char()) {
		p.incPos()
	}
	res := new(valInt)
	res.val, _ = strconv.ParseInt(p.markValString(0), 10, 64)
	return res, nil
}

type valInt struct {
	val int64
}

func (s *valInt) Val() interface{} {
	return s.val
}

func (s *valInt) Data() (res []byte, err error) {
	return []byte(strconv.FormatInt(s.val, 10)), nil
}

func (s *valInt) String() string     { return "[int :: {" + strconv.FormatInt(s.val, 10) + "}]" }
func (s *valInt) IsExecutable() bool { return false }

//////////////////////////////////////////////////////////

// Конструктор числа с плавающей точкой - тут так же всё просто, находим ряд чисел и точку
func newValFloat(p *parser) (token, error) {
	if !checkNumber(p.char()) {
		err := fmt.Errorf("Unexpected float value [%c]", p.char())
		return nil, err
	}
	p.incPos()
	for !p.isEndDocument() && (checkNumber(p.char()) || p.char() == '.') {
		p.incPos()
	}
	res := new(valFloat)
	res.val, _ = strconv.ParseFloat(p.markValString(0), 64)
	return res, nil
}

type valFloat struct {
	val float64
}

func (s *valFloat) Val() interface{} {
	return s.val
}

func (s *valFloat) Data() ([]byte, error) {
	return []byte(strconv.FormatFloat(s.val, 'F', -1, 64)), nil
}

func (s *valFloat) String() string {
	return "[float :: {" + strconv.FormatFloat(s.val, 'f', -1, 64) + "}]"
}
func (s *valFloat) IsExecutable() bool { return false }
