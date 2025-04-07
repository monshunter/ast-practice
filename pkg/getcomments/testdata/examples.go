package testdata

// 导入标准库
import "fmt"

// ExampleGetComments 示例
// 输入:
// 一段golang代码
// 输出:
// 文件名:行号: []string{comments...}
func ExampleGetComments() { // 函数名: 函数注释 ExampleGetComments
	// 定义两个变量
	x, y := 1, 2
	// 打印变量
	fmt.Println("x, y\n", x, y)
	// 判断变量大小
	if x > y { // if 分支
		fmt.Println("x is greater than y")
		// 否则进入 else 分支
	} else { // else 分支
		// 打印输出： x is less than y
		fmt.Println("x is less than y") // 行尾注释打印输出： x is less than y
	}
}
