# golang main 入口发现

## 功能介绍
利用golang ast分析从一个golang project中找到所有main 入口

## 输入与输出
### 输入: 
- 一个golang文件路径

### 输出: 
JSON格式的键值对:
```json
[{"dir":"main入口所在目录（相对于project）", "file":"main入口所在的文件"},]
```

## 安装

```bash
go install github.com/monshunter/ast-practice/cmd/findmain@latest
```

或者从源码编译：

```bash
git clone https://github.com/your-username/ast-practice.git
cd ast-practice
go build -o findmain ./cmd/findmain
```

## 使用方法

```bash
findmain [选项] <项目路径>
```

### 选项:
- `-v`: 启用详细输出模式，显示更多的分析过程信息
- `-o`: 输出文件路径（默认为标准输出）
- `-h`: 显示帮助信息

### 示例:

基本用法：
```bash
findmain /path/to/go/project
```

将结果保存到文件：
```bash
findmain -o result.json /path/to/go/project
```

启用详细输出模式：
```bash
findmain -v /path/to/go/project
```

## 工作原理

该工具通过以下步骤找到golang项目中的所有main入口：

1. 遍历项目中的所有目录和文件
2. 筛选出所有Go源文件（.go后缀）
3. 使用Go的AST（抽象语法树）分析每个文件
4. 检查文件是否同时满足以下条件：
   - 包声明为`package main`
   - 包含`func main()`函数声明
5. 收集符合条件的文件信息，包括相对目录路径和文件名
6. 以JSON格式输出结果

## 性能考虑

- 自动跳过常见的非Go代码目录（vendor, .git, node_modules）
- 对于访问权限不足或格式错误的文件，会记录警告并继续处理其他文件
- 对于大型项目，可能需要一定的处理时间，请耐心等待