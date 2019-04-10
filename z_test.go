package metla

import (
	"bytes"
	_ "errors"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

type Child struct {
	One int64
}

func (s *Child) Min10(left, right int) bool {
	log.Println("MIN10", left, right)
	return left < right
}

type Test struct {
	ChildObj *Child
}

func (s *Test) GetChild() *Child {
	return s.ChildObj
}

func printMethod(s interface{}) {
	log.Println(s)
}

func printTwink(x, y interface{}) {
	log.Println("XY", x, y)
}

func inc(v int64) int64 {
	return v + 1
}

func incTwo(l int64, r int64) (int64, int64) {
	return l + 1, r + 1
}

func cooler(one, two, three int) {
	log.Println(one, two, three)
}

func check(path string, marker *time.Time) (res UpdateState) {
	if info, err := os.Stat(path); err == nil {
		if marker != nil {
			if (*marker).Equal(info.ModTime()) {
				res = UpdateNotNeeded
			} else {
				res = UpdateNeeded
			}
		} else {
			res = UpdateNeeded
		}
	} else {
		res = ResourceNotFound
	}
	return
}

func content(path string, marker *time.Time) (res []byte, newMarker time.Time, state UpdateState) {
	readContent := func() {
		var err error
		if res, err = ioutil.ReadFile(path); err != nil {
			state = ResourceNotFound
		}
	}
	if info, err := os.Stat(path); err == nil {
		newMarker = info.ModTime()
		if marker == nil {
			state = UpdateNeeded
			readContent()
		} else {
			if marker != nil && (*marker).Equal(info.ModTime()) {
				state = UpdateNotNeeded
			} else {
				state = UpdateNeeded
				readContent()
			}
		}
	} else {
		state = ResourceNotFound
	}
	return
}

func TestParser(t *testing.T) {
	root := New(check, content)

	var buf bytes.Buffer

	data := map[string]interface{}{
		"one":     1,
		"three":   3,
		"colonel": "Hello, WORLD!",
		"print":   printMethod,
		"twink":   printTwink,
		"inc":     inc,
		"incTwo":  incTwo,
		"sli":     []byte{1, 2, 3, 4},
		"cooler":  cooler,
		"tr":      true,
		"cli":     map[string]string{"one": "over one"},
		"tst":     &Test{&Child{One: 5}},
		"req":     httptest.NewRequest("GET", "/index.html", nil),
	}

	for {

		if modified, err := root.Content("z_script", &buf, data); err != nil {
			log.Println("ERR", err)
		} else {
			log.Println("OK", modified)
			buf.WriteTo(os.Stdout)
		}

		time.Sleep(time.Second * 3)
	}
}
