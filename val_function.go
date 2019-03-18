package metla

import (
	"fmt"
	"reflect"
)

func newValFunction(name string, p *parser) (res interface{}, err error) {
	fmt.Println("FUNCTION_PARSE")
	p.IncPos()
	//fmt.Println("NEW_FUNCTION", string(p.EndLineContent()))
	p.SetupMark()
	stackOffset, info := p.stack.Len(), p.infoRecordFromMark()
	for !p.IsEndDocument() {
		p.PassSpaces()
		//fmt.Println(string(p.Char()))
		switch p.Char() {
		case ',':
			p.IncPos()
		case ')':
			p.stack.Push(&valVariable{info, name})
			p.stack.Push(&execCommand{info, execFunction, p.stack.Len() - stackOffset + 2})
			p.IncPos()
			fmt.Println("PARSE_COMPLETED")
			return
		default:
			if _, err = initCodeVal(p); err != nil {
				return
			}
		}
	}
	err = info.fatalError("Unexpected end of function exec")
	return
}

func execFunction(exec *tplExec, info *rawInfoRecord) (err error) {
	fmt.Println("FUNCTION_EXEC", exec.st.Len())
	f := reflect.ValueOf(exec.st.Pop().(*variable).value)
	fType := f.Type()
	fmt.Println("FFFFF", f, exec.st.Len())
	args := make([]reflect.Value, 0, exec.st.Len())
	for exec.st.Len() > 0 {
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		val := exec.st.Pop()
		fmt.Println("VAL", val)
		if command, check := val.(*execCommand); check {
			if err = command.method(exec, command.rawInfoRecord); err != nil {
				return
			}
		} else {
			if len(args) >= fType.NumIn() {
				err = info.positionWarning(fmt.Sprintf("Too many function arguments. Expected [%v]", fType.NumIn()))
			}
			args = append([]reflect.Value{reflect.ValueOf(val).Convert(fType.In(len(args)))}, args...)
		}

	}
	if len(args) != fType.NumIn() {
		err = info.positionWarning(fmt.Sprintf("Function agrs count isn't match - %v given, %v expected", len(args), fType.NumIn()))
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Function call fatal error (panic) :: %v", r)
		}
	}()
	if res := f.Call(args); len(res) > 0 {
		for _, val := range res {
			fmt.Println("VAL =======================", val)
			exec.st.Push(val.Interface())
		}
	}
	//err = fmt.Errorf("IIIIIIIIIIIIIIIIIIIIIIIIIIIIII")
	return
}

//////////////////////////////////////////////////////

type valFunction struct {
	*rawInfoRecord
	name string
	args []token
}

func (s *valFunction) Val() interface{} {
	return s.name
}

func (s *valFunction) IsExecutable() bool { return true }

func (s *valFunction) String() string {
	res := fmt.Sprintf("[function :: { %v }, args : { %v }]", string(s.name), s.args)
	return res
}

func (s *valFunction) posInfo() *rawInfoRecord { return s.rawInfoRecord }

/*func (s *valFunction) execObject(sto *storage, tpl *template, parent execObject) (res execObject, err error) {
	fObj := valFunctionExec{rawInfoRecord: s.rawInfoRecord}
	if f, check := sto.findVariable(s.name); check {
		if f.Kind() != reflect.Func {
			err = s.fatalError(fmt.Sprintf("Unexpected variable type [%v], [Func] expected", f.Kind()))
		} else {
			args := make([]execObject, len(s.args))
			for i, v := range s.args {
				if args[i], err = v.execObject(sto, tpl, &fObj); err != nil {
					return
				}
			}
			fObj.f, fObj.args, res = f, args, &fObj
			//res = &valFunctionExec{s.rawInfoRecord, f, args}
		}
	} else {
		err = s.fatalError(fmt.Sprintf("Function [%s] not found", s.name))
	}
	return
}*/

//////////////////////////////////////////////////////////

type valFunctionExec struct {
	*rawInfoRecord
	f    *variable
	args []executor
}

/*func (s *valFunctionExec) Data(w io.Writer) (err error) {
	if _, err = s.call(); err != nil {
		content := s.positionWarning(err.Error())
		_, err = w.Write([]byte(content.Error()))
	}
	return
}*/

func (s *valFunctionExec) IsNil() bool {
	return s.f.IsNil()
}

func (s *valFunctionExec) String() string {
	return "[function { " + s.f.key + " }]"
}

func (s *valFunctionExec) Type() reflect.Kind {
	return reflect.Func
}

/*func (s *valFunctionExec) Val() (interface{}, error) {
	if vals, err := s.Vals(); err == nil && len(vals) > 0 {
		return vals[0], err
	} else {
		return nil, err
	}
}*/

/*func (s *valFunctionExec) Vals() (res []interface{}, err error) {
	var rRes []reflect.Value
	if rRes, err = s.call(); err == nil {
		if len(rRes) > 0 {
			res = make([]interface{}, len(rRes))
			for i, v := range rRes {
				res[i] = v.Interface()
			}
		}
	}
	return
}*/

func (s *valFunctionExec) ValSingle() bool {
	return false
}

/*func (s *valFunctionExec) call() (res []reflect.Value, err error) {
	fVal := reflect.ValueOf(s.f.value)
	if fVal.Kind() != reflect.Func {
		err = fmt.Errorf("Function exec error :: variable [%v] is not a function", s.f.key)
	} else {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("Function call fatal error (panic) :: %v", r)
			}
		}()
		rVal := reflect.ValueOf(s.f.value)
		rType := rVal.Type()
		rArgs := make([]reflect.Value, 0, rVal.Type().NumIn())
		for _, v := range s.args {
			if v.ValSingle() {
				var val interface{}
				if val, err = v.Val(); err != nil {
					return
				} else {
					rArgs = append(rArgs, reflect.ValueOf(val).Convert(rType.In(len(rArgs))))
				}
			} else {
				var vals []interface{}
				if vals, err = v.Vals(); err != nil {
					return
				} else {
					for _, val := range vals {
						rArgs = append(rArgs, reflect.ValueOf(val).Convert(rType.In(len(rArgs))))
					}
				}
			}
		}
		res = fVal.Call(rArgs)
	}
	return
}*/

func (s *valFunctionExec) receiveEvent(name string, params []interface{}) bool { return false }
