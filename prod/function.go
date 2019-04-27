package prod

import (
	"fmt"
	"reflect"
)

func parseFunctionArgs(p *parser, fPos *position) (list []interface{}, err *parseError) {
	valAccepted := true
	for !p.IsEndDocument() {
		switch {
		case p.isEndCode():
			return nil, p.initParseError(p.Line(), p.LinePos(), "Uncosed function")
		case p.Char() == ',' || p.Char() == ')':
			if valAccepted && len(list) > 0 {
				return nil, p.initParseError(p.Line(), p.LinePos(), "Unexpected comma - value expected")
			}
			valAccepted = true
			if p.Char() == ')' {
				return
			} else {
				p.IncPos()
			}
		default:
			if err = p.initCodeVal(); err != nil {
				return
			}
			list = append(list, p.stack.Pop())
			valAccepted = false
		}
	}
	return nil, (*fPos).parseError("Unexpected end of document")
}

func newFunction(p *parser) (err *parseError) {
	f, check := function{position: p.posObject()}, false
	if f.nameVar, check = p.stack.Pop().(*iName); !check {
		return f.parseError("Function parse error :: expected variable in prefix")
	}
	//return p.initParseError(10, 10, "Function error test")
	p.IncPos()
	if f.args, err = parseFunctionArgs(p, &f.position); err != nil {
		return err
	}
	fmt.Println(">>>>", f.args)
	p.stack.Push(&f)
	p.IncPos()
	return
}

type function struct {
	position
	nameVar getter
	args    []interface{}
}

func (s *function) exec(exec *tplExec) *execError {
	//return s.execError("TTTTT")
	stackLen := exec.stack.Len()
	rName := reflect.ValueOf(s.nameVar.get(exec))
	if rName.Kind() != reflect.Func {
		return s.execError("Variable is not a function")
	}
	rArgs, rType := make([]reflect.Value, 0, len(s.args)), rName.Type()

	for _, v := range rArgs {
		switch v.(type) {
		case getter:
			index := len(rArgs)

		}
	}
}
