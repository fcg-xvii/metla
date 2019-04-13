package prod

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/fcg-xvii/containers"
)

func TestParser(t *testing.T) {
	exVals := map[string]interface{}{
		"one": 1,
	}

	src, _ := ioutil.ReadFile("z_content")
	//log.Println(string(src))
	parser := initParser("z_content", src)
	if err := parser.parseDocument(); err == nil {
		log.Println(parser.execList)
		var buf bytes.Buffer
		ex := &tplExec{
			"z_content",
			parser.execList,
			&buf,
			parser.store.execStorage(exVals),
			new(containers.Stack),
		}
		if err := ex.exec(); err != nil {
			log.Println(err)
		} else {
			log.Println(ex.sto.values)
			log.Println("======================")
			buf.WriteTo(os.Stdout)
			log.Println("======================")
		}
	} else {
		log.Println(err)
	}
}
