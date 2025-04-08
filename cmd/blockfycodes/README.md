## 代码分支化工具 (Blockfycodes)

### 功能
- 分析Go代码的抽象语法树(AST)
- 在代码中插入新的语句（fmt.Println打印语句）
- 在代码中插入注释
- 支持各种Go语言结构（if-else、for、switch等）的处理
- 接受文件路径或直接代码内容作为输入

### 用法
```
blockfycodes <文件路径或代码内容>
```

### 工作原理
1. **解析输入**：接受文件路径或直接的代码内容作为输入
2. **语句插入**：在各种Go语言结构（如if语句、for循环、函数等）前插入fmt.Println语句
3. **注释插入**：在代码的各个部分插入"// + insert comment test"注释
4. **代码输出**：将修改后的代码输出到控制台

### 支持的语法结构
- 赋值语句 (AssignStmt)
- 条件语句 (IfStmt)
- 循环语句 (ForStmt, RangeStmt)
- 分支语句 (SwitchStmt, TypeSwitchStmt)
- 代码块 (BlockStmt)
- 返回语句 (ReturnStmt)
- 延迟执行 (DeferStmt)
- 并发执行 (GoStmt)
- 通道选择 (SelectStmt)
- 表达式语句 (ExprStmt)

### 实现细节
- 使用Go标准库的`go/ast`、`go/parser`和`go/token`进行代码解析
- 使用`go/printer`将修改后的AST转换回代码
- 递归遍历AST，支持嵌套的语法结构
- 保持代码格式和缩进风格