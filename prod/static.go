package prod

import (
	"fmt"
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
		p.stack.Push(static{pos, p.MarkValString(0)[1:]}) // Обрезаем кавычки для результирующего значения
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
		p.stack.Push(static{pos, res})
	} else {
		res, _ := strconv.ParseFloat(p.MarkValString(0), 64)
		p.stack.Push(static{pos, res})
	}
	return nil
}

func newArray(p *parser) *parseError {
	var list []coordinator
	pos, stackLen := position{p.tplName, p.Line(), p.LinePos()}, p.stack.Len()

	flushArg := func() *parseError {
		if stackLen+1 != p.stack.Len() {
			return p.initParseError(p.Line(), p.LinePos(), "Expected one value")
		}
		list = append(list, p.stack.Pop().(coordinator))
		return nil
	}
	valArrived := false
	p.IncPos()
	for !p.IsEndDocument() {
		p.PassSpaces()
		switch p.Char() {
		case ']':
			if valArrived {
				flushArg()
			}
			p.stack.Push(static{pos, list})
			p.IncPos()
			return nil
		case ',', '\n', ';':
			if valArrived {
				if err := flushArg(); err != nil {
					return err
				}
				valArrived = false
			}
			p.IncPos()
		default:
			valArrived = true
			if err := p.initCodeVal(); err != nil {
				return err
			}
		}
	}
	return pos.parseError("Unclosed array token")
}

func newObject(p *parser) *parseError {
	fmt.Println("OBJECT")
	rMap := make(map[string]coordinator)
	key, stackLen, pos := "", p.stack.Len(), position{p.tplName, p.Line(), p.LinePos()}

	flushValue := func() *parseError {
		if key == "" && stackLen == p.stack.Len() {
			return nil
		} else if key == "" {
			return p.initParseError(p.Line(), p.LinePos(), "Unexpected value, name expected")
		} else if p.stack.Len() != stackLen+1 {
			return p.initParseError(p.Line(), p.LinePos(), "Expected single token")
		}
		rMap[key], key = p.stack.Pop().(coordinator), ""
		return nil
	}

	p.IncPos()
	for !p.IsEndDocument() {
		p.PassSpaces()
		fmt.Println("CHHH", string(p.Char()))
		switch ch := p.Char(); ch {
		case ':':
			fmt.Println("!!!!!!!!!!!!!!")
			if key != "" {
				return p.initParseError(p.Line(), p.LinePos(), "Unexpected ':' splitter, value expected")
			}
			if p.stack.Len() != stackLen+1 {
				return p.initParseError(p.Line(), p.LinePos(), "Expected string token 1")
			}
			if g, check := p.stack.Pop().(static); !check {
				return p.initParseError(p.Line(), p.LinePos(), "Expected string token 2")
			} else if key, check = g.get(nil).(string); !check {
				return p.initParseError(p.Line(), p.LinePos(), "Expected string token 3")
			}
			fmt.Println("KEY")
			p.IncPos()

		case '\n', ',', '}':
			if err := flushValue(); err != nil {
				return err
			}
			p.IncPos()
			if ch == '}' {
				p.stack.Push(static{pos, rMap})
				return nil
			}
		default:
			if err := p.initCodeVal(); err != nil {
				return err
			}
		}
	}
	return pos.parseError("Unclosed object token")
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
