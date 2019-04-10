package prod

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	src, _ := ioutil.ReadFile("z_content")
	//log.Println(string(src))
	parser := initParser("z_conntent", src)
	if err := parser.parseDocument(); err == nil {
		var buf bytes.Buffer
		ex := &tplExec{
			parser.execList,
			&buf,
			newStorage(map[string]interface{}{}),
		}
		if err := ex.exec(); err != nil {
			log.Println(err)
		} else {
			log.Println("======================")
			buf.WriteTo(os.Stdout)
			log.Println("======================")
		}
	} else {
		log.Println(err)
	}
}
