package metla

import (
	"fmt"
	"io"
)

func newTemplate(content ContentMethod, filePath string) (*template, error) {

	return nil, nil
}

type template struct {
	storage    *storage
	content    ContentMethod
	files      []string
	exec       []token
	updateMark interface{}
}

func (s *template) execute(w io.Writer, vals map[string]interface{}) ([]byte, error) {
	//fmt.Println("EXEC...")
	return nil, nil
}
