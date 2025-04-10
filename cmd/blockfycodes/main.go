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
	"path/filepath"

	"golang.org/x/tools/go/ast/astutil"
)

func main() {
	filename, content := initConfig()
	buf, err := runInsertImport(filename, content)
	if err != nil {
		log.Fatalf("Failed to insert import: %v", err)
	}
	buf, err = runInsertStmt(filename, buf)
	if err != nil {
		log.Fatalf("Failed to insert expr: %v", err)
	}

	buf, err = runInsertComment(filename, buf)
	if err != nil {
		log.Fatalf("Failed to insert comment: %v", err)
	}
	// ---
	fmt.Println("Modified Code:")
	fmt.Println("----------------")
	fmt.Println(string(buf))
}

func initConfig() (string, []byte) {
	var content []byte
	var filename string
	// 检查参数
	if len(os.Args) != 2 {
		fmt.Println("用法: blockfycodes <文件路径或代码内容>")
		os.Exit(1)
	}

	input := os.Args[1]
	// 判断输入是文件路径还是代码内容
	if _, err := os.Stat(input); err == nil {
		// 输入是文件路径
		filename = filepath.Base(input)
		content, err = os.ReadFile(input)
		if err != nil {
			log.Fatalf("读取文件失败: %v\n", err)
		}
	} else {
		// 输入是代码内容
		filename = "testdata/example.go"
		content = []byte(input)
	}
	return filename, content
}

func runInsertImport(filename string, content []byte) ([]byte, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, content, parser.ParseComments)
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

	var buf bytes.Buffer
	cfg := printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 8}
	err = cfg.Fprint(&buf, fset, f)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func runInsertStmt(filename string, content []byte) ([]byte, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, content, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.FuncDecl); ok {
			decl.Body.List = extraStmt(decl.Body.List, fset)
		}
	}
	var buf bytes.Buffer
	cfg := printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 8}
	err = cfg.Fprint(&buf, fset, f)
	if err != nil {
		log.Fatalf("Failed to print AST: %v", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func extraStmt(statList []ast.Stmt, fset *token.FileSet) []ast.Stmt {
	// 遍历函数体中的语句
	newStatList := make([]ast.Stmt, 0, len(statList))
	for _, stmt := range statList {
		//可以根据语句类型进一步处理
		switch s := stmt.(type) {
		case *ast.AssignStmt:
			newStatList = append(newStatList, createStatementNode())
			s.Rhs = extraExpr(s.Rhs, fset)
		case *ast.IfStmt:
			s.Body.List = extraStmt(s.Body.List, fset)
			if s.Else != nil {
				switch s.Else.(type) {
				case *ast.IfStmt:
					s.Else = extraStmt([]ast.Stmt{s.Else.(*ast.IfStmt)}, fset)[0]
				case *ast.BlockStmt:
					block := s.Else.(*ast.BlockStmt)
					block.List = extraStmt(block.List, fset)
					s.Else = block
				}
			}
		case *ast.ForStmt:
			newStatList = append(newStatList, createStatementNode())
			s.Body.List = extraStmt(s.Body.List, fset)
		case *ast.RangeStmt:
			newStatList = append(newStatList, createStatementNode())
			s.Body.List = extraStmt(s.Body.List, fset)
		case *ast.SwitchStmt:
			newStatList = append(newStatList, createStatementNode())
			s.Body.List = extraStmt(s.Body.List, fset)
		case *ast.CommClause:
			s.Body = extraStmt(s.Body, fset)
		case *ast.CaseClause:
			s.Body = extraStmt(s.Body, fset)
		case *ast.BlockStmt:
			newStatList = append(newStatList, createStatementNode())
			s.List = extraStmt(s.List, fset)
		case *ast.ReturnStmt:
			newStatList = append(newStatList, createStatementNode())
			for i, result := range s.Results {
				s.Results[i] = extraExpr([]ast.Expr{result}, fset)[0]
			}
		case *ast.DeferStmt:
			newStatList = append(newStatList, createStatementNode())
			if s.Call != nil && s.Call.Fun != nil {
				s.Call.Fun = extraExpr([]ast.Expr{s.Call.Fun}, fset)[0]
			}
		case *ast.SelectStmt:
			newStatList = append(newStatList, createStatementNode())
			s.Body.List = extraStmt(s.Body.List, fset)
		case *ast.GoStmt:
			newStatList = append(newStatList, createStatementNode())
			if s.Call != nil && s.Call.Fun != nil {
				s.Call.Fun = extraExpr([]ast.Expr{s.Call.Fun}, fset)[0]
			}
		case *ast.TypeSwitchStmt:
			newStatList = append(newStatList, createStatementNode())
			s.Body.List = extraStmt(s.Body.List, fset)
		case *ast.ExprStmt:
			switch s.X.(type) {
			case *ast.CallExpr:
				newStatList = append(newStatList, createStatementNode())
				expr := s.X.(*ast.CallExpr)
				if expr.Fun != nil {
					expr.Fun = extraExpr([]ast.Expr{expr.Fun}, fset)[0]
				}
			default:
				newStatList = append(newStatList, createStatementNode())
			}
		default:
			newStatList = append(newStatList, createStatementNode())
		}
		newStatList = append(newStatList, stmt)
	}
	return newStatList
}

func extraExpr(exprList []ast.Expr, fset *token.FileSet) []ast.Expr {
	newExprList := make([]ast.Expr, 0, len(exprList))
	for _, expr := range exprList {
		switch expr := expr.(type) {
		case *ast.FuncLit:
			expr.Body.List = extraStmt(expr.Body.List, fset)
		}
		newExprList = append(newExprList, expr)
	}
	return newExprList
}

func createStatementNode() ast.Stmt {
	return &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("fmt"),
				Sel: ast.NewIdent("Println"),
			},
			Args: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"hello world"`,
				},
			},
		},
	}
}

func runInsertComment(filename string, content []byte) ([]byte, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, content, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(f.Name.Name)
	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.FuncDecl); ok {
			insertComment(f, decl.Pos()-1)
			decl.Body.List = extraStmtAndInsertComment(f, decl.Body.List, fset)
		}
	}
	var buf bytes.Buffer
	cfg := printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 8}
	err = cfg.Fprint(&buf, fset, f)
	if err != nil {
		log.Fatalf("Failed to print AST: %v", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func extraStmtAndInsertComment(f *ast.File, statList []ast.Stmt, fset *token.FileSet) []ast.Stmt {
	// 遍历函数体中的语句
	newStatList := make([]ast.Stmt, 0, len(statList))
	for _, stmt := range statList {
		if _, ok := stmt.(*ast.IfStmt); !ok {
			insertComment(f, stmt.Pos()-1)
		}
		// insertComment(f, stmt.Pos()-1)
		//可以根据语句类型进一步处理
		switch s := stmt.(type) {
		case *ast.AssignStmt:
			s.Rhs = extraExprAndInsertComment(f, s.Rhs, fset)
		case *ast.IfStmt:
			insertComment(f, s.Body.Lbrace)
			s.Body.List = extraStmtAndInsertComment(f, s.Body.List, fset)
			if s.Else != nil {
				switch s.Else.(type) {
				case *ast.IfStmt:
					s.Else = extraStmtAndInsertComment(f, []ast.Stmt{s.Else.(*ast.IfStmt)}, fset)[0]
				case *ast.BlockStmt:
					insertComment(f, s.Else.Pos())
					block := s.Else.(*ast.BlockStmt)
					block.List = extraStmtAndInsertComment(f, block.List, fset)

					s.Else = block
				}
			}
		case *ast.ForStmt:
			s.Body.List = extraStmtAndInsertComment(f, s.Body.List, fset)
		case *ast.RangeStmt:
			s.Body.List = extraStmtAndInsertComment(f, s.Body.List, fset)
		case *ast.SwitchStmt:
			s.Body.List = extraStmtAndInsertComment(f, s.Body.List, fset)
		case *ast.CommClause:
			s.Body = extraStmtAndInsertComment(f, s.Body, fset)
		case *ast.CaseClause:
			s.Body = extraStmtAndInsertComment(f, s.Body, fset)
		case *ast.BlockStmt:
			s.List = extraStmtAndInsertComment(f, s.List, fset)
		case *ast.ReturnStmt:
			for i, result := range s.Results {
				s.Results[i] = extraExprAndInsertComment(f, []ast.Expr{result}, fset)[0]
			}
		case *ast.DeferStmt:
			if s.Call != nil && s.Call.Fun != nil {
				s.Call.Fun = extraExprAndInsertComment(f, []ast.Expr{s.Call.Fun}, fset)[0]
			}
		case *ast.SelectStmt:
			s.Body.List = extraStmtAndInsertComment(f, s.Body.List, fset)
		case *ast.GoStmt:
			if s.Call != nil && s.Call.Fun != nil {
				s.Call.Fun = extraExprAndInsertComment(f, []ast.Expr{s.Call.Fun}, fset)[0]
			}
		case *ast.TypeSwitchStmt:
			s.Body.List = extraStmtAndInsertComment(f, s.Body.List, fset)
		case *ast.ExprStmt:
			switch s.X.(type) {
			case *ast.CallExpr:
				expr := s.X.(*ast.CallExpr)
				if expr.Fun != nil {
					expr.Fun = extraExprAndInsertComment(f, []ast.Expr{expr.Fun}, fset)[0]
				}
			}
		}
		newStatList = append(newStatList, stmt)
	}
	return newStatList
}

func extraExprAndInsertComment(f *ast.File, exprList []ast.Expr, fset *token.FileSet) []ast.Expr {
	newExprList := make([]ast.Expr, 0, len(exprList))
	for _, expr := range exprList {
		switch expr := expr.(type) {
		case *ast.FuncLit:
			expr.Body.List = extraStmtAndInsertComment(f, expr.Body.List, fset)
		}
		newExprList = append(newExprList, expr)
	}
	return newExprList
}

func insertComment(f *ast.File, pos token.Pos) {
	f.Comments = append(f.Comments, &ast.CommentGroup{
		List: []*ast.Comment{
			{
				Slash: pos,
				Text:  "// + insert comment test",
			},
		},
	})
}
