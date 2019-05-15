package metla

import (
	"fmt"
)

func init() {
	keywords["echoln"] = newEcholn
}

func newEcholn(p *parser) *parseError {
	pos := position{p.tplName, p.Line(), p.LinePos() - 5}
	for !p.IsEndLine() {
		if err := p.initCodeVal(); err != nil {
			return err
		}
		p.PassSpaces()
		if p.Char() == ',' {
			p.IncPos()
		}
		p.PassSpaces()
	}
	p.stack.Push(echoln{pos, p.stack.PopAllReverse()})
	return nil
}

type echoln struct {
	position
	items []interface{}
}

func (s echoln) execType() execType {
	return execEcholn
}

func (s echoln) exec(exec *tplExec) *execError {
	for _, v := range s.items {
		switch v.(type) {
		case executer:
			if err := v.(executer).exec(exec); err != nil {
				return err
			}
		default:
			exec.stack.Push(v)
		}
	}
	for exec.stack.Len() > 0 {
		if err := exec.Write([]byte(fmt.Sprint(exec.stack.Pop().(getter).get(exec)) + " ")); err != nil {
			return err
		}
	}
	return exec.Write([]byte{'\n'})
}

func newPrint(p *parser) *parseError {
	pos, stackLen, closeTag := position{p.tplName, p.Line(), p.LinePos()}, p.stack.Len(), []byte("}}")
	p.ForwardPos(2)
	for !p.IsEndDocument() {
		if err := p.initCodeVal(); err != nil {
			return err
		}
		p.PassSpaces()
		//fmt.Println(string(p.Char()))
		if p.PosMatchSlice(closeTag) {
			if stackLen+1 != p.stack.Len() {
				return pos.parseError("Expected single token")
			}
			p.stack.Push(print{pos, p.stack.Pop().(coordinator)})
			p.ForwardPos(2)
			return nil
		}
	}
	return pos.parseError("Unclosed print tag")
}

type print struct {
	position
	item coordinator
}

func (s print) execType() execType {
	return execPrint
}

func (s print) exec(exec *tplExec) *execError {
	//stackLen := exec.stack.Len()
	switch s.item.(type) {
	case getter:
		//exec.stack.Push(s.item.(getter).get(exec))
		if err := exec.Write([]byte(fmt.Sprint(s.item.(getter).get(exec)))); err != nil {
			return nil
		}
	case executer:
		if err := s.item.(executer).exec(exec); err != nil {
			return nil
		}
		for exec.stack.Len() > 0 {
			if err := exec.Write([]byte(fmt.Sprint(exec.stack.Pop().(getter).get(exec)) + " ")); err != nil {
				return nil
			}
		}
	}
	return nil
}
