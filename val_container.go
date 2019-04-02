package metla

import (
	"fmt"
	_ "io"
	_ "reflect"
)

func newValArray(p *parser) (res interface{}, err error) {
	p.IncPos()
	res = &execCommand{p.infoRecordFromMark(), initArray, 0}
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
	res = &execCommand{p.infoRecordFromMark(), initObject, 0}
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
		if _, err = initCodeVal(p); err != nil {
			return
		}
		p.stack.Push(&execMarker{"pair-split"})
	}
	return nil, p.positionError("Unclosed object")
}

func initObject(exec *tplExec, info *rawInfoRecord) (err error) {
	fmt.Println("INIT_OBJECT")
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
	tmp := p.stack.Pop()
	p.stack.Push(&execMarker{"endfield"})
	p.stack.Push(tmp)
	var val interface{}
	for !p.IsEndLine() {
		p.PassSpaces()
		if p.Char() != '.' {
			break
		}
		p.IncPos()
		if _, err = initCodeVal(p); err != nil {
			return
		}
	}
	res = &execCommand{p.infoRecordFromMark(), execField, 0}
	p.stack.Push(res)
	return
}

func execField(exec *tplExec, info *rawInfoRecord) (err error) {
	owner := exec.st.Pop()
	for {

	}
	err = info.fatalError("EXEC FIELD error")
	return
}
