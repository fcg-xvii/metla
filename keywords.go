package metla

import (
	"fmt"
	"io"
	"reflect"
)

type keywordConstructor func(*parser) (interface{}, error)

var (
	keywords  = make(map[string]keywordConstructor)
	functions = map[string]interface{}{
		"len":  coreLen,
		"echo": coreEcho,
	}
)

func getKeywordConstructor(name string) (result keywordConstructor, check bool) {
	result, check = keywords[name]
	return
}

// keywords ////////////////////////////////////////////////////////////////////

func keywordEcho(p *parser) (interface{}, error) {
	info := p.infoRecordFromPos()
	for !p.IsEndDocument() {
		p.PassSpaces()
		
		switch p.Char() {
			case 
		}
	}
	return nil, info.fatalError("Unexpected end of document")
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
