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
	}
	p.stack.Push(&echoln{&pos, p.stack.PopAll()})
	return nil
}

type echoln struct {
	*position
	items []interface{}
}

func (s *echoln) Exec(exec *tplExec) *execError {
	fmt.Println("EXEC_ECHOLN", len(s.items), s.items[0])
	for _, v := range s.items {
		switch v.(type) {
		case executer:
			if err := v.(executer).Exec(exec); err != nil {
				return err
			}
			for exec.stack.Len() > 0 {
				if err := exec.Write([]byte(fmt.Sprint(exec.stack.Pop()) + " ")); err == nil {
					return err
				}
			}
		case getter:
			fmt.Println("GETTER")
			if err := exec.Write([]byte(fmt.Sprint(v.(getter).Get(exec)))); err != nil {
				return err
			}
		}
	}
	return nil
}
