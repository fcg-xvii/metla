package metla

import (
	_ "bytes"
	"errors"
	"fmt"
)

func initSet(prefix []byte, p *parser) (res *set, err error) {
	p.incPos()
	p.passSpaces()
	fmt.Println("PEFIX", prefix)
	val, err := p.parseToEndLine()
	fmt.Println(">>>>>>>>>>>>>>>>>", val, err)
	if !p.isEndLine() {
		err = fmt.Errorf("Unexpected symbol [%c]", p.char())
	}
	return nil, errors.New("Test error")
}

type set struct {
	varName string
	value   token
}

func (s *set) Data() ([]byte, error) {
	return nil, nil
}

func (s *set) Type() operatorType {
	return opSet
}
