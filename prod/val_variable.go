package prod

import "fmt"

func newValName(p *parser, line, pos int, key string) *parseError {
	name := iName{line: line, pos: pos}
	if p.varFlag {
		index, err := p.store.setVariable(string(key))
		if err != nil {
			return p.initParseError(line, pos, err)
		}
		name.index = index
	} else {
		name.index = p.store.initVariable(string(key))
	}
	p.stack.Push(name)
	return nil
}

type iName struct {
	line, pos int
	index     int
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
		//p
	}
	return p.initParseError(p.Line(), p.Pos(), fmt.Errorf("Error init set"))
}

func newValArifmetic(p *parser) *parseError {
	return p.initParseError(p.Line(), p.Pos(), fmt.Errorf("Error init arifmetic"))
}
