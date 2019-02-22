package metla

import (
	_ "bytes"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func check(path string) (res bool) {
	if info, err := os.Stat(path); err == nil {
		return info.IsDir()
	}
	return
}

func content(path string, marker interface{}) (res []byte, check bool, newMarker interface{}, err error) {
	var info os.FileInfo
	if info, err = os.Stat(path); err == nil {
		if marker != nil {
			var markerTime time.Time
			if markerTime, check = marker.(time.Time); check {
				check = !markerTime.Equal(info.ModTime())
			} else {
				err = errors.New("Unexpected marker type. [time.Time] expected")
			}
		} else {
			check = true
		}
		if check {
			newMarker = info.ModTime()
			res, err = ioutil.ReadFile(path)
		}
	}
	return
}

func TestParser(t *testing.T) {
	if content, err := ioutil.ReadFile("source_script"); err != nil {
		t.Error(err)
	} else {
		log.Println(string(content))
		pObj := newParser(content)
		log.Println(pObj.parseCallBack(func(p *parser) bool {
			return !p.IsEndLine()
		}))
		//log.Println(parser.execList)
		//log.Println("[" + string(parser.availableData()) + "]")
	}

	//tpl, err := newTemplate("source_script")
	//log.Println(tpl, err)

	/*metla := New(check, content)
	data := map[string]interface{}{
		"one": 1,
	}
	var content []byte
	buf := bytes.NewBuffer(content)
	if err := metla.Content("source_script", buf, data); err != nil {
		t.Error(err)
	} else {
		log.Println("==========================================")
		log.Println(content)
	}*/
}
