package edittool

import (
	"fmt"
	"strings"
	"testing"
)

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

	// 测试1: 忽略空白字符
	ret := DiffByLine(a, b, Options{IgnoreWhitespace: true})
	fmt.Println("=== 忽略空白字符 ===")
	fmt.Println(ret)
	// 预期: 只显示内容真正不同的行
	if !strings.Contains(ret, "x优先寻找") {
		t.Error("预期应该包含 'x优先寻找' 的修改")
	}

	// 测试2: 不忽略空白字符，显示所有行和行号
	ret = DiffByLine(a, b, Options{IgnoreWhitespace: false, ContextLines: -1, WithLineNo: true})
	fmt.Println("\n=== 不忽略空白字符，显示所有行 ===")
	fmt.Println(ret)
	// 预期: 显示所有差异，包括空白字符差异
	if !strings.Contains(ret, "-|") || !strings.Contains(ret, "+|") {
		t.Error("预期应该包含删除和添加标记")
	}
}

// 测试用例1: 简单的增删改
func TestDiffByLine_SimpleChanges(t *testing.T) {
	a := `line 1
line 2
line 3
line 4
line 5`
	b := `line 1
line 2 modified
line 4
line 5
line 6`

	fmt.Println("=== 测试1: 简单的增删改 (显示所有行) ===")
	ret := DiffByLine(a, b, Options{IgnoreWhitespace: false, ContextLines: -1, WithLineNo: true})
	fmt.Println(ret)

	// 预期检查
	expected := []string{
		" |0001|", // line 1 保持
		"-|0002|", // line 2 删除
		"+|0002|", // line 2 modified 添加
		"-|0003|", // line 3 删除
		" |0004|", // line 4 保持
		" |0005|", // line 5 保持
		"+|0006|", // line 6 添加
	}

	for _, exp := range expected {
		if !strings.Contains(ret, exp) {
			t.Errorf("预期包含 '%s'", exp)
		}
	}

	fmt.Println("\n=== 测试1: 简单的增删改 (上下文2行) ===")
	ret = DiffByLine(a, b, Options{IgnoreWhitespace: false, ContextLines: 2, WithLineNo: true})
	fmt.Println(ret)

	// 预期: 应该显示所有行（因为修改分散，上下文2行会覆盖全部）
	if !strings.Contains(ret, "line 1") || !strings.Contains(ret, "line 6") {
		t.Error("预期应该包含第一行和最后一行")
	}
}

// 测试用例2: 纯插入
func TestDiffByLine_OnlyInsert(t *testing.T) {
	a := `func main() {
	fmt.Println("Hello")
}`
	b := `func main() {
	fmt.Println("Hello")
	fmt.Println("World")
	fmt.Println("!")
}`

	fmt.Println("=== 测试2: 纯插入 ===")
	ret := DiffByLine(a, b, Options{IgnoreWhitespace: false, ContextLines: -1, WithLineNo: true})
	fmt.Println(ret)

	// 预期检查
	insertCount := strings.Count(ret, "+|")
	if insertCount != 2 {
		t.Errorf("预期有2个插入操作，实际有 %d 个", insertCount)
	}

	deleteCount := strings.Count(ret, "-|")
	if deleteCount != 0 {
		t.Errorf("预期没有删除操作，实际有 %d 个", deleteCount)
	}
}

// 测试用例3: 纯删除
func TestDiffByLine_OnlyDelete(t *testing.T) {
	a := `func calculate() {
	x := 1
	y := 2
	z := 3
	return x + y + z
}`
	b := `func calculate() {
	x := 1
	return x
}`

	fmt.Println("=== 测试3: 纯删除 ===")
	ret := DiffByLine(a, b, Options{IgnoreWhitespace: false, ContextLines: -1, WithLineNo: true})
	fmt.Println(ret)

	// 预期检查
	deleteCount := strings.Count(ret, "-|")
	if deleteCount != 3 {
		t.Errorf("预期有3个删除操作，实际有 %d 个", deleteCount)
	}

	insertCount := strings.Count(ret, "+|")
	if insertCount != 0 {
		t.Errorf("预期没有插入操作，实际有 %d 个", insertCount)
	}
}

// 测试用例4: 空白字符差异
func TestDiffByLine_WhitespaceIgnore(t *testing.T) {
	a := `	if (condition) {
		doSomething();
	}`
	b := `if (condition) {
    doSomething();
}`

	fmt.Println("=== 测试4: 忽略空白字符 ===")
	ret := DiffByLine(a, b, Options{IgnoreWhitespace: true, ContextLines: -1, WithLineNo: false})
	fmt.Println(ret)

	// 预期: 忽略空白字符后应该没有差异
	if ret != "" {
		t.Error("预期忽略空白字符后没有差异，但实际有差异")
	}

	fmt.Println("\n=== 测试4: 不忽略空白字符 ===")
	ret = DiffByLine(a, b, Options{IgnoreWhitespace: false, ContextLines: -1, WithLineNo: false})
	fmt.Println(ret)

	// 预期: 不忽略空白字符应该有差异
	if ret == "" {
		t.Error("预期不忽略空白字符应该有差异")
	}
	if !strings.Contains(ret, "-|") || !strings.Contains(ret, "+|") {
		t.Error("预期应该包含删除和添加标记")
	}
}

// 测试用例5: 多处修改，测试上下文显示
func TestDiffByLine_MultipleChangesWithContext(t *testing.T) {
	a := `line 1
line 2
line 3
line 4
line 5
line 6
line 7
line 8
line 9
line 10
line 11
line 12
line 13
line 14
line 15
line 16
line 17
line 18
line 19
line 20`

	b := `line 1
line 2 modified
line 3
line 4
line 5
line 6
line 7
line 8
line 9
line 10
line 11
line 12 modified
line 13
line 14
line 15
line 16
line 17
line 18 modified
line 19
line 20`

	fmt.Println("=== 测试5: 多处修改，上下文1行 ===")
	ret := DiffByLine(a, b, Options{IgnoreWhitespace: false, ContextLines: 1, WithLineNo: true})
	fmt.Println(ret)

	// 预期检查
	ellipsisCount := strings.Count(ret, "...")
	if ellipsisCount < 2 {
		t.Errorf("预期至少有2个省略标记，实际有 %d 个", ellipsisCount)
	}

	// 应该包含修改的行
	if !strings.Contains(ret, "line 2 modified") {
		t.Error("预期包含 'line 2 modified'")
	}
	if !strings.Contains(ret, "line 12 modified") {
		t.Error("预期包含 'line 12 modified'")
	}
	if !strings.Contains(ret, "line 18 modified") {
		t.Error("预期包含 'line 18 modified'")
	}

	fmt.Println("\n=== 测试5: 多处修改，上下文3行 ===")
	ret = DiffByLine(a, b, Options{IgnoreWhitespace: false, ContextLines: 3, WithLineNo: true})
	fmt.Println(ret)

	// 预期: 上下文3行应该显示更多内容
	lineCount := strings.Count(ret, "\n")
	if lineCount < 15 {
		t.Errorf("预期至少显示15行，实际显示 %d 行", lineCount)
	}
}

// 测试用例6: 代码重构场景
func TestDiffByLine_CodeRefactor(t *testing.T) {
	a := `package main

import "fmt"

func oldFunction(x int) int {
	return x * 2
}

func main() {
	result := oldFunction(5)
	fmt.Println(result)
}`

	b := `package main

import "fmt"

func newFunction(x int) int {
	// 新的实现
	return x * 3
}

func main() {
	result := newFunction(5)
	fmt.Println("Result:", result)
}`

	fmt.Println("=== 测试6: 代码重构，上下文2行 ===")
	ret := DiffByLine(a, b, Options{IgnoreWhitespace: false, ContextLines: 2, WithLineNo: true})
	fmt.Println(ret)

	// 预期检查
	if !strings.Contains(ret, "oldFunction") {
		t.Error("预期包含 'oldFunction' 的删除")
	}
	if !strings.Contains(ret, "newFunction") {
		t.Error("预期包含 'newFunction' 的添加")
	}
	if !strings.Contains(ret, "...") {
		t.Error("预期包含省略标记")
	}
}

// 测试用例7: 空文件和非空文件
func TestDiffByLine_EmptyFile(t *testing.T) {
	a := ``
	b := `line 1
line 2
line 3`

	fmt.Println("=== 测试7: 空文件 -> 非空文件 ===")
	ret := DiffByLine(a, b, Options{IgnoreWhitespace: false, ContextLines: -1, WithLineNo: true})
	fmt.Println(ret)

	// 预期: 全部是插入
	insertCount := strings.Count(ret, "+|")
	if insertCount != 3 {
		t.Errorf("预期有3个插入操作，实际有 %d 个", insertCount)
	}

	fmt.Println("\n=== 测试7: 非空文件 -> 空文件 ===")
	ret = DiffByLine(b, a, Options{IgnoreWhitespace: false, ContextLines: -1, WithLineNo: true})
	fmt.Println(ret)

	// 预期: 全部是删除
	deleteCount := strings.Count(ret, "-|")
	if deleteCount != 3 {
		t.Errorf("预期有3个删除操作，实际有 %d 个", deleteCount)
	}
}

// 测试用例8: 完全相同的文件
func TestDiffByLine_Identical(t *testing.T) {
	a := `line 1
line 2
line 3`
	b := `line 1
line 2
line 3`

	fmt.Println("=== 测试8: 完全相同的文件 ===")
	ret := DiffByLine(a, b, Options{IgnoreWhitespace: false, ContextLines: -1, WithLineNo: true})
	if ret == "" {
		fmt.Println("(无差异)")
	} else {
		fmt.Println(ret)
	}

	// 预期: 没有差异
	if ret != "" {
		t.Error("预期完全相同的文件应该没有差异")
	}
}

// 测试用例9: JSON格式化差异
func TestDiffByLine_JSONFormat(t *testing.T) {
	a := `{"name":"John","age":30,"city":"New York"}`
	b := `{
  "name": "John",
  "age": 30,
  "city": "New York"
}`

	fmt.Println("=== 测试9: JSON格式化 (忽略空白) ===")
	ret := DiffByLine(a, b, Options{IgnoreWhitespace: true, ContextLines: -1, WithLineNo: false})
	fmt.Println(ret)

	// 预期: 忽略空白后应该没有差异
	if ret != "" {
		t.Error("预期忽略空白字符后JSON格式化差异应该被忽略")
	}

	fmt.Println("\n=== 测试9: JSON格式化 (不忽略空白) ===")
	ret = DiffByLine(a, b, Options{IgnoreWhitespace: false, ContextLines: -1, WithLineNo: false})
	fmt.Println(ret)

	// 预期: 不忽略空白应该有差异
	if ret == "" {
		t.Error("预期不忽略空白字符应该检测到JSON格式化差异")
	}
}

// 测试用例10: 连续修改和间隔修改
func TestDiffByLine_ConsecutiveAndSeparateChanges(t *testing.T) {
	a := `line 1
line 2
line 3
line 4
line 5
line 6
line 7
line 8
line 9
line 10
line 11
line 12
line 13
line 14
line 15`

	b := `line 1
line 2 modified
line 3 modified
line 4 modified
line 5
line 6
line 7
line 8
line 9
line 10
line 11
line 12
line 13 modified
line 14
line 15`

	fmt.Println("=== 测试10: 连续修改和间隔修改，上下文1行 ===")
	ret := DiffByLine(a, b, Options{IgnoreWhitespace: false, ContextLines: 1, WithLineNo: true})
	fmt.Println(ret)

	// 预期检查
	// 应该有省略标记（在连续修改和单独修改之间）
	if !strings.Contains(ret, "...") {
		t.Error("预期包含省略标记")
	}

	// 应该包含所有修改
	modifiedCount := strings.Count(ret, "modified")
	if modifiedCount != 4 {
		t.Errorf("预期包含4处修改，实际有 %d 处", modifiedCount)
	}
}

// 测试用例11: 文件末尾修改
func TestDiffByLine_EndOfFile(t *testing.T) {
	a := `line 1
line 2
line 3
line 4
line 5`

	b := `line 1
line 2
line 3
line 4
line 5 modified
line 6 added`

	fmt.Println("=== 测试11: 文件末尾修改，上下文2行 ===")
	ret := DiffByLine(a, b, Options{IgnoreWhitespace: false, ContextLines: 2, WithLineNo: true})
	fmt.Println(ret)

	// 预期检查
	if !strings.Contains(ret, "line 5 modified") {
		t.Error("预期包含 'line 5 modified'")
	}
	if !strings.Contains(ret, "line 6 added") {
		t.Error("预期包含 'line 6 added'")
	}
	// 应该包含上下文
	if !strings.Contains(ret, "line 3") || !strings.Contains(ret, "line 4") {
		t.Error("预期包含上下文行 line 3 和 line 4")
	}
}

// 测试用例12: 中文内容
func TestDiffByLine_Chinese(t *testing.T) {
	a := `第一行
第二行
第三行
第四行`

	b := `第一行
第二行已修改
第三行
第四行
第五行新增`

	fmt.Println("=== 测试12: 中文内容 ===")
	ret := DiffByLine(a, b, Options{IgnoreWhitespace: false, ContextLines: -1, WithLineNo: true})
	fmt.Println(ret)

	// 预期检查
	if !strings.Contains(ret, "第二行已修改") {
		t.Error("预期包含 '第二行已修改'")
	}
	if !strings.Contains(ret, "第五行新增") {
		t.Error("预期包含 '第五行新增'")
	}

	// 检查中文字符没有被破坏
	if !strings.Contains(ret, "第一行") {
		t.Error("预期正确显示中文字符")
	}
}

// 测试用例13: 边界情况 - 只有一行
func TestDiffByLine_SingleLine(t *testing.T) {
	a := `single line`
	b := `single line modified`

	fmt.Println("=== 测试13: 单行修改 ===")
	ret := DiffByLine(a, b, Options{IgnoreWhitespace: false, ContextLines: 1, WithLineNo: true})
	fmt.Println(ret)

	// 预期检查
	if !strings.Contains(ret, "-|") || !strings.Contains(ret, "+|") {
		t.Error("预期包含删除和添加标记")
	}
}

// 测试用例14: 上下文为0
func TestDiffByLine_ZeroContext(t *testing.T) {
	a := `line 1
line 2
line 3
line 4
line 5`

	b := `line 1
line 2
line 3 modified
line 4
line 5`

	fmt.Println("=== 测试14: 上下文为0 ===")
	ret := DiffByLine(a, b, Options{IgnoreWhitespace: false, ContextLines: 0, WithLineNo: true})
	fmt.Println(ret)

	// 预期: 只显示修改的行，不显示上下文
	lineCount := strings.Count(ret, "\n")
	if lineCount > 3 { // 最多应该是删除行+添加行+可能的省略标记
		t.Errorf("预期上下文为0时只显示修改行，实际显示了 %d 行", lineCount)
	}

	if !strings.Contains(ret, "line 3 modified") {
		t.Error("预期包含修改的行")
	}
}
