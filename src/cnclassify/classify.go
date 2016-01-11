//中文分类
package cnclassify

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//
type Classify struct {
	Name  string //类别名称
	Rules map[string][]Rule
}

//从文件目录中加载配置，初始化的时候调用
func (c *Classify) LoadRulesByDir(filedir /*文件目录*/ string) {
	c.Rules = make(map[string][]Rule)
	filepath.Walk(filedir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		} else if !strings.HasSuffix(info.Name(), ".rul") {
			return nil
		}
		//处理文件
		fi, _ := os.Open(path)
		bs, _ := ioutil.ReadAll(fi)
		fi.Close()
		arr := []Rule{}
		tmp := strings.Split(string(bs), "\r\n")
		//
		for _, v := range tmp {
			scan := NewScanner(strings.NewReader(v))
			scan.Split(MyStopWord)
			exp := []string{}
			for scan.Scan() {
				if scan.Text() != "" {
					exp = append(exp, scan.Text())
				}
				if scan.Stopchar() > 0 {
					exp = append(exp, string(scan.Stopchar()))
				}
			}
			arr = append(arr, Rule{Expression: exp})
		}
		//log.Println(path, info.Name())
		rulename := path[len(filedir)+1 : strings.LastIndex(path, ".")]
		rulename = strings.Replace(rulename, "\\", "/", -1)
		//log.Println(rulename)
		c.Rules[rulename] = arr
		return nil
	})
}

//从字符，转换为规则
func (c *Classify) LoadRulesByString(name, rule string) {
	c.Rules = make(map[string][]Rule)
	scan := NewScanner(strings.NewReader(rule))
	scan.Split(MyStopWord)
	exp := []string{}
	for scan.Scan() {
		if scan.Text() != "" {
			exp = append(exp, scan.Text())
		}
		if scan.Stopchar() > 0 {
			exp = append(exp, string(scan.Stopchar()))
		}
	}
	c.Rules[name] = []Rule{Rule{Expression: exp}}
}

//归类
func (c *Classify) Classification(text string) []string {
	result := []string{} //记录满足条件的规则名称
	for k, v := range c.Rules {
		//每个规则文件中有多个规则，只要有一个匹配成功，则认为满足此分类
		for _, r := range v {
			r.Point = 0 //重置指针位置
			if r.Compute(text) {
				result = append(result, k)
				break
			}
		}
	}
	return result
}
