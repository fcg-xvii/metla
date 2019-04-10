package metla

import (
	"fmt"
	"io"
	"reflect"
)

type keywordConstructor func(*parser) (interface{}, error)

func init() {
	keywords["echo"] = keywordEcho
	keywords["echoln"] = keywordEcholn
	keywords["print"] = keywordPrint
	keywords["println"] = keywordPrintln
}

var (
	keywords = map[string]keywordConstructor{
		"nil": func(p *parser) (interface{}, error) {
			p.stack.Push(nil)
			return nil, nil
		}, "true": func(p *parser) (interface{}, error) {
			p.stack.Push(true)
			return true, nil
		}, "false": func(p *parser) (interface{}, error) {
			p.stack.Push(false)
			return false, nil
		},
	}
	functions = map[string]interface{}{
		"len":     coreLen,
		"defined": coreDefined,
	}
)

func getKeywordConstructor(name string) (result keywordConstructor, check bool) {
	result, check = keywords[name]
	return
}

// keywords ////////////////////////////////////////////////////////////////////

func parseKeywordArgs(p *parser, info *rawInfoRecord) (err error) {
	for !p.IsEndDocument() {
		p.PassSpaces()
		if _, err = initCodeVal(p); err != nil {
			return
		}
		switch p.Char() {
		case '\n', ';':
			return
		case ',':
			//p.pushSplitter()
			p.IncPos()
		default:
			p.positionError(fmt.Sprintf("Unexpected symbol '%c'", p.Char()))
		}
	}
	return info.fatalError("Unexpected end of document")
}

func keywordPrint(p *parser) (res interface{}, err error) {
	info := p.infoRecordFromPos()
	res = &execCommand{info, execKeywordPrint, "print", nil}
	p.stack.Push(res)
	p.pushSplitter()
	return res, parseKeywordArgs(p, info)
}

func keywordPrintln(p *parser) (res interface{}, err error) {
	info := p.infoRecordFromPos()
	res = &execCommand{info, execKeywordPrintln, "println", nil}
	p.stack.Push(res)
	p.pushSplitter()
	return res, parseKeywordArgs(p, info)
}

func keywordEcho(p *parser) (res interface{}, err error) {
	info := p.infoRecordFromPos()
	res = &execCommand{info, execEcho, "echo", nil}
	p.stack.Push(res)
	p.pushSplitter()
	return res, parseKeywordArgs(p, info)
}

func keywordEcholn(p *parser) (res interface{}, err error) {
	info := p.infoRecordFromPos()
	res = &execCommand{info, execEcholn, "echoln", nil}
	p.stack.Push(res)
	p.pushSplitter()
	return res, parseKeywordArgs(p, info)
}

func execEcho(exec *tplExec, info *rawInfoRecord) (err error) {
	for exec.st.Len() > 0 {
		//fmt.Println(exec.st.Peek())
		exec.w.Write([]byte(fmt.Sprint(exec.st.Pop())))
		if exec.st.Len() >= 1 {
			exec.w.Write([]byte{',', ' '})
		}
	}
	return
}

func execEcholn(exec *tplExec, info *rawInfoRecord) (err error) {
	for exec.st.Len() > 0 {
		//fmt.Println(exec.st.Peek())
		exec.w.Write([]byte(fmt.Sprint(exec.st.Pop())))
		if exec.st.Len() >= 1 {
			exec.w.Write([]byte{',', ' '})
		}
	}
	exec.w.Write([]byte{'\n'})
	return
}

func execKeywordPrint(exec *tplExec, info *rawInfoRecord) (err error) {
	for exec.st.Len() > 0 {
		exec.w.Write([]byte(fmt.Sprint(exec.st.Pop())))
	}
	return
}

func execKeywordPrintln(exec *tplExec, info *rawInfoRecord) (err error) {
	for exec.st.Len() > 0 {
		exec.w.Write([]byte(fmt.Sprint(exec.st.Pop())))
	}
	exec.w.Write([]byte{'\n'})
	return
}

// core functions ///////////////////////////////////////////////////////////////

func coreDefined(w io.Writer, val interface{}) (res bool) {
	if sVar, check := val.(*variable); check {
		res = sVar.stored
	}
	return
}

func coreLen(w io.Writer, val interface{}) int {
	if sVar, check := val.(*variable); check {
		val = sVar.value
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
		return v.Len()
	default:
		return 0
	}
}
