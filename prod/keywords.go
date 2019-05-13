package prod

import (
	"reflect"
)

type keywordConstructor func(*parser) *parseError

func init() {
	keywords["len"] = func(p *parser) *parseError {
		return initCoreFunction(coreLen, p)
	}
	//keywords["echo"] = keywordEcho
	//keywords["echoln"] = keywordEcholn
	//keywords["print"] = keywordPrint
	//keywords["println"] = keywordPrintln
}

var (
	keywords = map[string]keywordConstructor{
		"nil": func(p *parser) *parseError {
			p.stack.Push(initStatic(p, nil, 3))
			return nil
		}, "true": func(p *parser) *parseError {
			p.stack.Push(initStatic(p, true, 4))
			return nil
		}, "false": func(p *parser) *parseError {
			p.stack.Push(initStatic(p, false, 5))
			return nil
		}, "var": func(p *parser) *parseError {
			if p.varFlag {
				return p.initParseError(p.Line(), p.LinePos()-3, "Unexpected var keyword")
			}
			p.varFlag = true
			return nil
		}, "endfor": func(p *parser) *parseError {
			if p.cycleLayout == 0 {
				return p.initParseError(p.Line(), p.Pos(), "Unexpected endfor token")
			}
			for i := len(p.execList) - 1; i >= 0; i-- {
				switch obj := p.execList[i]; obj.(type) {
				case cycler:
					cycle := obj.(cycler)
					if !cycle.isClosed() {
						commands := make([]executer, len(p.execList)-i-1)
						copy(commands, p.execList[i+1:])
						cycle.setCommands(commands)
						p.execList = p.execList[:i+1]
						cycle.closeCycle()
						p.cycleLayout--
						return nil
					}
				}
			}
			return p.initParseError(p.Line(), p.LinePos(), "Unexpected endfor token")
		},
	}
	functions = map[string]interface{}{
		//"len": coreLen,
		//"defined": coreDefined,
	}
)

func getKeywordConstructor(name string) (result keywordConstructor, check bool) {
	result, check = keywords[name]
	return
}

func initCoreFunction(method func(*tplExec, position, interface{}) *execError, p *parser) *parseError {
	r := coreFunc{position: position{p.tplName, p.Line(), p.LinePos()}, method: method}
	p.PassSpaces()
	if p.Char() != '(' {
		return p.initParseError(p.Line(), p.LinePos(), "Expected '(' character")
	}
	p.IncPos()
	if err := p.initCodeVal(); err != nil {
		return err
	}
	r.arg = p.stack.Pop()
	p.PassSpaces()
	if p.Char() != ')' {
		return p.initParseError(p.Line(), p.LinePos(), "Expected ')' character")
	}
	p.IncPos()
	p.stack.Push(r)
	return nil
}

type coreFunc struct {
	position
	arg    interface{}
	method func(*tplExec, position, interface{}) *execError
}

func (s coreFunc) exec(exec *tplExec) *execError {
	var val interface{}
	switch s.arg.(type) {
	case getter:
		val = s.arg.(getter).get(exec)
	case executer:
		stackLen := exec.stack.Len()
		if err := s.arg.(executer).exec(exec); err != nil {
			return err
		}
		if stackLen+1 != exec.stack.Len() {
			return s.execError("Not one return value in argument")
		}
		val = exec.stack.Pop().(getter).get(exec)
	}
	return s.method(exec, s.position, val)
}

func (s coreFunc) execType() execType {
	return execFunction
}

func coreLen(exec *tplExec, pos position, arg interface{}) *execError {
	rVal := reflect.ValueOf(arg)
	switch rVal.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array:
		exec.stack.Push(static{pos, rVal.Len()})
		return nil
	default:
		return pos.execError("Expected map, slice or array argument")
	}
}
