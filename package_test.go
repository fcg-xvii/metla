package metla

import (
	"log"
	"os"
	"io/ioutil"
	"testing"
)

type CheckMetod func(string) bool
type ContentMethod func(string) ([]byte, interface{}, error)
type UpdateMethod func(string, interface{}) ([]byte, bool, interface{}, error)

func check(path string) (res bool) {
	if info, err := os.Stat(path); err == nil {
		return info.IsDir()
	}
	return
}

func content(path string) (res []byte, marker interface{}, err error) {
	var info os.FileInfo
	if info, err = os.Stat(path); err == nil {
		marker = info.ModTime()
		res, err = ioutil.ReadAll()
	}
	return
}

func update(path string, interface{}) (res []byte, check bool, marker interface{}, err error) {
	
}

func TestParser(t *testing.T) {
	/*if content, err := ioutil.ReadFile("source_script"); err != nil {
		t.Error(err)
	} else {
		log.Println(string(content))
		parser := newParser(content)
		log.Println(parser.parseDocument())
		log.Println(parser.execList)
		//log.Println("[" + string(parser.availableData()) + "]")
	}*/

	tpl := newTemplate("source_script")
}
