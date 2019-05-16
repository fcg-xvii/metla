package metla

import (
	"fmt"
	"reflect"
	"time"
)

func init() {
	keywords["for"] = newFor
	keywords["range"] = newRange
	keywords["each"] = newEach
}

type cycler interface {
	isClosed() bool
	closeCycle()
	getCommands() []executer
	setCommands([]executer)
	cLayout() byte
	tLayout() byte
}

type cycle struct {
	position
	commands     []executer
	closed       bool
	cycleLayout  byte
	threadLayout byte
}

func (s *cycle) isClosed() bool                  { return s.closed }
func (s *cycle) closeCycle()                     { s.closed = true }
func (s *cycle) setCommands(commands []executer) { s.commands = commands }
func (s *cycle) getCommands() []executer         { return s.commands }
func (s *cycle) cLayout() byte                   { return s.cycleLayout }
func (s *cycle) tLayout() byte                   { return s.threadLayout }

func newForCheck(cycle *cCycle, p *parser) *parseError {
	p.PassSpaces()
	if p.Char() != ';' && p.Char() != '\n' {
		return cycle.parseError("Expected endline")
	}
	p.stack.Push(cycle)
	return nil
}

func newFor(p *parser) (err *parseError) {
	c := &cCycle{cycle: &cycle{position: position{p.tplName, p.Line(), p.LinePos()}, cycleLayout: p.cycleLayout, threadLayout: p.threadLayout}}
	p.cycleLayout++
	p.store.incLayout()

	for !p.IsEndLine() {
		if err := p.initCodeVal(); err != nil {
			return err
		}
		p.PassSpaces()
	}
	switch token := p.stack.Pop(); token.(type) {
	case rpn:
		c.checkPN = token.(rpn)
		return newForCheck(c, p)
	default:
		return token.(coordinator).parseError(fmt.Sprintf("Unexpected cycle token type '%T'", token))
	}
}

type cCycle struct {
	*cycle
	checkPN rpn
}

func (s *cCycle) String() string {
	return "{ for }"
}

func (s *cycle) checkMaxTime(exec *tplExec) bool {
	return time.Now().After(exec.execStop)
}

func (s *cCycle) exec(exec *tplExec) (err *execError) {
	var checkRes bool
	for {
		if s.checkMaxTime(exec) {
			return s.execError("Fatal error :: Excess maximum execute time")
		}
		if err = s.checkPN.execToBool(exec, &checkRes); err != nil {
			return err
		}
		if checkRes {
			for _, v := range s.commands {
				if err = v.exec(exec); err != nil {
					return err
				} else if exec.breakFlag {
					exec.breakFlag = false
					return nil
				}
			}
		} else {
			break
		}
	}
	return
}

func (s *cCycle) execType() execType {
	return execFor
}

////////////////////////////////////

func newRange(p *parser) *parseError {
	fmt.Println("NEW_RANGE", p.cycleLayout, p.threadLayout)
	r := &cRange{cycle: &cycle{position: position{p.tplName, p.Line(), p.LinePos()}, cycleLayout: p.cycleLayout, threadLayout: p.threadLayout}}
	p.cycleLayout++
	p.store.incLayout()
	p.varFlag = true
	if err := p.initCodeVal(); err != nil {
		return err
	}
	p.varFlag = false
	cVar, check := p.stack.Pop().(iName)
	if !check {
		return r.parseError("Variable name expected")
	}
	r.countVar = cVar
	p.PassSpaces()
	if !p.PosMatchSlice([]byte("in")) {
		return r.parseError("Expected 'in' keyword after count variable name")
	}
	p.ForwardPos(2)
	if err := p.initCodeVal(); err != nil {
		return err
	}
	p.PassSpaces()
	if p.Char() != ':' {
		return p.initParseError(p.Line(), p.LinePos(), "Expected ':' token")
	}
	r.min = p.stack.Pop()
	p.IncPos()
	if err := p.initCodeVal(); err != nil {
		return err
	}
	r.max = p.stack.Pop()
	p.PassSpaces()
	if !p.IsEndLine() {
		return p.initParseError(p.Line(), p.LinePos(), "Ecpected endline")
	}
	p.stack.Push(r)
	return nil
}

type cRange struct {
	*cycle
	min, max interface{}
	countVar iName
}

func (s *cRange) String() string {
	return "{ range }"
}

func (s *cRange) setMinMax(obj interface{}, result *int, exec *tplExec) *execError {
	var rVal reflect.Value
	switch obj.(type) {
	case getter:
		rVal = reflect.ValueOf(obj.(getter).get(exec))
	case executer:
		stackLen := exec.stack.Len()
		if err := obj.(executer).exec(exec); err != nil {
			return err
		}
		if stackLen+1 != exec.stack.Len() {
			return obj.(coordinator).execError("Expected one return value")
		}
		rVal = reflect.ValueOf(exec.stack.Pop().(getter).get(exec))
	}

	resType := reflect.ValueOf(*result).Type()
	if rType := rVal.Type(); !rType.ConvertibleTo(resType) {
		return obj.(coordinator).execError("Expected integer friendly type")
	} else {
		*result = int(rVal.Convert(resType).Int())
		return nil
	}
}

func (s *cRange) exec(exec *tplExec) *execError {
	var min, max int
	if err := s.setMinMax(s.min, &min, exec); err != nil {
		return err
	}
	if err := s.setMinMax(s.max, &max, exec); err != nil {
		return err
	}
	if min > max {
		return s.dec(min, max, exec)
	} else {
		return s.inc(min, max, exec)
	}
}

func (s *cRange) dec(min, max int, exec *tplExec) *execError {
	for i := min; i > max; i-- {
		if s.checkMaxTime(exec) {
			return s.execError("Fatal error :: Excess maximum execute time")
		}
		s.countVar.set(exec, i)
		for _, v := range s.commands {
			if err := v.exec(exec); err != nil {
				return err
			}
			if exec.breakFlag {
				exec.breakFlag = false
				return nil
			}
		}
	}
	return nil
}

func (s *cRange) inc(min, max int, exec *tplExec) *execError {
	for i := min; i < max; i++ {
		if s.checkMaxTime(exec) {
			return s.execError("Fatal error :: Excess maximum execute time")
		}
		s.countVar.setRaw(exec, i)
		for _, v := range s.commands {
			if err := v.exec(exec); err != nil {
				return err
			}
			if exec.breakFlag {
				exec.breakFlag = false
				return nil
			}
		}
	}
	return nil
}

func (s *cRange) execType() execType {
	return execFor
}

func newEach(p *parser) *parseError {
	fmt.Println("NEW_EACH", p.cycleLayout, p.threadLayout, len(p.execList))
	r := &cEach{cycle: &cycle{position: position{p.tplName, p.Line(), p.LinePos()}, cycleLayout: p.cycleLayout, threadLayout: p.threadLayout}}
	p.cycleLayout++
	p.store.incLayout()
	if err := r.parseVar(p, &r.keyVar); err != nil {
		return err
	}
	p.PassSpaces()
	if p.Char() != ',' {
		return p.initParseError(p.Line(), p.LinePos(), "Expected ',' splitter")
	}
	p.IncPos()
	if err := r.parseVar(p, &r.valVar); err != nil {
		return err
	}
	p.PassSpaces()
	if !p.PosMatchSlice([]byte("->")) {
		return r.parseError("Each :: expected '->' keyword after count variable name")
	}
	p.ForwardPos(2)
	stackLen := p.stack.Len()
	for !p.IsEndDocument() {
		p.PassSpaces()
		if p.IsEndLine() {
			if stackLen+1 != p.stack.Len() {
				return r.parseError("Each :: single token expected")
			}
			r.objVar = p.stack.Pop()
			fmt.Println("PUSH", r)
			p.stack.Push(r)
			return nil
		} else {
			if err := p.initCodeVal(); err != nil {
				return err
			}
		}
	}
	return r.parseError("Each :: unexpected end of document")
}

type cEach struct {
	*cycle
	keyVar *iName
	valVar *iName
	objVar interface{}
}

func (s *cEach) String() string {
	return "{ each }"
}

func (s *cEach) parseVar(p *parser, kVar **iName) *parseError {
	p.PassSpaces()
	if p.Char() == '_' {
		p.IncPos()
		return nil
	} else {
		p.varFlag = true
		if err := p.initCodeVal(); err != nil {
			return err
		}
		p.varFlag = false
		crd := p.stack.Pop().(coordinator)
		if vVar, check := crd.(iName); !check {
			return crd.parseError("Expected variable token")
		} else {
			*kVar = &vVar
			return nil
		}
	}
}

func (s *cEach) exec(exec *tplExec) *execError {
	var rVal reflect.Value
	switch s.objVar.(type) {
	case getter:
		rVal = reflect.ValueOf(s.objVar.(getter).get(exec))
	case executer:
		if err := s.objVar.(executer).exec(exec); err != nil {
			return err
		}
		rVal = reflect.ValueOf(exec.stack.Pop().(getter).get(exec))
	}
	switch rVal.Kind() {
	case reflect.Slice, reflect.Array:
		{
			for i := 0; i < rVal.Len(); i++ {
				if s.checkMaxTime(exec) {
					return s.execError("Fatal error :: Excess maximum execute time")
				}
				if s.keyVar != nil {
					s.keyVar.setRaw(exec, i)
				}
				if s.valVar != nil {
					s.valVar.setRaw(exec, rVal.Index(i).Interface())
				}
				for _, v := range s.commands {
					if err := v.exec(exec); err != nil {
						return err
					}
				}
			}
			return nil
		}
	case reflect.Map:
		{
			iterator := rVal.MapRange()
			for iterator.Next() {
				if s.checkMaxTime(exec) {
					return s.execError("Fatal error :: Excess maximum execute time")
				}
				if s.keyVar != nil {
					s.keyVar.setRaw(exec, iterator.Key().Interface())
				}
				if s.valVar != nil {
					s.valVar.setRaw(exec, iterator.Value().Interface())
				}
				for _, v := range s.commands {
					if err := v.exec(exec); err != nil {
						return err
					}
					if exec.breakFlag {
						exec.breakFlag = false
						return nil
					}
				}
			}
			return nil
		}
	default:
		return s.objVar.(coordinator).execError("Unexpected each object type, map, slice or array needed")
	}
}

func (s *cEach) execType() execType {
	return execFor
}
