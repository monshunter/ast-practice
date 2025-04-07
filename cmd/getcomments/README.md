# golang文件注释提取工具

## 功能介绍
这个工具可以从Go源代码文件中提取每一行代码所关联的注释。它通过Go的AST（抽象语法树）工具实现，能够准确关联代码与其相关注释。

## 输入与输出
### 输入: 
- 一个golang文件路径
- 或者一段golang代码

### 输出: 
JSON格式的键值对:
```json
"文件名:执行代码的行号": ["注释1", "注释2", ...]
```

## 安装方法
1. 克隆本仓库
2. 进入cmd/getcomments目录
3. 执行build.sh脚本构建可执行文件:
```bash
chmod +x build.sh
./build.sh
```

## 使用方法
```bash
# 从文件中提取注释
./getcomments /path/to/your/gofile.go

# 在测试环境中运行
cd examples
go test -v
```

## 示例

### 输入文件
```golang
package examples

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
	fmt.Println(x, y)
	// 判断变量大小
	if x > y { // if 分支
		fmt.Println("x is greater than y")
		// 否则进入 else 分支
	} else { // else 分支
		// 打印输出： x is less than y
		fmt.Println("x is less than y") // 行尾注释打印输出： x is less than y
	}
}
```

### 输出结果

```json
{
    "examples.go:11": [
        "// ExampleGetComments 示例",
        "// 输入:",
        "// 一段golang代码",
        "// 输出:",
        "// 文件名:行号: []string{comments...}",
        "// 函数名: 函数注释 ExampleGetComments"
    ],
    "examples.go:13": [
        "// 定义两个变量"
    ],
    "examples.go:15": [
        "// 打印变量"
    ],
    "examples.go:17": [
        "// 判断变量大小",
        "// if 分支"
    ],
    "examples.go:20": [
        "// 否则进入 else 分支",
        "// else 分支"
    ],
    "examples.go:22": [
        "// 打印输出： x is less than y",
        "// 行尾注释打印输出： x is less than y"
    ]
}
```

## 实现原理
该工具使用Go语言的抽象语法树(AST)功能解析源代码，实现了以下关键功能：

1. **代码解析**: 使用`go/parser`包解析Go源代码，生成AST
2. **节点遍历**: 使用`ast.Inspect`遍历AST中的所有节点
3. **注释关联**: 根据代码节点的位置信息，关联相关的注释组
4. **类型识别**: 识别函数声明、赋值语句、表达式语句等不同类型的代码节点
5. **上下文关联**: 能够关联代码行附近的相关注释，包括行内注释和上方的注释块

该工具特别关注以下类型的节点：
- 函数声明 (*ast.FuncDecl)
- 赋值语句 (*ast.AssignStmt)
- 表达式语句 (*ast.ExprStmt)
- 条件语句 (*ast.IfStmt)
- 循环语句 (*ast.ForStmt)
- 等其他常见的代码结构

## 注意事项
1. 工具会尝试智能判断哪些注释与代码相关，但在复杂情况下可能存在误判
2. 注释与代码的关联基于行号和位置信息，调整代码格式可能影响结果
3. 目前支持行注释(`//`)，块注释(`/* */`)也会被解析，但可能不如行注释精确 