//规则定义
package cnclassify

import (
	"strings"
)

type Rule struct {
	Point      int      //当前解析词位置
	Expression []string //拆分过的表达式
}

//var rules map[string][]Rule //规则集合

//解释规则，递归解析，支持括号嵌套
func (r *Rule) Compute(text string) bool {
	if len(r.Expression) == 0 {
		return false
	}
	result := true
	op := "+"
	for ; r.Point < len(r.Expression); r.Point++ {
		tmp := r.Expression[r.Point]
		switch tmp {
		case "+", "|", "^":
			op = tmp
			continue
		}
		tmpresult := false
		if tmp == "(" {
			r.Point++
			tmpresult = r.Compute(text)
		} else if tmp == ")" {
			return result
		} else {
			tmpresult = strings.Index(text, tmp) > -1
		}
		switch op {
		case "+":
			result = result && tmpresult
		case "|":
			result = result || tmpresult
		case "^":
			result = result && (!tmpresult)
		}

	}
	return result
}
