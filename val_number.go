/*
 *  Объекты целого числа и числа с плавающей точкой,
 *  а так же их креаторы с методами поверок типа
 *
 */
package metla

import (
	"strconv"

	"github.com/fcg-xvii/lineman"
)

func newValNumber(p *parser) (err error) {
	p.SetupMark()
	intVal := true
	for lineman.CheckNumber(p.Char()) || p.Char() == '.' {
		if p.Char() == '.' {
			if !intVal {
				err = p.positionError("Unexpected float point")
			} else {
				intVal = false
			}
		}
		p.IncPos()
	}
	if intVal {
		res, _ := strconv.ParseInt(p.MarkValString(0), 10, 64)
		p.stack.Push(res)
	} else {
		res, _ := strconv.ParseFloat(p.MarkValString(0), 64)
		p.stack.Push(res)
	}
	return
}
