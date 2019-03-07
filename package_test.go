package metla

import (
	"bytes"
	_ "errors"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func printMethod(s interface{}) {
	log.Println(s)
}

func printTwink(x, y interface{}) {
	log.Println("XY", x, y)
}

func inc(v int64) int64 {
	return v + 1
}

func check(path string, marker interface{}) (res UpdateState) {
	if info, err := os.Stat(path); err == nil {
		if marker != nil {
			if mTime := marker.(time.Time); mTime.Equal(info.ModTime()) {
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

func content(path string, marker interface{}) (res []byte, newMarker interface{}, state UpdateState) {
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
			if markerTime, check := marker.(time.Time); check && markerTime.Equal(info.ModTime()) {
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
		"colonel": "Hello, WORLD!",
		"print":   printMethod,
		"twink":   printTwink,
		"inc":     inc,
	}

	if err := root.Content("source_script", &buf, data); err != nil {
		log.Println("ERR", err)
	} else {
		log.Println("OK")
		buf.WriteTo(os.Stdout)
	}
}
