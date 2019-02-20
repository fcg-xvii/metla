package metla

import (
	"fmt"
)

func init() {
	creators = append(creators, &valueCreator{
		checker:     checkValArray,
		constructor: newValArray,
	}, &valueCreator{
		checker:     checkValObject,
		constructor: newValObject,
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
	p.IncPos()
	var t token
loop:
	for !p.IsEndDocument() {
		p.PassEndLines()
		if t, err = initVal(p); err != nil {
			return nil, err
		} else {
			vals = append(vals, t)
			p.PassEndLines()
			switch p.Char() {
			case ',':
				p.IncPos()
			case ']':
				p.IncPos()
				break loop
			}
		}
	}
	res = &valArray{vals}
	fmt.Println(res)
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

func (s *valArray) String() string {
	return fmt.Sprintf("[array :: { %v }]", s.vals)
}

func (s *valArray) IsExecutable() bool { return false }

//////////////////////////////////////////////////////////////////

func newValObject(p *parser) (res token, err error) {
	m := make(map[string]token)
	var (
		check bool
		key   []byte
		val   token
	)
	p.IncPos()
	for !p.IsEndDocument() {
		// Парсим
		p.PassEndLines()
		if key, check = p.ReadName(); check {
			fmt.Println("KEY.....................", string(key))
			p.PassSpaces()
			if p.Char() != ':' {
				err = fmt.Errorf("Unexpected symbol '%c', expected ':'", p.Char())
				return
			}
			p.IncPos()
			if val, err = initVal(p); err != nil {
				return
			}
			m[string(key)] = val
			p.PassEndLines()
			switch p.Char() {
			case ',':
				p.IncPos()
			case '}':
				p.IncPos()
				break
			}
		} else {
			err = fmt.Errorf("Unexpected variable name...")
			return
		}
	}
	res = &valObject{m}
	fmt.Println("==================", res)
	return
}

type valObject struct {
	vals map[string]token
}

func (s *valObject) Val() interface{} {
	return s.vals
}

func (s *valObject) Data() (res []byte, err error) {
	//return []byte(s.val), nil
	return
}

func (s *valObject) String() string {
	return fmt.Sprintf("[object :: { %v }]", s.vals)
}

func (s *valObject) IsExecutable() bool { return false }
