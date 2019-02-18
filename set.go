package metla

import (
	"bytes"
	_ "errors"
	"fmt"
)

func initSet(prefix []byte, p *parser) (res *set, err error) {
	//fmt.Println("INIT_SET")
	// Парсим наименования переменных
	var names []string
	sNames := bytes.Split(prefix, []byte{','})
	names = make([]string, len(sNames))
	for i, v := range sNames {
		names[i] = string(bytes.TrimSpace(v))
	}
	fmt.Println(names)
	/////////////////////////////////////
	p.incPos()
	// Парсим значения (их может быть меньше, чем наименований из-за возвращаемых знаечений функции), но не больше
	var values []token
	for {
		var t token
		if t, err = p.parseToEndLine(); err == nil {
			if val, check := t.(value); check {
				values = append(values, val)
				p.passSpaces()
				if p.isEndLine() {
					break
				} else if p.char() == ',' {
					p.incPos()
				} else {
					err = fmt.Errorf("Unexpected symbol [%c], expected [',' or endline]", p.char())
					return
				}
			} else {
				err = fmt.Errorf("Set token error :: Value token expected...")
				fmt.Println(err)
				return
			}
		} else {
			return
		}
	}
	p.passSpaces()
	if !p.isEndLine() {
		err = fmt.Errorf("Unexpected symbol [%c]", p.char())
	} else {
		res = &set{names, values}
		p.incPos()
	}
	fmt.Println("SSSEETTT >> ", res, p.pos)
	return
}

type set struct {
	names  []string
	values []token
}

func (s *set) Data() ([]byte, error) {
	return nil, nil
}

func (s *set) Type() operatorType {
	return opSet
}

func (s *set) IsExecutable() bool { return true }
