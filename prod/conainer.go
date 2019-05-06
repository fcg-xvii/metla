package prod

import (
	"fmt"
	"reflect"

	"github.com/fcg-xvii/lineman"
)

func newField(p *parser) *parseError {
	pos := position{p.tplName, p.Line(), p.LinePos()}
	if p.stack.Len() == 0 {
		return p.initParseError(pos.line, pos.pos, "Expected field owner")
	}
	list := []interface{}{p.stack.Pop()}
	stackLen := p.stack.Len()
	p.IncPos()
	p.fieldFlag = true
mainLoop:
	for !p.IsEndDocument() {
		fmt.Println("FIELD_FLAG", p.fieldFlag)
		p.PassSpaces()
		switch {
		case !lineman.CheckLetter(p.Char()) && !lineman.CheckNumber(p.Char()) && p.Char() != '(':
			fmt.Println(stackLen, p.stack.Len(), string(p.Char()))
			if stackLen != p.stack.Len()-1 {
				return p.initParseError(p.Line(), p.LinePos(), "Field parse error :: unexpected value")
			}
			list = append(list, p.stack.Pop())
			fmt.Println(p.Char())
			if p.Char() != '.' {
				break mainLoop
			} else {
				p.IncPos()
			}
		default:
			if err := p.initCodeVal(); err != nil {
				return err
			}
		}
	}
	p.fieldFlag = false
	p.stack.Push(&field{pos, list})
	fmt.Println("FIELD_ENDDD")
	return nil
}

type field struct {
	position
	list []interface{}
}

func (s *field) execType() execType {
	return execField
}

func (s *field) exec(exec *tplExec) *execError {
	fmt.Println("field_exec")
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
			if err := owner.(executer).exec(exec); err != nil {
				return err
			}
		case getter:
			exec.stack.Push(owner.(getter).get(exec))
		default:
			rOwner := reflect.ValueOf(owner)
			switch rOwner.Kind() {
			case reflect.Struct, reflect.Ptr:
				sVal := l[0]
				switch sVal.(type) {
				case static:
					if rOwner.Kind() == reflect.Ptr {
						rOwner = rOwner.Elem()
					}
					fieldName := sVal.(static).get(exec).(string)
					rVal := rOwner.FieldByName(fieldName)
					if rVal.Kind() == reflect.Invalid {
						return s.execError(fmt.Sprintf("Field %v not found in owner", fieldName))
					}
					exec.stack.Push(static{s.position, rVal.Interface()})
					l = l[1:]
				case *method:
					exec.stack.Push(owner)
					if err := sVal.(*method).exec(exec); err != nil {
						return err
					}
					l = l[1:]
				default:
					return pos.execError("!!!!!!!!!!!!!!!!") // Релазовать методы
				}
			case reflect.Map:
				if st, check := l[0].(static); check {
					rVal := reflect.ValueOf(st.get(exec))
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
						return st.execError(fmt.Sprintf("Map index [%v] not found", st.get(exec)))
					} else {
						exec.stack.Push(static{s.position, rResult.Interface()})
						l = l[1:]
					}
				} else {
					return l[0].(coordinator).execError("Fieldmap static value expected")
				}
			}
		}
	}
	//fmt.Println("TTTTTTTTTTTTTTTTTTTTTT", exec.stack.Peek())
	return nil
}
