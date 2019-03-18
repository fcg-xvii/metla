package metla

import (
	_ "fmt"
	_ "io"

	_ "github.com/golang-collections/collections/stack"
)

func newValArifmetic(p *parser) (res interface{}, err error) {
	var pn []interface{}
	p.SetupMark()
	info := p.infoRecordFromMark()
	if pn, err = parseRPN(p); err == nil {
		if pn, err = simpleRPN(pn); err == nil {
			//fmt.Println("EEEEEEEEE", pn)
			if len(pn) == 1 {
				p.stack.Push(pn[0])
			} else {
				for _, v := range pn {
					p.stack.Push(v)
				}
				p.stack.Push(&execCommand{info, execArifmetic, len(pn) + 1})
			}
			//p.codeStack.Push(&arifmetic{info, pn})
		}
	}
	//fmt.Println(execRPN(pn))
	return
}

func execArifmetic(exec *tplExec) (err error) {
	//fmt.Println("EXEC_ARIFMETIC", st.Len())
	pn := make([]interface{}, 0, exec.st.Len())
	for exec.st.Len() > 0 {
		//pn = append([]interface{}{st.Pop()}, pn...)
		v := exec.st.Pop()
		if k, check := v.(*variable); check {
			pn = append([]interface{}{k.value}, pn...)
		} else {
			pn = append([]interface{}{v}, pn...)
		}
	}
	//fmt.Println("PN", pn)
	if res, err := execRPN(pn); err == nil {
		exec.st.Push(res)
	}
	return
}
