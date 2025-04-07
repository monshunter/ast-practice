package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// 保存代码行及其关联的注释
type CommentsMap map[string][]string

type comment struct {
	line    int
	content string
}

// 从文件路径或代码内容提取注释
func extractComments(input string) (CommentsMap, error) {
	var content []byte
	var err error
	var filename string

	// 判断输入是文件路径还是代码内容
	if _, err := os.Stat(input); err == nil {
		// 输入是文件路径
		filename = filepath.Base(input)
		content, err = os.ReadFile(input)
		if err != nil {
			return nil, fmt.Errorf("读取文件失败: %v", err)
		}
	} else {
		// 输入是代码内容
		filename = "code.go"
		content = []byte(input)
	}

	// 解析Go代码
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, content, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("解析代码失败: %v", err)
	}

	// 将内容分割为行
	lines := strings.Split(string(content), "\n")

	// 创建注释映射
	commentsMap := make(CommentsMap)

	// 通过TokenFile直接处理每一行
	tf := fset.File(f.Pos())
	if tf == nil {
		return nil, fmt.Errorf("获取token文件失败")
	}

	// 创建行注释映射
	lineComments := make(map[int]string)

	// 收集所有注释
	for _, cg := range f.Comments {
		for _, comment := range cg.List {
			pos := fset.Position(comment.Pos())
			lineComments[pos.Line] = comment.Text
		}
	}

	// 处理函数声明
	for _, decl := range f.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			funcPos := fset.Position(funcDecl.Pos())
			lineKey := fmt.Sprintf("%s:%d", filename, funcPos.Line)

			var comments []comment

			// 收集函数声明行的注释
			if lineComments[funcPos.Line] != "" {
				comments = append(comments, comment{line: funcPos.Line, content: lineComments[funcPos.Line]})
				delete(lineComments, funcPos.Line)
			}

			// 收集上方的文档注释
			// 先尝试doc注释
			if funcDecl.Doc != nil {
				for _, cm := range funcDecl.Doc.List {
					if cm.Text != "" {
						thisLine := fset.Position(cm.Pos()).Line
						comments = append(comments, comment{line: thisLine, content: cm.Text})
						delete(lineComments, thisLine)
					}
				}
			} else {
				// 尝试查找上方连续的行注释
				startLine := funcPos.Line - 1
				for startLine > 0 {
					if cmts, ok := lineComments[startLine]; ok && cmts != "" {
						comments = append([]comment{{line: startLine, content: cmts}}, comments...)
						delete(lineComments, startLine)
					} else if startLine < len(lines) && strings.TrimSpace(lines[startLine-1]) != "" {
						// 非空行且没有注释，停止搜索
						break
					}
					startLine--
				}
			}

			// 去重并存储
			if len(comments) > 0 {
				commentsMap[lineKey] = comment2String(uniqueComments(comments))
			}
		}
	}

	// 处理其他语句
	ast.Inspect(f, func(n ast.Node) bool {
		if n == nil || isNodeType(n, "FuncDecl") {
			return true // 跳过nil和已处理的函数声明
		}

		if isCodeNode(n) {
			pos := fset.Position(n.Pos())
			if pos.Line == 0 {
				return true // 跳过无效位置
			}

			lineKey := fmt.Sprintf("%s:%d", filename, pos.Line)
			// 如果已经存在，则跳过
			if commentsMap[lineKey] != nil {
				return true
			}

			var comments []comment

			// 收集当前行的注释
			if lineComments[pos.Line] != "" {
				comments = append(comments, comment{line: pos.Line, content: lineComments[pos.Line]})
				delete(lineComments, pos.Line)
			}

			// 收集上一行的注释
			if pos.Line > 1 {
				prevLine := pos.Line - 1
				// 只往上寻找最近的注释，遇到非空行停止
				for prevLine > 0 {
					if cmts, ok := lineComments[prevLine]; ok && cmts != "" {
						comments = append([]comment{{line: prevLine, content: cmts}}, comments...)
						delete(lineComments, prevLine)
						break // 找到第一个注释行就停止
					} else if prevLine < len(lines) && strings.TrimSpace(lines[prevLine-1]) != "" {
						// 非空行且没有注释，停止搜索
						break
					}
					prevLine--
				}
			}

			// 去重并存储
			if len(comments) > 0 {
				commentsMap[lineKey] = comment2String(uniqueComments(comments))
			}
		}

		return true
	})

	return commentsMap, nil
}

// 判断节点是否为指定类型
func isNodeType(n ast.Node, typeName string) bool {
	nodeType := fmt.Sprintf("%T", n)
	return strings.Contains(nodeType, typeName)
}

// 判断是否为关键代码节点
func isCodeNode(n ast.Node) bool {
	switch n.(type) {
	case *ast.AssignStmt, *ast.ExprStmt, *ast.IfStmt, *ast.BranchStmt,
		*ast.ReturnStmt, *ast.DeclStmt, *ast.BlockStmt, *ast.ForStmt:
		return true
	}
	return false
}

func uniqueComments(comments []comment) []comment {
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].line < comments[j].line
	})

	slow, fast := 0, 0
	for fast < len(comments) {
		if comments[slow].line == comments[fast].line {
			fast++
		} else {
			slow++
			comments[slow] = comments[fast]
		}
	}
	return comments[:slow+1]
}

func comment2String(comments []comment) []string {
	result := []string{}
	for _, comment := range comments {
		result = append(result, comment.content)
	}
	return result
}

func main() {
	// 检查参数
	if len(os.Args) != 2 {
		fmt.Println("用法: getcomments <文件路径或代码内容>")
		os.Exit(1)
	}

	input := os.Args[1]

	// 提取注释
	commentsMap, err := extractComments(input)
	if err != nil {
		fmt.Printf("提取注释失败: %v\n", err)
		os.Exit(1)
	}

	// 输出JSON格式结果
	output, err := json.MarshalIndent(commentsMap, "", "    ")
	if err != nil {
		fmt.Printf("序列化结果失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}
