package prod

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/fcg-xvii/containers"
)

type TestSingle struct {
	Single string
	Number int
}

func (s *TestSingle) Methodd(i, j int) (int, int) {
	return i + 1, j + 1
}

type Test struct {
	Val *TestSingle
}

func (s *Test) ValStruct() *TestSingle {
	return s.Val
}

func Inc(val int) int {
	return val + 1
}

func min(left, right int) int {
	if left < right {
		return left
	}
	return right
}

func TestParser(t *testing.T) {
	var null *Test

	exVals := map[string]interface{}{
		"inc":  Inc,
		"min":  min,
		"null": null,
		"one":  1,
		"list": []string{"one", "two", "three", "four", "five"},
		"mmap": map[string]interface{}{"one": 1, "two": 2, "three": 3},
		"rMap": map[string]int{"min": 1, "max": 100000},
		"tVal": &Test{&TestSingle{"SINGLE___", 777}},
		"map": map[string]interface{}{
			"one": map[int]interface{}{
				100: 200,
			},
		},
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
			false,
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
