package metla

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

type splitter struct {
	position
}

func (s *splitter) execType() execType {
	return execSplitter
}

func (s *splitter) exec(exec *tplExec) *execError {
	exec.stack.PopAll()
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
	//fmt.Println("EXEC_METHOD")
	pos, stackLen := s.position, exec.stack.Len()
	exec.stack.Push(s.list[0])
	l := s.list[1:]
	for len(l) > 0 {
		//fmt.Println(exec.stack.Len(), stackLen, len(l))
		if exec.stack.Len()-1 > stackLen {
			return pos.execError("field item returned more that one value")
		}
		owner := exec.stack.Pop()
		switch owner.(type) {
		case executer:
			if err := owner.(executer).exec(exec); err != nil {
				return err
			}
		case getter:
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
						exec.stack.Push(static{s.position, nil})
						return nil
						//return st.execError(fmt.Sprintf("Map index [%v] not found", st.get(exec)))
					} else {
						switch iface := rResult.Interface(); iface.(type) {
						case static:
							exec.stack.Push(iface)
						case iName:
							exec.stack.Push(static{s.position, iface.(getter).get(exec)})
						default:
							exec.stack.Push(static{s.position, rResult.Interface()})
						}
						l = l[1:]
					}
				} else {
					return l[0].(coordinator).execError("Fieldmap static value expected")
				}
			default:
				//return s.execError("map, array or slice token expected")
				exec.stack.Push(static{s.position, nil})
				return nil
			}

		}
	}
	//fmt.Println(exec.stack.Peek())
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
}

func (s objIndex) execType() execType {
	return execIndex
}

///////////////////////////////////////////////////////////

func newArray(p *parser) *parseError {
	var list []coordinator
	pos, stackLen := position{p.tplName, p.Line(), p.LinePos()}, p.stack.Len()

	flushArg := func() *parseError {
		if stackLen+1 != p.stack.Len() {
			return p.initParseError(p.Line(), p.LinePos(), "Expected one value")
		}
		list = append(list, p.stack.Pop().(coordinator))
		return nil
	}
	valArrived := false
	p.IncPos()
	for !p.IsEndDocument() {
		p.PassSpaces()
		switch p.Char() {
		case ']':
			if valArrived {
				flushArg()
			}
			p.stack.Push(array{pos, list})
			p.IncPos()
			return nil
		case ',', '\n', ';':
			if valArrived {
				if err := flushArg(); err != nil {
					return err
				}
				valArrived = false
			}
			p.IncPos()
		default:
			valArrived = true
			if err := p.initCodeVal(); err != nil {
				return err
			}
		}
	}
	return pos.parseError("Unclosed array token")
}

type array struct {
	position
	vals []coordinator
}

func (s array) execType() execType {
	return execArray
}

func (s array) exec(exec *tplExec) *execError {
	arr := make([]interface{}, len(s.vals))
	for i, v := range s.vals {
		if obj, err := execOneReturn(v, exec); err != nil {
			return err
		} else {
			arr[i] = obj
		}
	}
	exec.stack.Push(static{s.position, arr})
	return nil
}

////////////////////////////////////////////////////////////

func newObject(p *parser) *parseError {
	rMap := make(map[string]coordinator)
	key, stackLen, pos := "", p.stack.Len(), position{p.tplName, p.Line(), p.LinePos()}

	flushValue := func() *parseError {
		if key == "" && stackLen == p.stack.Len() {
			return nil
		} else if key == "" {
			return p.initParseError(p.Line(), p.LinePos(), "Unexpected value, name expected")
		} else if p.stack.Len() != stackLen+1 {
			return p.initParseError(p.Line(), p.LinePos(), "Expected single token")
		}
		rMap[key], key = p.stack.Pop().(coordinator), ""
		return nil
	}

	p.IncPos()
	p.fieldFlag = true
	for !p.IsEndDocument() {
		p.PassSpaces()
		switch ch := p.Char(); ch {
		case ':':
			p.fieldFlag = false
			if key != "" {
				return p.initParseError(p.Line(), p.LinePos(), "Unexpected ':' splitter, value expected")
			}
			if p.stack.Len() != stackLen+1 {
				return p.initParseError(p.Line(), p.LinePos(), "Expected string token")
			}
			switch crd := p.stack.Pop().(coordinator); crd.(type) {
			case static:
				var check bool
				if key, check = crd.(static).get(nil).(string); !check {
					return p.initParseError(p.Line(), p.LinePos(), "Expected string token")
				}
			//case iName:
			//key = crd.(iName).name
			default:
				//fmt.Println("CRD", crd)
				return p.initParseError(p.Line(), p.LinePos(), "Expected string token")
			}
			/*if g, check := p.stack.Pop().(static); !check {
				return p.initParseError(p.Line(), p.LinePos(), "Expected string token 2")
			} else if key, check = g.get(nil).(string); !check {
				return p.initParseError(p.Line(), p.LinePos(), "Expected string token 3")
			}*/
			p.IncPos()

		case '\n', ',', '}':
			if err := flushValue(); err != nil {
				return err
			}
			p.IncPos()
			if ch == '}' {
				p.stack.Push(object{pos, rMap})
				p.fieldFlag = false
				return nil
			}
			p.fieldFlag = true
		default:
			if err := p.initCodeVal(); err != nil {
				return err
			}
		}
	}
	return pos.parseError("Unclosed object token")
}

type object struct {
	position
	vals map[string]coordinator
}

func (s object) execType() execType {
	return execObject
}

func (s object) exec(exec *tplExec) *execError {
	res := make(map[string]interface{})
	for k, v := range s.vals {
		if obj, err := execOneReturn(v, exec); err != nil {
			return err
		} else {
			res[k] = obj
		}
	}
	exec.stack.Push(static{s.position, res})
	return nil
}
