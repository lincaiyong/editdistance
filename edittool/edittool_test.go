package edittool

import (
	"fmt"
	"testing"
)

func TestDiffByLine(t *testing.T) {
	a := `- 文档处理：处理格式略有差异的文本
	- 代码重构工具：容忍空白字符差异的替换
	- 模板引擎：灵活的文本替换
	
函数会：
	1. 优先寻找完全匹配（编辑距离为0）
	2. 如果没有完全匹配，选择编辑距离最小的匹配进行替换
	3. 通过规范化忽略空白字符差异
`
	b := `- 文档处理：处理格式略有差异的文本
- 代码重构工具：容忍空白字符差异的替换
- 模板引擎：灵活的文本替换

函数会：
1. x优先寻找完全匹配（编辑距离为0）
2. 如果没有完全匹配，  选择编辑距离最小的匹配进行替换
3. 通过规范化忽略空白字符差异
`
	ret := DiffByLine(a, b, Options{IgnoreWhitespace: true})
	fmt.Println(ret)
}

func TestDiffByLine2(t *testing.T) {
	a := `- 文档处理：处理格式略有差异的文本
	- 代码重构工具：容忍空白字符差异的替换
	- 模板引擎：灵活的文本替换
	
函数会：
	1. 优先寻找完全匹配（编辑距离为0）
	2. 如果没有完全匹配，选择编辑距离最小的匹配进行替换
	3. 通过规范化忽略空白字符差异
`
	b := `- 文档处理：处理格式略有差异的文本
- 代码重构工具：容忍空白字符差异的替换
- 模板引擎：灵活的文本替换

函数会：
1. x优先寻找完全匹配（编辑距离为0）
2. 如果没有完全匹配，  选择编辑距离最小的匹配进行替换
3. 通过规范化忽略空白字符差异
`
	ret := DiffByLine(a, b, Options{IgnoreWhitespace: true})
	fmt.Println(ret)

	ret = DiffByLine(a, b, Options{IgnoreWhitespace: false, ContextLines: -1, WithLineNo: true})
	fmt.Println(ret)
}
