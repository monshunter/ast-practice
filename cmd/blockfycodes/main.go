package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"slices"

	"golang.org/x/tools/go/ast/astutil"
)

func main() {
	content := initConfig()
	content, err := runInsertImport(content)
	if err != nil {
		log.Fatalf("Failed to insert import: %v", err)
	}
	content, err = runInsertStmt(content)
	if err != nil {
		log.Fatalf("Failed to insert expr: %v", err)
	}

	content, err = runInsertComment(content)
	if err != nil {
		log.Fatalf("Failed to insert comment: %v", err)
	}
	// ---
	fmt.Println("Modified Code:")
	fmt.Println("----------------")
	fmt.Println(string(content))
}

func initConfig() []byte {
	var content []byte
	// 检查参数
	if len(os.Args) != 2 {
		fmt.Println("用法: blockfycodes <文件路径或代码内容>")
		os.Exit(1)
	}

	input := os.Args[1]
	// 判断输入是文件路径还是代码内容
	if _, err := os.Stat(input); err == nil {
		// 输入是文件路径
		content, err = os.ReadFile(input)
		if err != nil {
			log.Fatalf("读取文件失败: %v\n", err)
		}
	} else {
		content = []byte(input)
	}
	return content
}

func getAstTree(content []byte) (*token.FileSet, *ast.File, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	return fset, f, err
}

func formatFile(fset *token.FileSet, f *ast.File) ([]byte, error) {
	var buf bytes.Buffer
	cfg := printer.Config{Mode: printer.UseSpaces, Tabwidth: 4}
	err := cfg.Fprint(&buf, fset, f)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func runInsertImport(content []byte) ([]byte, error) {
	fset, f, err := getAstTree(content)
	if err != nil {
		return nil, err
	}

	// 检查是否已存在"fmt" import
	hasFmt := false
	for _, ipt := range f.Imports {
		if ipt.Path.Value == `"fmt"` {
			hasFmt = true
			break
		}
	}
	if !hasFmt {
		astutil.AddNamedImport(fset, f, "f", "fmt")
	}
	return formatFile(fset, f)
}

func runInsertStmt(content []byte) ([]byte, error) {
	fset, f, err := getAstTree(content)
	if err != nil {
		return nil, err
	}

	positionsToInsert := []int{}
	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.FuncDecl); ok {
			extraStmt(decl.Body.List, fset, &positionsToInsert)
		}
	}
	printStmtStr := `fmt.Println("hello world")` + "\n"
	return doInsert(printStmtStr, fset, f, positionsToInsert)
}

func doInsert(printStmtStr string, fset *token.FileSet, f *ast.File, positionsToInsert []int) ([]byte, error) {
	// 创建需要插入的打印语句的AST节点
	// printStmtStr := `fmt.Println("hello world")` + "\n"

	// 将源代码转为字符串
	var src bytes.Buffer
	if err := printer.Fprint(&src, fset, f); err != nil {
		return nil, err
	}
	srcStr := src.String()
	slices.Sort(positionsToInsert)
	// 对于每个插入位置，将打印语句插入到源代码字符串中
	var buf bytes.Buffer
	buf.Grow(len(srcStr) + len(printStmtStr)*len(positionsToInsert))
	posIdx := 0
	lines := 0
	i, j := 0, 0
	for ; i < len(srcStr) && posIdx < len(positionsToInsert); i++ {
		if lines == positionsToInsert[posIdx]-1 {
			buf.WriteString(srcStr[j:i])
			buf.WriteString(printStmtStr)
			j = i
			posIdx++
		} else if srcStr[i] == '\n' {
			lines++
		}
	}
	buf.WriteString(srcStr[i:])
	newFset, newF, err := getAstTree(buf.Bytes())
	if err != nil {
		return nil, err
	}
	return formatFile(newFset, newF)
}

func extraStmt(statList []ast.Stmt, fset *token.FileSet, positionsToInsert *[]int) {
	// 遍历函数体中的语句
	for _, stmt := range statList {
		//可以根据语句类型进一步处理
		switch s := stmt.(type) {
		case *ast.AssignStmt:
			*positionsToInsert = append(*positionsToInsert, fset.Position(stmt.Pos()).Line)
			extraExpr(s.Rhs, fset, positionsToInsert)
		case *ast.IfStmt:
			extraStmt(s.Body.List, fset, positionsToInsert)
			if s.Else != nil {
				switch s.Else.(type) {
				case *ast.IfStmt:
					extraStmt([]ast.Stmt{s.Else.(*ast.IfStmt)}, fset, positionsToInsert)
				case *ast.BlockStmt:
					block := s.Else.(*ast.BlockStmt)
					extraStmt(block.List, fset, positionsToInsert)
					s.Else = block
				}
			}
		case *ast.ForStmt:
			*positionsToInsert = append(*positionsToInsert, fset.Position(stmt.Pos()).Line)
			extraStmt(s.Body.List, fset, positionsToInsert)
		case *ast.RangeStmt:
			*positionsToInsert = append(*positionsToInsert, fset.Position(stmt.Pos()).Line)
			extraStmt(s.Body.List, fset, positionsToInsert)
		case *ast.SwitchStmt:
			*positionsToInsert = append(*positionsToInsert, fset.Position(stmt.Pos()).Line)
			extraStmt(s.Body.List, fset, positionsToInsert)
		case *ast.SelectStmt:
			*positionsToInsert = append(*positionsToInsert, fset.Position(stmt.Pos()).Line)
			extraStmt(s.Body.List, fset, positionsToInsert)
		case *ast.TypeSwitchStmt:
			*positionsToInsert = append(*positionsToInsert, fset.Position(stmt.Pos()).Line)
			extraStmt(s.Body.List, fset, positionsToInsert)
		case *ast.CommClause:
			extraStmt(s.Body, fset, positionsToInsert)
		case *ast.CaseClause:
			extraStmt(s.Body, fset, positionsToInsert)
		case *ast.BlockStmt:
			*positionsToInsert = append(*positionsToInsert, fset.Position(stmt.Pos()).Line)
			extraStmt(s.List, fset, positionsToInsert)

		case *ast.ReturnStmt:
			*positionsToInsert = append(*positionsToInsert, fset.Position(stmt.Pos()).Line)
			for _, result := range s.Results {
				extraExpr([]ast.Expr{result}, fset, positionsToInsert)
			}
		case *ast.DeferStmt:
			*positionsToInsert = append(*positionsToInsert, fset.Position(stmt.Pos()).Line)
			if s.Call != nil && s.Call.Fun != nil {
				extraExpr([]ast.Expr{s.Call.Fun}, fset, positionsToInsert)
			}

		case *ast.GoStmt:
			*positionsToInsert = append(*positionsToInsert, fset.Position(stmt.Pos()).Line)
			if s.Call != nil && s.Call.Fun != nil {
				extraExpr([]ast.Expr{s.Call.Fun}, fset, positionsToInsert)
			}
		case *ast.ExprStmt:
			switch s.X.(type) {
			case *ast.CallExpr:
				*positionsToInsert = append(*positionsToInsert, fset.Position(stmt.Pos()).Line)
				expr := s.X.(*ast.CallExpr)
				if expr.Fun != nil {
					extraExpr([]ast.Expr{expr.Fun}, fset, positionsToInsert)
				}
			default:
				*positionsToInsert = append(*positionsToInsert, fset.Position(stmt.Pos()).Line)
			}
		default:
			*positionsToInsert = append(*positionsToInsert, fset.Position(stmt.Pos()).Line)
		}
	}
}

func extraExpr(exprList []ast.Expr, fset *token.FileSet, positionsToInsert *[]int) {
	for _, expr := range exprList {
		switch expr := expr.(type) {
		case *ast.FuncLit:
			extraStmt(expr.Body.List, fset, positionsToInsert)
		}
	}
}

func runInsertComment(content []byte) ([]byte, error) {
	fset, f, err := getAstTree(content)
	if err != nil {
		return nil, err
	}
	positionsToInsert := []int{}
	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.FuncDecl); ok {
			positionsToInsert = append(positionsToInsert, fset.Position(decl.Pos()).Line)
			extraStmtAndInsertComment(f, decl.Body.List, fset, &positionsToInsert)
		}
	}
	printStmtStr := `// + insert comment` + "\n"
	return doInsert(printStmtStr, fset, f, positionsToInsert)
}

func extraStmtAndInsertComment(f *ast.File, statList []ast.Stmt, fset *token.FileSet, positionsToInsert *[]int) {
	// 遍历函数体中的语句
	for _, stmt := range statList {
		if _, ok := stmt.(*ast.IfStmt); !ok {
			*positionsToInsert = append(*positionsToInsert, fset.Position(stmt.Pos()-1).Line)
		}
		//可以根据语句类型进一步处理
		switch s := stmt.(type) {
		case *ast.AssignStmt:
			extraExprAndInsertComment(f, s.Rhs, fset, positionsToInsert)
		case *ast.IfStmt:
			*positionsToInsert = append(*positionsToInsert, fset.Position(s.Body.Lbrace).Line)
			extraStmtAndInsertComment(f, s.Body.List, fset, positionsToInsert)
			if s.Else != nil {
				switch s.Else.(type) {
				case *ast.IfStmt:
					extraStmtAndInsertComment(f, []ast.Stmt{s.Else.(*ast.IfStmt)}, fset, positionsToInsert)
				case *ast.BlockStmt:
					*positionsToInsert = append(*positionsToInsert, fset.Position(s.Else.Pos()).Line)
					block := s.Else.(*ast.BlockStmt)
					extraStmtAndInsertComment(f, block.List, fset, positionsToInsert)
				}
			}
		case *ast.ForStmt:
			extraStmtAndInsertComment(f, s.Body.List, fset, positionsToInsert)
		case *ast.RangeStmt:
			extraStmtAndInsertComment(f, s.Body.List, fset, positionsToInsert)
		case *ast.SwitchStmt:
			extraStmtAndInsertComment(f, s.Body.List, fset, positionsToInsert)
		case *ast.SelectStmt:
			extraStmtAndInsertComment(f, s.Body.List, fset, positionsToInsert)
		case *ast.TypeSwitchStmt:
			extraStmtAndInsertComment(f, s.Body.List, fset, positionsToInsert)
		case *ast.CommClause:
			extraStmtAndInsertComment(f, s.Body, fset, positionsToInsert)
		case *ast.CaseClause:
			extraStmtAndInsertComment(f, s.Body, fset, positionsToInsert)
		case *ast.BlockStmt:
			extraStmtAndInsertComment(f, s.List, fset, positionsToInsert)
		case *ast.ReturnStmt:
			for _, result := range s.Results {
				extraExprAndInsertComment(f, []ast.Expr{result}, fset, positionsToInsert)
			}
		case *ast.DeferStmt:
			if s.Call != nil && s.Call.Fun != nil {
				extraExprAndInsertComment(f, []ast.Expr{s.Call.Fun}, fset, positionsToInsert)
			}

		case *ast.GoStmt:
			if s.Call != nil && s.Call.Fun != nil {
				extraExprAndInsertComment(f, []ast.Expr{s.Call.Fun}, fset, positionsToInsert)
			}

		case *ast.ExprStmt:
			switch s.X.(type) {
			case *ast.CallExpr:
				expr := s.X.(*ast.CallExpr)
				if expr.Fun != nil {
					extraExprAndInsertComment(f, []ast.Expr{expr.Fun}, fset, positionsToInsert)
				}
			}
		}
	}
}

func extraExprAndInsertComment(f *ast.File, exprList []ast.Expr, fset *token.FileSet, positionsToInsert *[]int) {
	for _, expr := range exprList {
		switch expr := expr.(type) {
		case *ast.FuncLit:
			extraStmtAndInsertComment(f, expr.Body.List, fset, positionsToInsert)
		}
	}
}
