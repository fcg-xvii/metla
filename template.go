package metla

import (
	"fmt"
	"io"
	"sync"
)

func newTemplate(root *Metla, objPath string) *template {
	return &template{
		root:    root,
		objPath: objPath,
		locker:  new(sync.Mutex),
	}
}

type template struct {
	root       *Metla
	objPath    string
	tokenList  []token
	updateMark interface{}
	locker     *sync.Mutex
	//lastRequest time.Time

}

func (s *template) execute(w io.Writer, vals map[string]interface{}) (err error) {
	fmt.Println("TEMPLATE_EXECUTE...")
	var (
		check  bool
		source []byte
	)
	s.locker.Lock()
	if source, check, s.updateMark, err = s.root.content(s.objPath, s.updateMark); err == nil {
		if check {
			fmt.Println("TEMPLATE_UPDATE")
			err = s.parse(source)
		}
		if err == nil {
			fmt.Println("TEMPLATE_EXECUTE")
			err = s.exec(w, vals)
		}
	}
	s.locker.Unlock()
	return
}

func (s *template) parse(src []byte) error {
	parser := newParser(src, s, s.root)
	return parser.parseDocument()
}

func (s *template) exec(w io.Writer, vals map[string]interface{}) error {
	// todo template exec tokens
	return nil
}
