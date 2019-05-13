package prod

import (
	_ "reflect"
)

/*func rpnValueNumber(val interface{}, pos position, exec *tplExec) (float64, *execError) {
	var rVal reflect.Value
	switch val.(type) {
	case getter:
		rVal = reflect.ValueOf(val.(getter).get(exec))
	case executer:
		stackLen := exec.stack.Len()
		if err := val.(executer).exec(tplExec); err == nil {

		} else {
			return 0, err
		}
	}
}*/
