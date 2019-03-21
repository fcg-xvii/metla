package metla

import (
	"fmt"
	"io"
	"reflect"
)

type keywordConstructor func(*parser) (interface{}, error)

func init() {
	keywords["echo"] = keywordEcho
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
		"len": coreLen,
	}
)

func getKeywordConstructor(name string) (result keywordConstructor, check bool) {
	result, check = keywords[name]
	return
}

// keywords ////////////////////////////////////////////////////////////////////

func keywordEcho(p *parser) (res interface{}, err error) {
	info := p.infoRecordFromPos()
	res = &execCommand{info, execEcho, 0}
	p.stack.Push(res)
	p.pushSplitter()
	for !p.IsEndDocument() {
		p.PassSpaces()
		if _, err = initCodeVal(p); err != nil {
			return
		}
		switch p.Char() {
		case '\n', ';':
			return
		case ',':
			p.pushSplitter()
			p.IncPos()
		default:
			p.positionError(fmt.Sprintf("Unexpected symbol '%c'", p.Char()))
		}
	}
	return nil, info.fatalError("Unexpected end of document")
}

func execEcho(exec *tplExec, info *rawInfoRecord) (err error) {
	for exec.st.Len() > 0 {
		fmt.Println(exec.st.Peek())
		exec.w.Write([]byte(fmt.Sprint(exec.st.Pop())))
		if exec.st.Len() >= 1 {
			exec.w.Write([]byte{',', ' '})
		}

	}
	return
}

// core functions ///////////////////////////////////////////////////////////////

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
