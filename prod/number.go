package prod

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/fcg-xvii/lineman"
)

func newValNumber(p *parser) *parseError {
	p.SetupMark()
	line, pos := p.Line(), p.LinePos()
	intVal := true
	for lineman.CheckNumber(p.Char()) || p.Char() == '.' {
		if p.Char() == '.' {
			if !intVal {
				return p.initParseError(line, pos, errors.New("Unexpected float point"))
			} else {
				intVal = false
			}
		}
		p.IncPos()
	}
	if intVal {
		res, _ := strconv.ParseInt(p.MarkValString(0), 10, 64)
		fmt.Println("RES", res)
		p.stack.Push(res)
	} else {
		res, _ := strconv.ParseFloat(p.MarkValString(0), 64)
		fmt.Println("RES", res)
		p.stack.Push(res)
	}
	return nil
}
