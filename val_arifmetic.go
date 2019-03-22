package metla

import _ "fmt"

func newValArifmetic(p *parser) (res interface{}, err error) {
	var pn []interface{}
	if pn, err = parseRPN(p); err == nil {
		if pn, err = simpleRPN(pn); err == nil {
			if len(pn) == 1 {
				p.stack.Push(pn[0])
			} else {
				//p.stack.Push(&execCommand{p.infoRecordFromPos(), execArifmetic, 0})
				/*for _, v := range pn {
					p.stack.Push(v)
				}*/
				for i := len(pn) - 1; i >= 0; i-- {
					p.stack.Push(pn[i])
				}
				//p.stack.Push(initMarker())
			}
		}
	}
	//fmt.Println("PNNNNNNN!!!!!!!!!!!", pn)
	return
}
