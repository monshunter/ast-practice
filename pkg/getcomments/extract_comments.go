package getcomments

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
// ExtractComments// 从文件路径或代码内容提取注释并返回注释映射
// 算法思路：
// 1. 判断输入类型（文件路径或代码内容）
// 2. 解析Go代码并获取所有注释
// 3. 遍历AST节点收集关联注释
// 4. 处理特殊节点（函数、结构体、接口等）
//
// 时间复杂度分析：
// - 文件读取：O(n)
// - 代码解析：O(n)
// - AST遍历：O(n)
// - 注释收集：O(m) m为注释数量
// 总体复杂度：O(n + m)
//
// 优化建议：
// 1. 可缓存解析结果避免重复解析相同文件
// 2. 并行处理多个文件
// 3. 优化行注释收集逻辑，减少内存分配
// 4. 添加注释过滤选项（如只收集文档注释）
func ExtractComments(input string) (CommentsMap, error) {
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
	// 创建行注释映射
	lineComments := make([]string, len(lines)+1)
	lineCommentsVisited := make([]bool, len(lines)+1)

	// 收集所有注释
	for _, cg := range f.Comments {
		for _, comment := range cg.List {
			lineComments[fset.Position(comment.Pos()).Line] = comment.Text
		}
	}

	// 处理函数声明
	for _, decl := range f.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			funcPos := fset.Position(funcDecl.Pos())
			lineKey := fmt.Sprintf("%s:%d", filename, funcPos.Line)
			// 收集本行的文档注释
			currentLines := []int{funcPos.Line}
			// 先尝试doc注释
			if funcDecl.Doc != nil {
				for _, cm := range funcDecl.Doc.List {
					currentLines = append(currentLines, fset.Position(cm.Pos()).Line)
				}
			} else {
				// 尝试查找上方连续的行注释
				startLine := funcPos.Line - 1
				for startLine > 0 && strings.TrimSpace(lines[startLine]) != "" {
					currentLines = append(currentLines, startLine)
					startLine--
				}
			}
			comments := collectCurrentLineComments(currentLines, lineComments, lineCommentsVisited)
			// 去重并存储
			if len(comments) > 0 && commentsMap[lineKey] == nil {
				commentsMap[lineKey] = comment2String(comments)
			}
		} else if genDecl, ok := decl.(*ast.GenDecl); ok {
			// 处理常量、变量和类型声明
			for _, spec := range genDecl.Specs {
				var pos token.Pos

				switch s := spec.(type) {
				case *ast.TypeSpec:
					pos = s.Pos()
				case *ast.ValueSpec:
					pos = s.Pos()
				default:
					continue
				}

				specPos := fset.Position(pos)
				lineKey := fmt.Sprintf("%s:%d", filename, specPos.Line)
				currentLines := []int{specPos.Line}

				// 尝试使用Doc注释
				if genDecl.Doc != nil {
					for _, cm := range genDecl.Doc.List {
						currentLines = append(currentLines, fset.Position(cm.Pos()).Line)
					}
				} else {
					// 尝试查找上方连续的行注释
					startLine := specPos.Line - 1
					for startLine > 0 && strings.TrimSpace(lines[startLine]) != "" {
						currentLines = append(currentLines, startLine)
						startLine--
					}
				}
				comments := collectCurrentLineComments(currentLines, lineComments, lineCommentsVisited)
				// 去重并存储
				if len(comments) > 0 && commentsMap[lineKey] == nil {
					commentsMap[lineKey] = comment2String(comments)
				}

				// 对于结构体类型，处理字段注释
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						collectFieldsComments(fset, filename, structType.Fields.List, lineComments, lineCommentsVisited, commentsMap)
					}
				}
			}
		}
	}

	// 处理其他语句
	ast.Inspect(f, func(n ast.Node) bool {
		if n == nil {
			return true
		}

		// 获取节点位置
		pos := fset.Position(n.Pos())
		if pos.Line == 0 {
			return true // 跳过无效位置
		}

		lineKey := fmt.Sprintf("%s:%d", filename, pos.Line)

		// 收集当前行的注释
		currentLines := []int{pos.Line, pos.Line - 1}
		// 收集上一行的注释
		comments := collectCurrentLineComments(currentLines, lineComments, lineCommentsVisited)
		// 去重并存储
		if len(comments) > 0 && commentsMap[lineKey] == nil {
			commentsMap[lineKey] = comment2String(comments)
		}

		// 处理特殊节点类型
		switch node := n.(type) {
		case *ast.InterfaceType:
			if node.Methods != nil {
				collectFieldsComments(fset, filename, node.Methods.List, lineComments, lineCommentsVisited, commentsMap)
			}
		case *ast.StructType:
			if node.Fields != nil {
				collectFieldsComments(fset, filename, node.Fields.List, lineComments, lineCommentsVisited, commentsMap)
			}
		}

		return true
	})

	return commentsMap, nil
}

// collectFieldsComments 评估正确性：
// 1. 正确性验证：
// - 正确遍历结构体/接口字段
// - 准确获取字段位置信息
// - 正确收集字段关联注释
// - 避免重复收集已处理注释
//
// 算法思路：
// 1. 遍历所有字段节点
// 2. 获取每个字段的位置信息
// 3. 收集字段所在行的注释
// 4. 去重后存入结果映射
//
// 时间复杂度分析：
// - 字段遍历：O(n) n为字段数量
// - 位置获取：O(1)
// - 注释收集：O(1) 因使用预计算的lineComments
// 总体复杂度：O(n)
//
// 优化建议：
// 1. 批量预计算字段位置信息
// 2. 并行处理字段注释收集
// 3. 添加字段注释类型过滤（如只收集文档注释）
// 4. 优化内存分配，预分配结果切片
// 5. 添加字段名称匹配过滤选项
func collectFieldsComments(fset *token.FileSet, filename string, fields []*ast.Field,
	lineComments []string, lineCommentsVisited []bool, commentsMap CommentsMap) {
	for _, field := range fields {
		fieldPos := fset.Position(field.Pos())
		fieldKey := fmt.Sprintf("%s:%d", filename, fieldPos.Line)
		fieldComments := collectCurrentLineComments([]int{fieldPos.Line}, lineComments, lineCommentsVisited)
		// 去重并存储
		if len(fieldComments) > 0 && commentsMap[fieldKey] == nil {
			commentsMap[fieldKey] = comment2String(fieldComments)
		}
	}
}

// collectCurrentLineComments 评估正确性：
// 1. 正确性验证：
// - 正确排序输入行号
// - 准确过滤已访问注释行
// - 正确收集非空注释内容
// - 维护已访问状态避免重复收集
//
// 算法思路：
// 1. 对输入行号进行排序
// 2. 预分配结果切片容量
// 3. 遍历所有行号
// 4. 检查行注释是否有效且未访问
// 5. 收集有效注释并标记已访问
//
// 时间复杂度分析：
// - 行号排序：O(n log n) n为输入行数
// - 行遍历：O(n)
// - 注释检查：O(1)
// 总体复杂度：O(n log n)
//
// 优化建议：
// 1. 如果输入行号已排序可跳过排序步骤
// 2. 使用位图替代bool切片减少内存占用
// 3. 并行处理行号遍历（需解决数据竞争）
// 4. 预计算有效注释行减少运行时检查
// 5. 添加注释内容过滤选项（如只收集特定前缀注释）
func collectCurrentLineComments(lines []int, lineComments []string, lineCommentsVisited []bool) []comment {
	sort.Ints(lines)
	comments := make([]comment, 0, len(lines))
	for _, line := range lines {
		if !lineCommentsVisited[line] && lineComments[line] != "" {
			comments = append(comments, comment{line: line, content: lineComments[line]})
			lineCommentsVisited[line] = true
		}
	}
	return comments
}

func comment2String(comments []comment) []string {
	result := make([]string, len(comments))
	for i, comment := range comments {
		result[i] = comment.content
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
	commentsMap, err := ExtractComments(input)
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
