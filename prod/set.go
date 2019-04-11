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
	return p.initParseError(p.Line(), p.LinePos(), fmt.Errorf("Unexpected endLine, expected ':' or '='"))
}

func newValSet(p *parser) *parseError {
	//stackLen := p.stack.Len()
	ex := set{}
	if p.Char() == ',' {
		if err := parseSetNames(p); err != nil {
			return err
		}
	}
	fmt.Println(p.stack)
	ex.names = p.stack.PopAll()
	fmt.Println("EXNNN", ex.names)
	fmt.Println("sto", p.store.list)
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
	fmt.Println("VALUES", ex.values)
	p.execList = append(p.execList, &ex)
	return nil
}

type set struct {
	names   []interface{}
	values  []interface{}
	uppdate bool
}

func (s *set) Exec(exec *tplExec) error {
	return fmt.Errorf("AAAAAAAAAAAAAAAAAAAAAAA")
}
