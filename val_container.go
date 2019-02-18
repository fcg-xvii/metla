package metla

import (
	"fmt"
)

func init() {
	creators = append(creators, &valueCreator{
		checker:     checkValArray,
		constructor: newValArray,
	})
}

func checkValArray(src []byte) bool {
	return src[0] == '['
}

func checkValObject(src []byte) bool {
	return src[0] == '{'
}

func newValArray(p *parser) (res token, err error) {
	var vals []token
	p.incPos()
	var t token
loop:
	for !p.isEndDocument() {
		p.passEndLines()
		if t, err = initVal(p); err != nil {
			return nil, err
		} else {
			fmt.Println("<<<", t, ">>>")
			vals = append(vals, t)
			p.passSpaces()
			switch p.char() {
			case ',':
				p.incPos()
			case ']':
				fmt.Println("RRRRRR")
				p.incPos()
				break loop
			}
		}
	}
	p.passSpaces()
	if !p.isEndLine() {
		err = fmt.Errorf("Unexpected symbol [%c]", p.char())
	} else {
		res = &valArray{vals}
		p.incPos()
	}
	return
}

type valArray struct {
	vals []token
}

func (s *valArray) Val() interface{} {
	return s.vals
}

func (s *valArray) Data() (res []byte, err error) {
	//return []byte(s.val), nil
	return
}

func (s *valArray) String() string     { return "[array :: { ... }]" }
func (s *valArray) IsExecutable() bool { return false }

type varObject struct {
	val map[string]token
}
