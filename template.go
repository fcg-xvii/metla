package metla

import (
	"fmt"
	"io"
)

func newTemplate(filePath string) (*template, error) {

}

type template struct {
	files []string
	exec  []token
}

func (s *template) exec(w io.Writer, vals map[string]interface{}) []byte {
	fmt.Println("EXEC...")
}
