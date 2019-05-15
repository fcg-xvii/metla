package metla

import "fmt"

func newValName(p *parser, line, pos int, key string) *parseError {
	name := iName{position{p.tplName, line, pos}, key, 0}
	if p.varFlag {
		index, err := p.store.setVariable(string(key))
		if err != nil {
			return p.initParseError(line, pos, err.Error())
		}
		name.index = index
	} else {
		name.index = p.store.initVariable(string(key))
	}
	p.stack.Push(name)
	return nil
}

type iName struct {
	position
	name  string
	index int
}

func (s iName) StorageIndex() int {
	return s.index
}

func (s iName) set(exec *tplExec, val interface{}) *execError {
	//fmt.Println("VAL.....", val)
	if g, check := val.(getter); check {
		exec.sto.setValue(s.index, g.get(exec))
	} else {
		return val.(coordinator).execError("Set variable error - expected getter right side")
	}
	return nil
}

func (s iName) setRaw(exec *tplExec, val interface{}) {
	exec.sto.setValue(s.index, val)
}

func (s iName) get(exec *tplExec) interface{} {
	return exec.sto.getValue(s.index)
}

func (s iName) String() string {
	return fmt.Sprintf("{ iName: %v, %v }", s.name, s.index)
}

func newValArifmetic(p *parser) *parseError {
	return p.initParseError(p.Line(), p.Pos(), "Error init arifmetic")
}
