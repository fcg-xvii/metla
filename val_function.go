package metla

import (
	"fmt"
	"reflect"
)

func parseFuncArgs(p *parser) (err error) {
	info := p.infoRecordFromMark()
	for !p.IsEndDocument() {
		p.PassSpaces()
		//fmt.Println(string(p.Char()))
		switch p.Char() {
		case ',':
			p.IncPos()
		case ')':
			p.IncPos()
			return
		default:
			if _, err = initCodeVal(p); err != nil {
				return
			}
		}
	}
	return info.fatalError("Unexpected end of function exec")
}

func popVariadicArgs(exec *tplExec, fType reflect.Type, info *rawInfoRecord) (args []reflect.Value, err error) {
	var tmp []interface{}
	for exec.st.Len() > 0 {
		val := exec.st.Pop()
		if _, check := val.(*execMarker); check {
			break
		} else if vVal, check := val.(*variable); check {
			val = vVal.value
		}
		tmp = append([]interface{}{val}, tmp...)
	}
	for i, v := range tmp {
		var t reflect.Type
		if i < fType.NumIn()-1 {
			t = fType.In(i)
		} else {
			t = fType.In(fType.NumIn() - 1).Elem()
		}
		if t != reflect.TypeOf(v) && !reflect.TypeOf(v).ConvertibleTo(t) {
			err = info.fatalError(fmt.Sprintf("Coudn't convert function arg [%v], [%T] to [%v]", i, v, t))
			return
		}
		args = append(args, reflect.ValueOf(v).Convert(t))
	}
	return
}

func popStaticArgs(exec *tplExec, fType reflect.Type, info *rawInfoRecord) (args []reflect.Value, err error) {
	for exec.st.Len() > 0 {
		val := exec.st.Pop()
		if _, check := val.(*execMarker); check {
			break
		} else if vVal, check := val.(*variable); check {
			val = vVal.value
		}
		t := fType.In(len(args))
		if !reflect.TypeOf(val).ConvertibleTo(t) {
			err = info.fatalError(fmt.Sprintf("Coudn't convert function arg [%v], [%T] to [%v]", len(args), val, t))
			return
		}
		args = append([]reflect.Value{reflect.ValueOf(val).Convert(t)}, args...)
	}
	return
}

func popExecFunctionArgs(exec *tplExec, fType reflect.Type, info *rawInfoRecord) ([]reflect.Value, error) {
	if fType.IsVariadic() {
		return popVariadicArgs(exec, fType, info)
	} else {
		return popStaticArgs(exec, fType, info)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////

func newValFunction(name string, p *parser) (res interface{}, err error) {
	info := p.infoRecordFromPos()
	p.IncPos()
	p.stack.Push(&execMarker{name})
	if err = parseFuncArgs(p); err == nil {
		p.stack.Push(&valVariable{info, name})
		p.stack.Push(&execCommand{info, execFunction, 0})
	}
	return
}

func execFunction(exec *tplExec, info *rawInfoRecord) (err error) {
	var f reflect.Value
	fVal := exec.st.Pop()
	if fRes, check := fVal.(*variable); check {
		if fRes.value == nil {
			err = info.fatalError(fmt.Sprintf("Function [%v] not found", fRes.key))
			return
		}
		f = reflect.ValueOf(fRes.value)
	} else {
		f = reflect.ValueOf(fVal)
	}

	fType := f.Type()
	var args []reflect.Value
	if args, err = popExecFunctionArgs(exec, fType, info); err != nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = info.fatalError(fmt.Sprintf("Function call fatal error (panic) :: %v", r))
		}
	}()
	if res := f.Call(args); len(res) > 0 {
		for _, val := range res {
			exec.st.Push(val.Interface())
		}
	}
	return
}

//////////////////////////////////////////////////////////////////////////////////////////////

func newStaticFunction(fIface interface{}, p *parser) (res interface{}, err error) {
	info := p.infoRecordFromPos()
	p.IncPos()
	p.stack.Push(&execMarker{""})
	if err = parseFuncArgs(p); err == nil {
		p.stack.Push(fIface)
		p.stack.Push(&execCommand{info, execFunction, 0})
	}
	return
}
