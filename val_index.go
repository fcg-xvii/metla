package metla

import (
	"fmt"
	"reflect"
)

func newValIndex(p *parser) (res interface{}, err error) {
	p.IncPos()
	tmp := p.stack.Pop()
	res = &execCommand{p.infoRecordFromMark(), execIndex, "index", nil}
	p.stack.Push(res)
	p.stack.Push(tmp)
	for !p.IsEndDocument() {
		p.PassSpaces()
		if p.Char() == ']' {
			p.IncPos()
			//p.stack.Push(res)
			return
		} else {
			if _, err = initCodeVal(p); err != nil {
				return
			}
		}
	}
	err = fmt.Errorf("Expected close index char ']'")
	return
}

func execIndex(exec *tplExec, info *rawInfoRecord) (err error) {
	obj, index := exec.st.Pop(), exec.st.Pop()
	if index == nil {
		exec.showStack()
	}
	if _, check := obj.(*variable); check {
		if v, check := index.(*variable); check {
			index = v.value
		}
		exec.st.Push(indexVariable{obj.(*variable), index})
	} else {
		rObj := reflect.ValueOf(obj)
		switch rObj.Kind() {
		case reflect.Array, reflect.Slice, reflect.String:
			if v, check := index.(*variable); check {
				index = v.value
			}
			rIndex := reflect.ValueOf(index)
			switch rIndex.Kind() {
			case reflect.Int64, reflect.Int32, reflect.Int, reflect.Int16, reflect.Int8:
				exec.st.Push(rObj.Index(int(rIndex.Int())).Interface())
			default:
				err = info.fatalError("Array index integer value expected")
			}
		case reflect.Map:
			if v, check := index.(*variable); check {
				index = v.key
			}
			exec.st.Push(rObj.MapIndex(reflect.ValueOf(index)).Interface())
		default:
			err = info.fatalError("Expected variable, array or map")
		}
	}
	return
}
