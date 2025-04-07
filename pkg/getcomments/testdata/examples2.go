package testdata

// 导入标准库
import (
	"fmt"
	"math"
)

// ExampleGetComments 示例
// 输入:
// 一段golang代码
// 输出:
// 文件名:行号: []string{comments...}
func ExampleGetComments2() { // 函数名: 函数注释 ExampleGetComments
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
	// 定义一个函数
	var f func(i int) bool
	// 函数赋值
	f = func(i int) bool { // 函数注释
		// 打印输出： i > 0
		return i > 0 // 函数返回值
	}
	// 打印输出： true
	fmt.Println(f(1))
	// 延迟执行
	defer func(x int) { // 延迟执行注释
		// 打印输出： defer func
		fmt.Println("defer func", x)
		// 结束延迟执行
	}(2)
	// 启动一个协程
	go func() { // 协程注释
		// 打印输出： go func
		fmt.Println("go func")
	}() // 协程注释2
}

// ExampleGetComments3 示例
// 输入:
// 一段golang代码
// 输出:
// 文件名:行号: []string{comments...}
func ExampleGetComments3() {
	// 定义一个函数
	var f func(i int) bool
	// 函数赋值
	f = func(i int) bool { // 函数注释
		return i > 0 // 函数返回值
	}
	// 打印输出： true
	fmt.Println(f(1))
	// 定义一个通道
	var ch chan int
	// 通道赋值
	ch = make(chan int, 1)
	// 通道写入
	ch <- 1
	// 通道读取
	fmt.Println(<-ch)
	// 通道选择
	select { // 通道选择注释
	case <-ch: // 通道选择注释2
		fmt.Println("ch")
	default: // 通道选择注释3
		// 打印输出： default
		fmt.Println("default")
	}

	// 定义一个结构体
	type S struct {
		Name string // 姓名
		Age  int    // 年龄
	}
	// 打印输出： { 0}
	fmt.Println(S{})
}

// INF 常量
const INF = math.MaxInt64

// StaticVar 静态变量
var StaticVar = 2

// 定义一个结构体
type Person struct {
	Name string // 姓名
	Age  int    // 年龄
}

// 定义一个接口
type Animal interface {
	Say() string // 说话
}

// 实现接口
func (p *Person) Say() string {
	// 打印输出： Hello,
	return "Hello, " + p.Name
}

// 定义一个数组
// 数组注释
var ARR = []int{1, 2, 3} // 行尾数组注释

var ARR2 = []int{
	1, // 这是1
	2, // 这是2
	3, // 这是3
} // 行尾数组注释2
// 定义一个map
// map注释
var MAP = map[string]int{
	"a": 1, // map注释a
	"b": 2, // map注释b
}
