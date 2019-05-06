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
	fmt.Println("SETTTTTTTTTTTTTTTTTTTTTT >>>>>>>>>>>>>>>>>>.")
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
			//ex.names = append(ex.names, v)
			ex.names = append(ex.names, v)
		} else {
			return val.parseError("Expected setter token")
		}
	}
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
	ex.values = p.stack.PopAllReverse()
	//fmt.Println("NAMES", ex.names, p.stack.Len())
	//fmt.Println("VALUES", ex.values, p.stack.Len())
	p.stack.Push(&ex)
	return nil
}

type set struct {
	position
	names   []setter
	values  []interface{}
	uppdate bool
}

func (s *set) String() string {
	return "{ set }"
}

func (s *set) execType() execType {
	return execSet
}

func (s *set) exec(exec *tplExec) *execError {
	fmt.Println("SET_EXEC...")
	valsIndex := 0
	fmt.Println("NAMES", s.names)
	fmt.Println("VALS", s.values)
	for _, v := range s.names {
		//fmt.Println("!!!!!!", exec.stack.Len(), v)
		if exec.stack.Len() == 0 {
			switch s.values[valsIndex].(type) {
			case executer:
				if err := s.values[valsIndex].(executer).exec(exec); err != nil {
					return err
				}
			case getter:
				exec.stack.Push(s.values[valsIndex])
			default:
				return s.values[valsIndex].(coordinator).execError("Expected executer or getter token")
			}
			valsIndex++
		}
		if err := v.set(exec, exec.stack.Pop()); err != nil {
			return err
		}
	}
	fmt.Println("SET_STACK", exec.stack)
	return nil
}
