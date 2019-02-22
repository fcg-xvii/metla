package metla

import (
	"bytes"
	_ "errors"
	"fmt"
)

func initSet(prefix []byte, p *parser) (res *set, err error) {
	// Парсим наименования переменных
	fmt.Println("PREFIX", prefix)
	var names []string
	sNames := bytes.Split(prefix, []byte{','})
	names = make([]string, len(sNames))
	for i, v := range sNames {
		names[i] = string(bytes.TrimSpace(v))
	}
	fmt.Println(names)
	/////////////////////////////////////
	p.IncPos()
	// Парсим значения (их может быть меньше, чем наименований из-за возвращаемых знаечений функции), но не больше
	var values []token
	for {
		var t token
		if t, err = p.parseToEndLine(); err == nil {
			if val, check := t.(value); check {
				values = append(values, val)
				p.PassSpaces()
				if p.IsEndLine() {
					break
				} else if p.Char() == ',' {
					p.IncPos()
				} else {
					err = fmt.Errorf("Unexpected symbol [%c], expected [',' or endline]", p.Char())
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
	p.PassSpaces()
	if !p.IsEndLine() {
		err = fmt.Errorf("Unexpected symbol [%c]", p.Char())
	} else {
		res = &set{names, values}
		p.IncPos()
	}
	fmt.Println("SSSEETTT >>", names, values)
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

func (s *set) String() string { return "[set {}}" }
