package metla

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type keywordConstructor func(*parser) *parseError

func init() {
	keywords["len"] = func(p *parser) *parseError {
		return initCoreFunction(coreLen, p)
	}
	keywords["string"] = func(p *parser) *parseError {
		return initCoreFunction(coreString, p)
	}
	keywords["cmp"] = func(p *parser) *parseError {
		return initCoreFunction(coreCMP, p)
	}
	keywords["int"] = func(p *parser) *parseError {
		return initCoreFunction(coreInt, p)
	}
	keywords["join"] = func(p *parser) *parseError {
		return initCoreFunction(coreJoin, p)
	}
	keywords["unixToDate"] = func(p *parser) *parseError {
		return initCoreFunction(coreTimestampToDate, p)
	}
	//keywords["echo"] = keywordEcho
	//keywords["echoln"] = keywordEcholn
	//keywords["print"] = keywordPrint
	//keywords["println"] = keywordPrintln
}

var (
	keywords = map[string]keywordConstructor{
		"return": func(p *parser) *parseError {
			pos := position{p.tplName, p.Line(), p.LinePos()}
			p.stack.Push(cReturn{pos})
			return nil
		}, "exit": func(p *parser) *parseError {
			pos := position{p.tplName, p.Line(), p.LinePos()}
			p.stack.Push(cExit{pos})
			return nil
		}, "nil": func(p *parser) *parseError {
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
						p.store.decLayout()
						/*if cycle.cLayout() != p.cycleLayout {
							for _, v := range cycle.getCommands() {
								if c, check := v.(cycler); check && !c.isClosed() {
									return c.(coordinator).parseError("Unclosed for tag")
								}
							}
						}*/
						if cycle.tLayout() != p.threadLayout {
							for _, v := range commands {
								if thread, check := v.(*thread); check && !thread.closed {
									return thread.parseError("Unclosed if tag")
								}
							}
						}
						return nil
					}
				}
			}
			return p.initParseError(p.Line(), p.LinePos(), "Unexpected endfor token")
		}, "endif": func(p *parser) *parseError {
			if ck, i, check := findThread(p); !check {
				return p.initParseError(p.Line(), p.LinePos(), "Unexpected endif token - 'if' token not found")
			} else {
				p.threadLayout--
				ck.closed = true
				lastBlock := ck.blocks[len(ck.blocks)-1]
				lastBlock.commands = make([]executer, len(p.execList)-i-1)
				copy(lastBlock.commands, p.execList[i+1:])
				p.execList = p.execList[:i+1]
				p.store.decLayout()
				if ck.cycleLayout != p.cycleLayout {
					for _, block := range ck.blocks {
						for _, v := range block.commands {
							if c, check := v.(cycler); check && !c.isClosed() {
								return c.(coordinator).parseError("Unclosed for tag")
							}
						}
					}
				}
				return nil
			}
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

func initCoreFunction(method func(*tplExec, position, ...interface{}) *execError, p *parser) *parseError {
	//r, stackLen := coreFunc{position: position{p.tplName, p.Line(), p.LinePos()}, method: method}, p.stack.Len()
	r := coreFunc{position: position{p.tplName, p.Line(), p.LinePos()}, method: method}
	p.PassSpaces()
	if p.Char() != '(' {
		return p.initParseError(p.Line(), p.LinePos(), "Expected '(' character")
	}
	p.IncPos()
	/*if err := p.initCodeVal(); err != nil {
		return err
	}
	r.arg = p.stack.Pop()
	p.PassSpaces()
	if p.Char() != ')' {
		return p.initParseError(p.Line(), p.LinePos(), "Expected ')' character")
	}*/
	for !p.IsEndDocument() {
		p.PassSpaces()
		switch p.Char() {
		case ')':
			r.args = append(r.args, p.stack.Pop())
			p.stack.Push(r)
			p.IncPos()
			return nil
		case ',':
			r.args = append(r.args, p.stack.Pop())
			p.IncPos()
		default:
			if err := p.initCodeVal(); err != nil {
				return err
			}
		}
	}
	return r.parseError("len function - unexpected end document")
}

type coreFunc struct {
	position
	args   []interface{}
	method func(*tplExec, position, ...interface{}) *execError
}

func (s coreFunc) exec(exec *tplExec) *execError {
	var val []interface{}
	for i := 0; i < len(s.args); i++ {
		switch s.args[i].(type) {
		case getter:
			val = append(val, s.args[i].(getter).get(exec))
		case executer:
			stackLen := exec.stack.Len()
			if err := s.args[i].(executer).exec(exec); err != nil {
				return err
			}
			for exec.stack.Len() > stackLen {
				val = append(val, exec.stack.Pop().(getter).get(exec))
			}
		}
	}
	return s.method(exec, s.position, val...)
}

func (s coreFunc) execType() execType {
	return execFunction
}

func coreLen(exec *tplExec, pos position, arg ...interface{}) *execError {
	if len(arg) != 1 {
		return pos.execError("coreLen - expected 1 argument")
	}
	rVal := reflect.ValueOf(arg[0])
	switch rVal.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array, reflect.String:
		exec.stack.Push(static{pos, rVal.Len()})
	default:
		exec.stack.Push(static{pos, 0})
	}
	return nil
}

func coreString(exec *tplExec, pos position, arg ...interface{}) *execError {
	if len(arg) != 1 {
		return pos.execError("coreLen - expected 1 argument")
	}
	exec.stack.Push(static{pos, fmt.Sprint(arg[0])})
	return nil
}

func coreCMP(exec *tplExec, pos position, arg ...interface{}) *execError {
	var str string
	for _, v := range arg {
		str = fmt.Sprintf("%v%v", str, v)
	}
	exec.stack.Push(static{pos, str})
	return nil
}

func coreInt(exec *tplExec, pos position, arg ...interface{}) *execError {
	if len(arg) != 1 {
		return pos.execError("coreLen - expected 1 argument")
	}
	res, err := strconv.ParseInt(fmt.Sprint(arg[0]), 10, 64)
	if err != nil {
		return pos.execError(err.Error())
	}
	exec.stack.Push(static{pos, res})
	return nil
}

func coreJoin(exec *tplExec, pos position, arg ...interface{}) *execError {
	if len(arg) != 2 {
		return pos.execError("coreLen - expected 2 arguments - array and string splitter")
	}
	splitter := fmt.Sprint(arg[1])
	switch arg[0].(type) {
	case []interface{}:
		sl := arg[0].([]interface{})
		tmp := make([]string, len(sl))
		for i, v := range sl {
			tmp[i] = fmt.Sprint(v)
		}
		exec.stack.Push(static{pos, strings.Join(tmp, splitter)})
	case []string:
		exec.stack.Push(static{pos, strings.Join(arg[0].([]string), splitter)})
	default:
		return pos.execError("First argument slice or array expected")
	}
	return nil
}

func coreTimestampToDate(exec *tplExec, pos position, arg ...interface{}) *execError {
	if len(arg) != 1 {
		return pos.execError("coreLen - expected 1 argument")
	}
	var tmps int64
	switch arg[0].(type) {
	case int64:
		tmps = arg[0].(int64)
	case float64:
		tmps = int64(arg[0].(float64))
	}
	layout := "2006-01-02"
	t := time.Unix(tmps, 0)
	exec.stack.Push(static{pos, t.Format(layout)})
	return nil
}
