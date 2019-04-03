package metla

import (
	"fmt"
	_ "io"
	"reflect"
)

func newValArray(p *parser) (res interface{}, err error) {
	p.IncPos()
	res = &execCommand{p.infoRecordFromMark(), initArray, "init-array"}
	p.stack.Push(res)
	for !p.IsEndDocument() {
		p.PassSpaces()
		switch p.Char() {
		case ',', '\n':
			p.IncPos()
		case ']':
			p.stack.Push(&execMarker{"endarr"})
			p.IncPos()
			return
		default:
			if _, err = initCodeVal(p); err != nil {
				return
			}
		}
	}
	return nil, p.positionError("Unclosed array, expected ']'")
}

func initArray(exec *tplExec, info *rawInfoRecord) (err error) {
	var arr []interface{}
	for exec.st.Len() > 0 {
		item := exec.st.Pop()
		if _, check := item.(*execMarker); check {
			exec.st.Push(arr)
			return
		} else {
			arr = append(arr, item)
		}
	}
	err = info.fatalError("INIT_ARRAY_ERROR")
	return
}

func newValObject(p *parser) (res interface{}, err error) {
	res = &execCommand{p.infoRecordFromMark(), initObject, "init-object"}
	p.stack.Push(res)
	p.IncPos()
loop:
	for !p.IsEndDocument() {
		p.PassSpaces()
		switch p.Char() {
		case ',', '\n':
			p.IncPos()
			continue loop
		case '}':
			p.stack.Push(&execMarker{"endobj"})
			p.IncPos()
			return
		}
		if p.Char() == ',' || p.Char() == '\n' {
			p.IncPos()
			continue
		}
		p.PassSpaces()
		if fieldName, check := p.ReadName(); !check {
			return nil, p.positionError("Object init error :: expected field name")
		} else {
			p.stack.Push(string(fieldName))
		}
		if p.Char() != ':' {
			return nil, p.positionError("Object init error :: expected value splitter - ':'")
		}
		p.IncPos()
		for !p.IsEndDocument() {
			p.PassSpaces()
			if p.Char() == ',' || p.Char() == '}' || p.Char() == '\n' {
				break
			}
			if _, err = initCodeVal(p); err != nil {
				return
			}
		}
		p.stack.Push(&execMarker{"pair-split"})
	}
	return nil, p.positionError("Unclosed object")
}

func initObject(exec *tplExec, info *rawInfoRecord) (err error) {
	//fmt.Println("INIT_OBJECT")
	m, pairAccepted, fieldName := make(map[string]interface{}), false, ""
	var val interface{}
loop:
	for exec.st.Len() > 0 {
		item := exec.st.Pop()
		switch item.(type) {
		case string:
			if !pairAccepted {
				if len(fieldName) == 0 {
					fieldName = item.(string)
				} else if !pairAccepted {
					val = item
					pairAccepted = true
				}
			}
		case *execMarker:
			marker := item.(*execMarker)
			if marker.name == "endobj" {
				break loop
			} else if !pairAccepted {
				return info.fatalError("Object pair value not arrived")
			}
			m[fieldName] = val
			fieldName, val, pairAccepted = "", nil, false
		default:
			if !pairAccepted {
				val = item
				pairAccepted = true
			}
		}
	}
	exec.st.Push(m)
	return
}

func newValField(p *parser) (res interface{}, err error) {
	//fmt.Println("VAL_FIELD.........")
	p.fieldFlag = true
	tmp := p.stack.Pop()
	p.stack.Push(&execCommand{p.infoRecordFromMark(), execFieldEnd, "field-end"})
	p.stack.Push(tmp)
	var val interface{}
	for !p.IsEndLine() {
		//fmt.Println("STEP....")
		p.PassSpaces()
		//fmt.Println("CHAR", string(p.Char()))
		if p.Char() != '.' {
			p.fieldFlag = false
			break
		}
		//fmt.Println("Continue")
		p.IncPos()
		if val, err = initCodeVal(p); err != nil {
			return
		} else if v, check := val.(*valVariable); check {
			p.stack.Pop()
			p.stack.Push(v.name)
		}
	}
	p.stack.Push(&execCommand{p.infoRecordFromMark(), execFieldStart, "field-start"})
	return
}

func execFieldStart(exec *tplExec, info *rawInfoRecord) (err error) {
	//fmt.Println("FIELD_START", exec.fieldLayout)
	exec.fieldLayout++
	exec.st.Push(&execMarker{"endField"})
	return
}

func execFieldEnd(exec *tplExec, info *rawInfoRecord) (err error) {
	//fmt.Println("EXEC_FIELD_END", exec.fieldLayout, exec.index)
	exec.fieldLayout--
	for exec.st.Len() > 0 {
		l := exec.st.Pop()
		if v, check := l.(*variable); check {
			l = v.value
		}
		switch exec.st.Peek().(type) {
		case *execMarker:
			if exec.st.Peek().(*execMarker).name == "endField" {
				exec.st.Pop()
				exec.st.Push(l)
				return
			}
		case *execCommand:
			ex := exec.st.Pop().(*execCommand)
			exec.st.Push(l)
			if err = ex.method(exec, ex.posInfo()); err != nil {
				return
			}
		case string:
			left := reflect.ValueOf(l)
			if left.Kind() == reflect.Ptr {
				left = left.Elem()
			}
			var val reflect.Value
			switch left.Kind() {
			case reflect.Struct:
				val = left.FieldByName(exec.st.Pop().(string))
			case reflect.Map:
				val = left.MapIndex(reflect.ValueOf(exec.st.Pop()))
			default:
				return info.fatalError(fmt.Sprintf("Field expected type struct or map, not %v", left.Kind()))
			}
			if !val.IsValid() {
				exec.st.Push(nil)
			} else {
				exec.st.Push(val.Interface())
			}
		default:
			{
				err = info.fatalError("Field uncnown error...")
				return
			}
		}

	}
	return
}
