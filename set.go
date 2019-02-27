package metla

import (
	"bytes"
	"errors"
	_ "errors"
	"fmt"
	"io"
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
		res = &set{names, values, false}
		p.IncPos()
	}
	fmt.Println("SSSEETTT >>", names, values)
	return
}

type set struct {
	names  []string
	values []token
	create bool
}

func (s *set) execObject(sto *storage, tpl *template) (res execObject, err error) {
	vars, vals := make([]*variable, len(s.names)), make([]execObject, len(s.values))
	if s.create {
		for i, v := range s.names {
			if vars[i], err = sto.appendValue(v, nil); err != nil {
				return
			}
		}
	} else {
		var check bool
		for i, v := range s.names {
			if vars[i], check = sto.findVariable(v); !check {
				if vars[i], err = sto.appendValue(v, nil); err != nil {
					return
				}
			}
		}
	}
	for i, v := range s.values {
		if vals[i], err = v.execObject(sto, tpl); err != nil {
			return
		}
	}
	return
}

func (s *set) Data(w io.Writer) error {
	return nil
}

func (s *set) Type() operatorType {
	return opSet
}

func (s *set) IsExecutable() bool { return true }

func (s *set) String() string { return "[set {}}" }

//////////////////////////////////////////////////////////////////////

type execObjectSet struct {
	vars   []*variable
	values []execObject
}

func (s *execObjectSet) Data(w io.Writer) (err error) {
	count := 0
	if len(s.values) == 1 {
		_, err = s.setupVariable(0, 0)
	} else {
		var varIndex, count int
		for i, v := range s.values {
			if count, err = s.setupVariable(varIndex, i); err == nil && varIndex+count < len(s.vars) {
				varIndex += count
			} else {
				return
			}
		}
	}
	return
}

func (s *execObjectSet) setupVariable(varIndex, valIndex int) (count int, err error) {
	val := s.values[valIndex]
	if val.ValSingle() {
		count = 1
		s.vars[varIndex] = val
	} else {
		var index int
		for i, v := range val.Vals() {
			if index = varIndex + i; index < len(s.vars) {
				s.vars[varIndex+i].value = v
				count++
			} else {
				return
			}
		}
	}
}

/*func (s *execObjectSet) Data(w io.Writer) (err error) {
	count := 0
	if len(s.values) == 1 {
		_, err = s.setupVariable(0, 0)
	} else {
		var varIndex, count int
		for i, v := range s.values {
			if count, err = s.setupVariable(varIndex, i); err == nil && varIndex+count < len(s.vars) {
				varIndex += count
			} else {
				return
			}
		}
	}
	return
}

func (s *execObjectSet) setupVariable(varIndex, valIndex int) (count int, err error) {
	val := s.values[valIndex]
	count = val.ValsCount()
	if count == 0 {
		err = errors.New("Set exec error :: token values result count is 0. Count must be > 0")
	} else if count == 1 {
		s.vars[varIndex].value = val.Val()
	} else {
		var index int
		for i, v := range val.Vals() {
			if index = varIndex + i; index < len(s.vars) {
				s.vars[varIndex+i].value = v
			} else {
				return
			}
		}
	}
	return
}*/
