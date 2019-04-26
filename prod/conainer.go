package prod

import (
	"fmt"
	"reflect"
)

func newField(p *parser) *parseError {
	pos := position{p.tplName, p.Line(), p.LinePos()}
	if p.stack.Len() == 0 {
		return p.initParseError(pos.line, pos.pos, "Expected field owner")
	}
	list := []interface{}{p.stack.Pop()}
	p.IncPos()
	p.fieldFlag = true

mainLoop:
	for !p.IsEndDocument() {
		if err := p.initCodeVal(); err != nil {
			return err
		}
		list = append(list, p.stack.Pop())
		/*p.PassSpaces()
		if p.Char() != '.' {
			break
		}
		p.IncPos()*/
		p.PassSpaces()
		switch p.Char() {
		case '.':
			fmt.Println("POINT")
			p.IncPos()
		default:
			break mainLoop
		}
	}
	p.fieldFlag = false
	p.stack.Push(&field{&pos, list})
	return nil
}

type field struct {
	*position
	list []interface{}
}

func (s *field) Exec(exec *tplExec) *execError {
	pos, stackLen := s.position, exec.stack.Len()
	exec.stack.Push(s.list[0])
	l := s.list[1:]
	for len(l) > 0 {
		//fmt.Println(exec.stack.Len(), stackLen, len(l))
		if exec.stack.Len()-1 > stackLen {
			return pos.execError("field item returned more that one value")
		}
		owner := exec.stack.Pop()
		//fmt.Println(owner)
		switch owner.(type) {
		case executer:
			if err := owner.(executer).Exec(exec); err != nil {
				return err
			}
		case getter:
			exec.stack.Push(owner.(getter).Get(exec))
		default:
			rOwner := reflect.ValueOf(owner)
			if rOwner.Kind() == reflect.Ptr {
				rOwner = rOwner.Elem()
			}
			switch rOwner.Kind() {
			case reflect.Struct:
				sVal := l[0]
				switch sVal.(type) {
				case *static:
					fieldName := sVal.(*static).Get(exec).(string)
					rVal := rOwner.FieldByName(fieldName)
					if rVal.Kind() == reflect.Invalid {
						return s.execError(fmt.Sprintf("Field %v not found in owner", fieldName))
					}
					exec.stack.Push(rVal.Interface())
					l = l[1:]
				default:
					return pos.execError("!!!!!!!!!!!!!!!!") // Релазовать методы
				}
			case reflect.Map:
				if st, check := l[0].(*static); check {
					rVal := reflect.ValueOf(st.Get(exec))
					rType, keyType := rVal.Type(), reflect.TypeOf(owner).Key()
					if rType != keyType {
						if !rType.ConvertibleTo(keyType) {
							return st.execError(fmt.Sprintf("Types map key and field is not comparable [%v, %v]", keyType, rType))
						} else {
							rVal = rVal.Convert(keyType)
						}
					}
					rResult := rOwner.MapIndex(rVal)
					if rResult.Kind() == reflect.Invalid {
						return st.execError(fmt.Sprintf("Map index [%v] not found", st.Get(exec)))
					} else {
						exec.stack.Push(rResult.Interface())
						l = l[1:]
					}
				} else {
					return l[0].(coordinator).execError("Fieldmap static value expected")
				}
			}
		}
	}
	return nil
}
