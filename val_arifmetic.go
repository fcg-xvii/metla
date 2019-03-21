package metla

import "fmt"

func newValArifmetic(p *parser) (res interface{}, err error) {
	var pn []interface{}
	if pn, err = parseRPN(p); err == nil {
		if pn, err = simpleRPN(pn); err == nil {
			if len(pn) == 1 {
				p.stack.Push(pn[0])
			} else {
				p.stack.Push(&execCommand{p.infoRecordFromPos(), execArifmetic, 0})
				for i := len(pn) - 1; i >= 0; i-- {
					p.stack.Push(pn[i])
				}
			}
		}
	}
	//fmt.Println("PNNNNNNN!!!!!!!!!!!", pn)
	return
}

func execArifmetic(exec *tplExec, info *rawInfoRecord) (err error) {
	pn := make([]interface{}, 0, exec.st.Len())
	fmt.Println("EXXX")
	for exec.st.Len() > 0 {
		v := exec.st.Pop()
		if k, check := v.(*variable); check {
			pn = append([]interface{}{k.value}, pn...)
		} else {
			pn = append([]interface{}{v}, pn...)
		}
	}
	var res interface{}
	if res, err = execRPN(pn); err == nil {
		exec.st.Push(res)
		//mt.Println("RESSSSS", res)
	} else {
		err = info.positionWarning(err.Error())
	}
	return
}
