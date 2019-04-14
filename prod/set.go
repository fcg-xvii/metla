package prod

import "fmt"

func parseSetNames(p *parser) *parseError {

	for !p.IsEndLine() {
		p.PassSpaces()
		switch p.Char() {
		case ',':
			p.IncPos()
		case '=':
			return nil
		}
		if err := p.initCodeVal(); err != nil {
			return err
		}
	}
	return p.initParseError(p.Line(), p.LinePos(), "Unexpected endLine, expected ':' or '='")
}

func newValSet(p *parser) *parseError {
	fmt.Println("NEW_VAL_SET")
	//stackLen := p.stack.Len()
	ex := set{}
	if p.Char() == ',' {
		if err := parseSetNames(p); err != nil {
			return err
		}
	}
	fmt.Println(p.stack)
	for p.stack.Len() > 0 {
		val := p.stack.Pop().(coordinator)
		if v, check := val.(setter); check {
			ex.names = append(ex.names, v)
		} else {
			return val.parseError("Expected setter token")
		}
	}
	//fmt.Println("EXNNN", ex.names)
	//fmt.Println("sto", p.store.list)
	if p.Char() == '=' {
		ex.uppdate = true
		p.IncPos()
	} else {
		p.ForwardPos(2)
	}

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
	ex.values = p.stack.PopAll()
	fmt.Println("NAMES", ex.names, p.stack.Len())
	fmt.Println("VALUES", ex.values, p.stack.Len())
	p.stack.Push(&ex)
	return nil
}

type set struct {
	*position
	names   []setter
	values  []interface{}
	uppdate bool
}

func (s *set) String() string {
	return "{ set }"
}

func (s *set) Exec(exec *tplExec) *execError {
	valsIndex := 0
	fmt.Println("NAMES", s.names)
	fmt.Println("VALS", s.values)
	for _, v := range s.names {
		//fmt.Println("!!!!!!", exec.stack.Len(), v)
		if exec.stack.Len() == 0 {
			switch s.values[valsIndex].(type) {
			case executer:
				//fmt.Println("EXECUTERT")
				if err := s.values[valsIndex].(executer).Exec(exec); err != nil {
					return err
				}
			case getter:
				//fmt.Println("GETTER")
				if err := v.Set(exec, s.values[valsIndex]); err != nil {
					return err
				}
			default:
				//fmt.Println("DEFAULT")
				return s.values[valsIndex].(coordinator).execError("Expected executer or getter token")
			}
			valsIndex++
		}
	}
	return nil
}
