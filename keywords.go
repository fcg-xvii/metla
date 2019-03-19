package metla

import (
	"fmt"
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

// core functions ///////////////////////////////////////////////////////////////

func coreLen(val interface{}) int {
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

func coreEcho(l string, vals ...interface{}) {
	fmt.Println("ECHO......", l, vals)
	for i, v := range vals {
		fmt.Println(i, v)
	}
}
