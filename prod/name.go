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

func newValArifmetic(p *parser) *parseError {
	return p.initParseError(p.Line(), p.Pos(), fmt.Errorf("Error init arifmetic"))
}
