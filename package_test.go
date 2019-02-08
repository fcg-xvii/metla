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
		log.Println(parser.parseToEndLine())
		a := []byte{1, 2, 3, 4, 5}
		log.Println(a[:1])
		//log.Println("[" + string(parser.availableData()) + "]")
	}
}
