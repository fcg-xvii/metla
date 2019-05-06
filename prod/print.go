package prod

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
	fmt.Println("ST_LENNN", exec.stack.Len())
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
	fmt.Println("STACKLEN", exec.stack.Len())
	for exec.stack.Len() > 0 {
		if err := exec.Write([]byte(fmt.Sprint(exec.stack.Pop().(getter).get(exec)) + " ")); err != nil {
			return err
		}
	}
	return exec.Write([]byte{'\n'})
}
