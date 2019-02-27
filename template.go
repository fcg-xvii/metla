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
			//fmt.Println("TEMPLATE_EXECUTE")
			sto := newStorage(vals)
			err = s.exec(w, sto)
		}
	}
	s.locker.Unlock()
	fmt.Println(s.tokenList)
	return
}

func (s *template) parse(src []byte) error {
	parser := newParser(src, s, s.root)
	return parser.parseDocument()
}

func (s *template) exec(w io.Writer, sto *storage) (err error) {
	var execObj execObject
	for _, v := range s.tokenList {
		if execObj, err = v.execObject(sto, s); err != nil {
			fmt.Println("Exec error ::", err)
			return
		}
		if err = execObj.Data(w); err != nil {
			fmt.Println("ExecObj error ::", err)
			return
		}
		/*if err = v.Data(w, sto); err != nil {
			return
		}*/
	}
	return
}
