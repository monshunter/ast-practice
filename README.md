# Go AST 实践与探索

## 简介

本项目旨在通过实践探索 Go 语言的抽象语法树 (AST) 能力。通过开发一系列基于 AST 的小工具，加深对 Go 代码结构的理解，并为未来利用 AST 实现更高级的代码分析、转换或生成功能奠定基础。

这是一个学习和验证想法的过程，目的是掌握 AST 相关的技术，以便将来能将其应用于更有趣、更有价值的场景。

## 动机

创建这些工具的主要目的是学习和掌握 Go AST 的使用方法，验证通过 AST 分析和修改 Go 代码的可行性。虽然 `getcomments` 和 `blockfycodes` 这两个工具本身功能较为基础，但它们是理解和应用 AST 的重要练习。

## 包含的工具

本项目目前包含以下基于 AST 的工具：

1.  **`cmd/getcomments`**:
    *   **功能**: 从 Go 源代码文件中提取每一行代码所关联的注释。
    *   **详情**: 请参阅 [`cmd/getcomments/README.md`](cmd/getcomments/README.md) 获取安装、使用和示例。

2.  **`cmd/blockfycodes`**:
    *   **功能**: 分析 Go 代码的 AST，并在代码的关键结构（如 `if`, `for`, `switch`, 函数体等）中插入指定的语句（例如 `fmt.Println`）或注释。
    *   **详情**: 请参阅 [`cmd/blockfycodes/README.md`](cmd/blockfycodes/README.md) 获取使用方法和实现细节。

## 未来方向

本项目中的实践经验将作为基础，用于未来开发更多基于 AST 的工具，可能性包括但不限于：

*   代码静态分析工具
*   自动化代码重构
*   领域特定语言 (DSL) 的解析与处理
*   代码生成器

## 如何运行

请参考各个子目录 (`cmd/getcomments`, `cmd/blockfycodes`) 下的 `README.md` 文件，获取对应工具的具体安装、构建和运行说明。
