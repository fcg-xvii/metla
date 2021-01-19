package metla

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var contentPath = "content"

func filePath(name string) string {
	return fmt.Sprintf("%v/%v", contentPath, name)
}

func modified(name string) (res int64) {
	if info, err := os.Stat(filePath(name)); err == nil {
		res = info.ModTime().Unix()
	}
	return
}

func content(name string) ([]byte, error) {
	return ioutil.ReadFile(filePath(name))
}

func TestParser(t *testing.T) {
	m := New(modified, content)
	log.Println(m)

	var b bytes.Buffer

	params := map[string]interface{}{
		"title": "Heya ))",
	}
	log.Println(m.Exec("z_source.html", params, &b))

	log.Println(b.String())
}
