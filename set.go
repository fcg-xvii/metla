package metla

import (
	"bytes"
	_ "errors"
	"fmt"
)

func initSet(prefix []byte, p *parser) (res *set, err error) {
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
		p.passSpaces() // Убираем пробелы ме
		if t, err = p.parseToEndLine(); err == nil {
			if val, check := t.(value); check {
				values = append(values, val)
				fmt.Println(values)
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
				err = fmt.Errorf("Token value expected...")
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
