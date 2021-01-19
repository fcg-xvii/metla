package metla

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/fcg-xvii/mjs"
)

func New(modifiedCallback func(string) int64, contentCallback func(string) ([]byte, error)) *Metla {
	m := &Metla{
		contentCallback: contentCallback,
	}
	m.js = mjs.New(modifiedCallback, m.content)
	return m
}

type Metla struct {
	js              *mjs.Mjs
	contentCallback func(string) ([]byte, error)
}

func (s *Metla) content(name string) (content []byte, err error) {
	if content, err = s.contentCallback(name); err == nil {
		if filepath.Ext(name) != ".script" {
			content, err = parseBytes(content, name)
		}
	}
	return
}

func (s *Metla) Exec(name string, data map[string]interface{}, w io.Writer) (modified int64, err error) {
	if data == nil {
		data = make(map[string]interface{})
	}
	data["flush"] = func(vals ...interface{}) {
		for _, v := range vals {
			w.Write([]byte(fmt.Sprint(v)))
		}
	}
	return s.js.Exec(name, data)
}
