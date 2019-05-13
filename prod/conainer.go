package prod

import (
	"fmt"
	"reflect"

	"github.com/fcg-xvii/lineman"
)

func newField(p *parser) *parseError {
	//fmt.Println("NEW_FIELD")
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
		//fmt.Println("FIELD_FLAG", p.fieldFlag)
		p.PassSpaces()
		switch {
		case !lineman.CheckLetter(p.Char()) && !lineman.CheckNumber(p.Char()) && p.Char() != '(':
			//fmt.Println(stackLen, p.stack.Len(), string(p.Char()))
			if stackLen != p.stack.Len()-1 {
				return p.initParseError(p.Line(), p.LinePos(), "Field parse error :: unexpected value")
			}
			list = append(list, p.stack.Pop())
			//fmt.Println(p.Char())
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
	//fmt.Println("FIELD_ENDDD", p.stack.Peek())
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
	//fmt.Println("field_exec")
	pos, stackLen := s.position, exec.stack.Len()
	exec.stack.Push(s.list[0])
	fmt.Println(exec.stack.Peek())
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
			fmt.Println("GETTER", owner)
			val := owner.(getter).get(exec)
			if val == nil {
				return owner.(coordinator).execError("Invalid field owner")
			}
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
				fmt.Println("MMMAAAA")
				if st, check := l[0].(static); check {
					fmt.Println("IF")
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
					fmt.Println(">>>>>>>>>>>>>>>>>>>>", rResult)
					if rResult.Kind() == reflect.Invalid {
						return st.execError(fmt.Sprintf("Map index [%v] not found", st.get(exec)))
					} else {
						fmt.Sprintf("STTTTT... %T", rResult.Interface())
						switch iface := rResult.Interface(); iface.(type) {
						case static:
							exec.stack.Push(iface)
						case iName:
							fmt.Println("INNNNAMMMMMMEEEEEEEEE")
							exec.stack.Push(static{s.position, iface.(getter).get(exec)})
							fmt.Println("><<<<<<<<<<<<<<<<<<<<<<<<<", exec.stack.Peek())
						default:
							exec.stack.Push(static{s.position, rResult.Interface()})
						}
						l = l[1:]
					}
				} else {
					return l[0].(coordinator).execError("Fieldmap static value expected")
				}
			}
		}
	}
	fmt.Println(exec.stack.Peek())
	return nil
}

///////////////////////////////////////////////////////////

func newIndex(p *parser) *parseError {
	r := objIndex{position: position{p.tplName, p.Line(), p.LinePos()}, owner: p.stack.Pop().(coordinator)}
	p.IncPos()
	for !p.IsEndDocument() {
		if err := p.initCodeVal(); err != nil {
			return err
		}
		p.PassSpaces()
		if p.Char() == ']' {
			r.index = p.stack.Pop().(coordinator)
			p.IncPos()
			p.stack.Push(r)
			return nil
		}
	}
	return p.initParseError(0, 0, "Unexpected end document")
}

type objIndex struct {
	position
	owner coordinator
	index coordinator
}

func (s objIndex) exec(exec *tplExec) *execError {
	owner, err := execOneReturn(s.owner, exec)
	if err != nil {
		return err
	}
	index, err := execOneReturn(s.index, exec)
	if err != nil {
		return err
	}
	oVal, iVal := reflect.ValueOf(owner), reflect.ValueOf(index)
	var rVal reflect.Value
	switch oVal.Kind() {
	case reflect.Map:
		{
			if index == nil {
				rVal = reflect.ValueOf(nil)
			} else {
				oType, iType := oVal.Type(), iVal.Type()
				if oType.Key() != iType {
					if !iType.ConvertibleTo(oType.Key()) {
						return s.index.execError(fmt.Sprintf("Expected index type %v", oType.Elem()))
					} else {
						iVal = iVal.Convert(oType.Key())
					}
				}
				rVal = oVal.MapIndex(iVal)
			}
		}
	case reflect.Slice, reflect.Array:
		{
			var index int
			switch iVal.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				index = int(iVal.Int())
			case reflect.Float32, reflect.Float64:
				index = int(iVal.Float())
			}
			rVal = oVal.Index(index)
		}
	default:
		return s.owner.execError("Expected array, slice or map object")
	}
	if rVal.Kind() == reflect.Invalid {
		exec.stack.Push(static{s.owner.getPosition(), nil})
	} else {
		switch iface := rVal.Interface(); iface.(type) {
		case coordinator:
			exec.stack.Push(iface)
		default:
			exec.stack.Push(static{s.owner.getPosition(), iface})
		}
	}
	return nil
	//fmt.Println(owner, index)

	/*switch s.owner.(type) {
	case getter:
		rVal = reflect.ValueOf(s.owner.(getter).get(exec))
	case executer:
		stackLen := exec.stack.Len()
		if err := s.owner.(executer).exec(exec); err != nil {
			return err
		}
		if stackLen+1 != exec.stack.Len() {
			return s.owner.execError("Expected one return value")
		}
		rVal = reflect.ValueOf(exec.stack.Pop())
	}
	switch rVal.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array:
		{
			var iVal reflect.Value
			switch s.index.(type) {
				case
			}
		}
	default:
		return s.owner.execError("Expected map, slice or array type")
	}*/
}

func (s objIndex) execType() execType {
	return execIndex
}
