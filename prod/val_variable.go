package prod

import "fmt"

type iName struct {
	line, pos int
	name      string
}

type set struct {
	names   []interface{}
	values  []interface{}
	uppdate bool
}

func parseSetNames(p *parser) *parseError {

	for !p.IsEndLine() {
		p.PassSpaces()
		switch p.Char() {
		case ',':
			p.IncPos()
		case ':', '=':
			return nil
		}
		if err := p.initCodeVal(); err != nil {
			return err
		}
	}
	return p.initParseError(p.Line(), p.LinePos(), fmt.Errorf("Unexpected endLine, expected ':' or '='"))
}

func newValSet(p *parser) *parseError {
	//fmt.Println(p.stack)
	stackLen := p.stack.Len()
	ex := set{}
	if p.Char() == ',' {
		if err := parseSetNames(p); err != nil {
			return err
		}
	}
	ex.names = p.stack.PopAllReverse()
	fmt.Println("EXNNN", ex.names)

	if p.Char() == '=' {
		ex.uppdate = true
		p.IncPos()
	} else {
		p.ForwardPos(2)
	}

	for !p.IsEndLine() {
		//p
	}
	return p.initParseError(p.Line(), p.Pos(), fmt.Errorf("Error init set"))
}

func newValArifmetic(p *parser) *parseError {
	return p.initParseError(p.Line(), p.Pos(), fmt.Errorf("Error init arifmetic"))
}
