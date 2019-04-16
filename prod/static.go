package prod

import (
	"strconv"

	"github.com/fcg-xvii/lineman"
)

func newValString(p *parser) *parseError {
	p.SetupMark()
	pos := position{p.tplName, p.Line(), p.LinePos()}
	charID := p.Char()     // Определяем, двойная или одинарная кавычка открыта
	p.IncPos()             // Смещаемся на начало строки
	if !p.ToChar(charID) { // Пытаемся найти закрывающую кавычку
		return p.initParseError(pos.line, pos.pos, "Unclosed string")
	} else {
		// Закрывающая кавычка найдена. Инициализируем результирующее значение
		p.stack.Push(&static{&pos, p.MarkValString(0)[1:]}) // Обрезаем кавычки для результирующего значения
		p.IncPos()
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
		p.stack.Push(&static{&pos, res})
	} else {
		res, _ := strconv.ParseFloat(p.MarkValString(0), 64)
		p.stack.Push(&static{&pos, res})
	}
	return nil
}

func initStatic(p *parser, val interface{}, offset int) *static {
	return &static{
		&position{p.tplName, p.Line(), p.LinePos() + offset},
		val,
	}
}

type static struct {
	*position
	val interface{}
}

func (s *static) Get(*tplExec) interface{} {
	return s.val
}

func (s *static) String() string {
	return "{ static }"
}
