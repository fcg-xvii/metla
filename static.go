package metla

import (
	"fmt"
	"strconv"

	"github.com/fcg-xvii/lineman"
)

func newValString(p *parser) *parseError {
	p.SetupMark()
	pos := position{p.tplName, p.Line(), p.LinePos()}
	charID := p.Char() // Определяем, двойная или одинарная кавычка открыта
	p.IncPos()         // Смещаемся на начало строки
	res, check := "", false
	for {
		if !p.ToChar(charID) { // Пытаемся найти закрывающую кавычку
			return p.initParseError(pos.line, pos.pos, "Unclosed string")
		} else {
			// Закрывающая кавычка найдена. Инициализируем результирующее значение
			if p.PrevChar() != '\\' {
				if check {
					res += p.MarkValString(0)[0:]
				} else {
					res += p.MarkValString(0)[1:]
				}
				p.stack.Push(static{pos, res}) // Обрезаем кавычки для результирующего значения
				p.IncPos()
				return nil
			} else {
				if check {
					res += p.MarkValString(0)[0 : p.Pos()-p.MarkPos()-1]
				} else {
					res += p.MarkValString(0)[1 : p.Pos()-p.MarkPos()-1]
				}
				check = true
				p.SetupMark()
				p.IncPos()
			}
		}
	}
	return nil
}

func newValNumber(p *parser) *parseError {
	p.SetupMark()
	pos := position{p.tplName, p.Line(), p.LinePos() - 1}
	intVal := true
	for lineman.CheckNumber(p.Char()) || p.Char() == '.' {
		if p.Char() == '.' {
			if !intVal {
				return p.initParseError(pos.line, pos.pos, "Unexpected float point")
			} else {
				intVal = false
			}
		}
		p.IncPos()
	}
	if intVal {
		res, _ := strconv.ParseInt(p.MarkValString(0), 10, 64)
		p.stack.Push(static{pos, res})
	} else {
		res, _ := strconv.ParseFloat(p.MarkValString(0), 64)
		p.stack.Push(static{pos, res})
	}
	return nil
}

func initStatic(p *parser, val interface{}, offset int) static {
	return static{
		position{p.tplName, p.Line(), p.LinePos() + offset},
		val,
	}
}

type static struct {
	position
	val interface{}
}

func (s static) get(*tplExec) interface{} {
	return s.val
}

func (s static) String() string {
	return fmt.Sprintf("{ static: %v }", s.val)
}
