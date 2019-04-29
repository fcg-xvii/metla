package prod

import (
	"fmt"
	"reflect"
)

func parseFunctionArgs(p *parser, fPos *position) (list []interface{}, err *parseError) {
	stackLen := p.stack.Len()
	for !p.IsEndDocument() {
		switch {
		case p.isEndCode():
			return nil, p.initParseError(p.Line(), p.LinePos(), "Uncosed function")
		case p.Char() == ',' || p.Char() == ')':
			if stackLen != p.stack.Len()-1 {
				return nil, fPos.parseError("Expected ',' or ')' character")
			}
			list = append(list, p.stack.Pop())
			if p.Char() == ')' {
				return
			} else {
				p.IncPos()
			}
		default:
			if err = p.initCodeVal(); err != nil {
				return
			}
			//list = append(list, p.stack.Pop())
			//valAccepted = false
		}
	}
	return nil, (*fPos).parseError("Unexpected end of document")
}

func execArgsPrepare(exec *tplExec, fType reflect.Type, args []interface{}) (rArgs []reflect.Value, err *execError) {
	fCount, stackLen := fType.NumIn(), exec.stack.Len()
	rArgs = make([]reflect.Value, 0, fCount)

	argsLenCheck := func(v interface{}) {
		if fCount < len(rArgs) {
			err = v.(coordinator).execError("Args count more than needed")
		}
	}

	convert := func(v interface{}, rType reflect.Type) (val reflect.Value) {
		fmt.Println("convert <<< ", v)
		lVal := reflect.ValueOf(v)
		lType := lVal.Type()
		fmt.Println("PPPP", lType, rType)
		if lType != rType {
			if lType.ConvertibleTo(rType) {
				val = lVal.Convert(rType)
			} else {
				err = v.(coordinator).execError(fmt.Sprintf("Wrong function arg %v", len(rArgs)))
			}
		} else {
			return lVal
		}
		return
	}

	for _, v := range args {
		if argsLenCheck(v); err != nil {
			return
		}
		fmt.Println("222", v)
		switch v.(type) {
		case getter:
			if rArg := convert(v.(getter).get(exec), fType.In(len(rArgs))); err != nil {
				return
			} else {
				rArgs = append(rArgs, rArg)
			}
		case executer:
			obj := v.(executer)
			if !functionCheckExecType(obj.execType()) {
				return nil, obj.execError(fmt.Sprintf("Type %v isn't valid function argument", obj.execType()))
			}
			if err = obj.exec(exec); err != nil {
				return
			}
			for exec.stack.Len() > stackLen {
				if argsLenCheck(v); err != nil {
					return
				}
				if rArg := convert(exec.stack.Pop().(getter).get(exec), fType.In(len(rArgs))); err != nil {
					return
				} else {
					rArgs = append(rArgs, rArg)
				}
			}
		}
	}
	return
}

func newFunction(p *parser) (err *parseError) {
	fmt.Println("NEW_FUNCTION")
	f, check, returnCall := function{position: p.posObject()}, false, p.resetFlags()
	if f.nameVar, check = p.stack.Pop().(*iName); !check {
		return f.parseError("Function parse error :: expected variable in prefix")
	}
	//return p.initParseError(10, 10, "Function error test")
	p.IncPos()
	if f.args, err = parseFunctionArgs(p, &f.position); err != nil {
		return err
	}
	fmt.Println(">>>>", f.args)
	p.stack.Push(&f)
	p.IncPos()
	returnCall()
	return
}

/////////////////////////////////////////////

type function struct {
	position
	nameVar getter
	args    []interface{}
}

func (s *function) execType() execType {
	return execFunction
}

func functionCheckExecType(eType execType) bool {
	return eType == execFunction
}

func (s *function) exec(exec *tplExec) (err *execError) {
	//return s.execError("TTTTT")
	rName := reflect.ValueOf(s.nameVar.get(exec))
	if rName.Kind() != reflect.Func {
		return s.execError("Variable is not a function")
	}
	//fRArgs, fRType := make([]reflect.Value, 0, len(s.args)), rName.Type()
	var args []reflect.Value
	if args, err = execArgsPrepare(exec, rName.Type(), s.args); err != nil {
		return
	}
	fmt.Println(len(args), rName.Type().NumIn())
	for _, v := range rName.Call(args) {
		exec.stack.Push(static{s.position, v.Interface()})
	}
	fmt.Println("F_STACK", exec.stack)
	return
}

/////////////////////////////////////

func newMethod(p *parser) (err *parseError) {
	fmt.Println("NEW_Method")
	f, check, returnCall := function{position: p.posObject()}, false, p.resetFlags()
	if f.nameVar, check = p.stack.Pop().(*iName); !check {
		return f.parseError("Function parse error :: expected variable in prefix")
	}
	//return p.initParseError(10, 10, "Function error test")
	p.IncPos()
	if f.args, err = parseFunctionArgs(p, &f.position); err != nil {
		return err
	}
	fmt.Println(">>>>", f.args)
	p.stack.Push(&f)
	p.IncPos()
	returnCall()
	return
}

type method struct {
	position
	nameVar getter
	args    []interface{}
}
