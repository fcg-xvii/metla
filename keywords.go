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
	keywords  = make(map[string]keywordConstructor)
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
	argsCount := 0
	for !p.IsEndDocument() {
		p.PassSpaces()
		if _, err = initCodeVal(p); err != nil {
			return
		}
		argsCount++
		switch p.Char() {
		case '\n', ';':
			res = &execCommand{info, execEcho, 0}
			p.stack.Push(argsCount)
			p.stack.Push(res)
			return
		case ',':
			p.IncPos()
		default:
			p.positionError(fmt.Sprintf("Unexpected symbol '%c'", p.Char()))
		}
	}
	return nil, info.fatalError("Unexpected end of document")
}

func execEcho(exec *tplExec, info *rawInfoRecord) (err error) {
	//err = fmt.Errorf("OKKK")
	vfvdfvdfs
	return
}

// core functions ///////////////////////////////////////////////////////////////

func coreLen(w io.Writer, val interface{}) int {
	if sVar, check := val.(*variable); check {
		val = sVar.value
	}
	v := reflect.ValueOf(val)
	//fmt.Println("VALLLLL", val, v.Kind())
	switch v.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
		return v.Len()
	default:
		return 0
	}
}
