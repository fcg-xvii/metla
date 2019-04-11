package prod

import "fmt"

func newValName(p *parser, line, pos int, key string) *parseError {
	name := iName{&position{p.tplName, line, pos}, 0}
	if p.varFlag {
		index, err := p.store.setVariable(string(key))
		if err != nil {
			return p.initParseError(line, pos, err)
		}
		name.index = index
	} else {
		name.index = p.store.initVariable(string(key))
	}
	p.stack.Push(&name)
	return nil
}

type iName struct {
	*position
	index int
}

func (s *iName) StorageIndex() int {
	return s.index
}

func newValArifmetic(p *parser) *parseError {
	return p.initParseError(p.Line(), p.Pos(), fmt.Errorf("Error init arifmetic"))
}
