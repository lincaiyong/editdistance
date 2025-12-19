package edittool

import (
	"fmt"
	"testing"
)

// 测试用例1: 基本的增删改操作
func TestDiffByLine_BasicOperations(t *testing.T) {
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

	fmt.Println("=== 测试1: 基本增删改 (显示所有行，带行号) ===")
	ret := DiffByLine(a, b, Options{IgnoreWhitespace: false, ContextLines: -1, WithLineNo: true})
	fmt.Println(ret)

	expected := ` |0001:0001|	line 1
-|0002:    |	line 2
-|0003:    |	line 3
+|    :0002|	line 2 modified
 |0004:0003|	line 4
 |0005:0004|	line 5
+|    :0005|	line 6
`
	if ret != expected {
		t.Errorf("输出不匹配\n期望:\n%s\n实际:\n%s", expected, ret)
	}

	fmt.Println("\n=== 测试1: 基本增删改 (上下文2行) ===")
	ret = DiffByLine(a, b, Options{IgnoreWhitespace: false, ContextLines: 2, WithLineNo: true})
	fmt.Println(ret)

	// 上下文2行应该显示所有内容（因为文件较短）
	if ret != expected {
		t.Errorf("上下文2行输出不匹配\n期望:\n%s\n实际:\n%s", expected, ret)
	}
}

// 测试用例2: 空白字符处理
func TestDiffByLine_WhitespaceHandling(t *testing.T) {
	a := `	if (condition) {
		doSomething();
	}`
	b := `if (condition) {
    doSomething();
}`

	fmt.Println("=== 测试2: 忽略空白字符 ===")
	ret := DiffByLine(a, b, Options{IgnoreWhitespace: true, ContextLines: 0, WithLineNo: false})
	fmt.Println(ret)

	expected := ""
	if ret != expected {
		t.Errorf("忽略空白字符输出不匹配\n期望: (空字符串)\n实际:\n%s", ret)
	}

	fmt.Println("\n=== 测试2: 不忽略空白字符 ===")
	ret = DiffByLine(a, b, Options{IgnoreWhitespace: false, ContextLines: -1, WithLineNo: true})
	fmt.Println(ret)

	expected = `-|0001:    |		if (condition) {
-|0002:    |			doSomething();
-|0003:    |		}
+|    :0001|	if (condition) {
+|    :0002|	    doSomething();
+|    :0003|	}
`
	if ret != expected {
		t.Errorf("不忽略空白字符输出不匹配\n期望:\n%s\n实际:\n%s", expected, ret)
	}
}

// 测试用例3: 上下文行数控制和省略标记
func TestDiffByLine_ContextControl(t *testing.T) {
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

	fmt.Println("=== 测试3: 多处修改，上下文1行 ===")
	ret := DiffByLine(a, b, Options{IgnoreWhitespace: false, ContextLines: 1, WithLineNo: true})
	fmt.Println(ret)

	expected := ` |0001:0001|	line 1
-|0002:    |	line 2
+|    :0002|	line 2 modified
 |0003:0003|	line 3
 | ...
 |0011:0011|	line 11
-|0012:    |	line 12
+|    :0012|	line 12 modified
 |0013:0013|	line 13
 | ...
 |0017:0017|	line 17
-|0018:    |	line 18
+|    :0018|	line 18 modified
 |0019:0019|	line 19
`
	if ret != expected {
		t.Errorf("上下文1行输出不匹配\n期望:\n%s\n实际:\n%s", expected, ret)
	}

	fmt.Println("\n=== 测试3: 完全相同的文件 ===")
	ret = DiffByLine(a, a, Options{IgnoreWhitespace: false, ContextLines: 1, WithLineNo: true})
	fmt.Println(ret)

	expected = ""
	if ret != expected {
		t.Errorf("完全相同文件输出不匹配\n期望: (空字符串)\n实际:\n%s", ret)
	}

	fmt.Println("\n=== 测试3: 上下文为0 ===")
	ret = DiffByLine(a, b, Options{IgnoreWhitespace: false, ContextLines: 0, WithLineNo: true})
	fmt.Println(ret)

	expected = `-|0002:    |	line 2
+|    :0002|	line 2 modified
 | ...
-|0012:    |	line 12
+|    :0012|	line 12 modified
 | ...
-|0018:    |	line 18
+|    :0018|	line 18 modified
`
	if ret != expected {
		t.Errorf("上下文0行输出不匹配\n期望:\n%s\n实际:\n%s", expected, ret)
	}
}

func TestParseDiff(t *testing.T) {
	a := `line 1
line 2 modified
line x
line 4
line 7
line 8
line 9
line 11`

	b := ``
	diff := GeneratePatch(a, b)
	fmt.Println(diff)
	bb := Patch(a, diff)
	if b != bb {
		t.Errorf("<UNK>\n<UNK>:\n%s\n<UNK>:\n%s", b, bb)
	}
}
