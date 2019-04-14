package prod

import (
	"strconv"

	"github.com/fcg-xvii/lineman"
)

func newValNumber(p *parser) *parseError {
	p.SetupMark()
	pos := position{p.tplName, p.Line(), p.LinePos()}
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
		&position{p.Line(), p.LinePos() + offset},
	}
}

type static struct {
	*position
	val interface{}
}

func (s *static) Get() (interface{}, *execError) {
	return s.val, nil
}

func (s *static) String() string {
	return "{ static }"
}
