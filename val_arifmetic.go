package metla

import (
	"fmt"
	"io"

	"github.com/golang-collections/collections/stack"
)

func newValArifmetic(p *parser) (err error) {
	var pn []interface{}
	p.SetupMark()
	info := p.infoRecordFromMark()
	if pn, err = parseRPN(p); err == nil {
		if pn, err = simpleRPN(pn); err == nil {
			fmt.Println("EEEEEEEEE", pn)
			if len(pn) == 1 {
				p.stack.Push(pn[0])
			} else {
				for _, v := range pn {
					p.stack.Push(v)
				}
				p.stack.Push(&execCommand{info, execArifmetic})
			}
			//p.codeStack.Push(&arifmetic{info, pn})
		}
	}
	//fmt.Println(execRPN(pn))
	return
}

func execArifmetic(com []interface{}, st *stack.Stack, sto *storage, w io.Writer) (newCom []interface{}, err error) {
	fmt.Println("EXEC_ARIFMETIC")
	pn := make([]interface{}, 0, st.Len())
	for st.Len() > 0 {
		pn = append(pn, st.Pop())
	}
	fmt.Println("PN", pn)
	return
	/*if _, err = w.Write(st.Pop().([]byte)); err == nil {
		newCom = com[1:]
	}
	return*/
}
