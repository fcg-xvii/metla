package metla

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestParser(t *testing.T) {
	if content, err := ioutil.ReadFile("source_script"); err != nil {
		t.Error(err)
	} else {
		log.Println(string(content))
		parser := newParser(content)
		log.Println(parser.parseDocument())
		log.Println(parser.execList)
		//log.Println("[" + string(parser.availableData()) + "]")
	}
}
