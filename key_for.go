package metla

import (
	"fmt"
)

func init() {
	keywords["for"] = newKeyFor
}

func newKeyFor(p *parser) (t token, err error) {
	fmt.Println("FOR_ERROR", string(p.EndLineContent()))
	var index token
	p.PassSpaces()
	if index, err = newValVariable(p); err != nil {
		return
	}
	p.PassSpaces()
	switch {
	case p.PosMatchSlice([]byte("in")):
		t, err = newForIn(index, p)
	default:
		err = fmt.Errorf("Keyword [for] parse error :: unexpected char after index [%c]", p.Char())
		return
	}
	Тут мы парсим тело цикла )))
	err = fmt.Errorf("For error")
	return
}

///////////////////////////////////////////////////////////////////////////////////////////////////////

func newForIn(index token, p *parser) (t token, err error) {

	fmt.Println("ELC", string(p.EndLineContent()))
	p.ForwardPos(2)
	var left token
	if left, err = initVal(p); err != nil {
		return
	}
	p.PassSpaces()
	if p.Char() != ':' {
		err = fmt.Errorf("Keyword [for] parse error :: unexpected char after left side '%c', ':' expected", p.Char())
	}
	p.IncPos()
	var right token
	if right, err = initVal(p); err != nil {
		return
	}
	fmt.Println("LRRRRR", left, right)
	t = &cycleIn{&cycle{index, nil}, left, right}
	return
}

type cycleIn struct {
	*cycle
	left, right token
}

type cycle struct {
	index  token
	childs []token
}

func (s *cycle) IsExecutable() bool {
	return true
}

func (s *cycle) String() string {
	return "[cycle...]"
}

func (s *cycle) execObject(*storage, *template) (execObject, error) {
	return nil, fmt.Errorf("Cycle is abstract object")
}

/*func rangeVal(p *parser) (token, error) {
	if !checkValInt(p.EndLineContent()) {
		return nil, errors.New("For (range) parse error :: Unexpected range token type, expected integer")
	}
	return newValInt(p)
}

func keywordFor(p *parser) (res token, err error) {
	var min, max []token
	p.PassSpaces()
	if min, err = rangeVal()
	return nil, nil
}

type keyForCount struct {
	min, max token
}*/
