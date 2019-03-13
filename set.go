package metla

import (
	_ "errors"
	"fmt"
	"io"
	"reflect"
)

func initSet(p *parser) (res *set, err error) {
	// Парсим наименования переменных
	var (
		vars []string
		v    token
	)
	// Пока стек не пустой
	for p.codeStack.Peek() != nil {
		v = p.codeStack.Pop().(token) // Достаём токен из стека
		// Если токен не является наименованием переменной, возвращаем ошибку или добавляем наименование переменной в список
		if val, check := v.(*valVariable); !check {
			err = v.fatalError(fmt.Sprintf("Unexpected variable type %v, variable name expected", v))
			return
		} else {
			vars = append([]string{val.name}, vars...)
		}
	}
	// Если список наименований переменных пуст, возвращаем ошибку
	if len(vars) == 0 {
		err = p.positionError("Set left side is empty")
		return
	}
	/////////////////////////////////////
	p.IncPos()
	// Парсим значения (их может быть меньше, чем наименований из-за возвращаемых знаечений функции), но не больше
	var values []token
	if _, err = p.parseToEndLine(); err == nil {
		for p.codeStack.Peek() != nil {
			values = append([]token{p.codeStack.Pop().(token)}, values...)
		}
	}
	p.PassSpaces()
	if !p.IsEndLine() {
		err = p.positionError(fmt.Sprintf("Unexpected symbol [%c]", p.Char()))
	} else {
		res = &set{p.infoRecordFromMark(), vars, values, false}
	}
	return
}

type set struct {
	*rawInfoRecord
	names  []string
	values []token
	create bool
}

func (s *set) posInfo() *rawInfoRecord { return s.rawInfoRecord }

func (s *set) execObject(sto *storage, tpl *template, parent execObject) (res execObject, err error) {
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
		if vals[i], err = v.execObject(sto, tpl, nil); err != nil {
			return
		}
	}
	res = &execObjectSet{s.rawInfoRecord, &eventExec{parent}, vars, vals}
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
	*rawInfoRecord
	*eventExec
	vars   []*variable
	values []execObject
}

func (s *execObjectSet) Data(w io.Writer) (err error) {
	fmt.Println("SET_DATA...", s.vars, s.values)
	if len(s.values) == 1 {
		_, err = s.setupVariable(0, 0)
	} else {
		var varIndex, count int
		for i, _ := range s.values {
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
	fmt.Println("SETUP_INDEX", s.vars, s.values)
	val := s.values[valIndex]
	if val.ValSingle() {
		count = 1
		var iface interface{}
		if iface, err = val.Val(); err == nil {
			s.vars[varIndex].value = iface
		}
	} else {
		var (
			index int
			vals  []interface{}
		)
		if vals, err = val.Vals(); err == nil {
			for i, v := range vals {
				if index = varIndex + i; index < len(s.vars) {
					s.vars[varIndex+i].value = v
					count++
				} else {
					return
				}
			}
		}
	}
	return
}

func (s *execObjectSet) IsNil() bool {
	return false
}

func (s *execObjectSet) String() string {
	return "[set...]"
}

func (s *execObjectSet) Type() reflect.Kind           { return reflect.Invalid }
func (s *execObjectSet) Val() (interface{}, error)    { return nil, nil }
func (s *execObjectSet) Vals() ([]interface{}, error) { return nil, nil }
func (s *execObjectSet) ValSingle() bool              { return true }
