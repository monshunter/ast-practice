package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: getblockscopes <filename>")
		os.Exit(1)
	}
	filename := os.Args[1]
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	blockScopes, err := BlockScopesOfGoAST(filename, content)
	if err != nil {
		fmt.Println("error:", err)
	}
	sort.Slice(blockScopes, func(i, j int) bool {
		if blockScopes[i].StartLine == blockScopes[j].StartLine {
			return blockScopes[i].EndLine < blockScopes[j].EndLine
		}
		return blockScopes[i].StartLine < blockScopes[j].StartLine
	})
	fmt.Println(blockScopes)
	lines := strings.Split(string(content), "\n")
	for i := range len(lines) {
		idx := blockScopes.Search(i + 1)
		fmt.Printf("%d: %s\n", i+1, blockScopes[idx].String())
	}
}

type BlockScope struct {
	StartLine int
	EndLine   int
}

func (b *BlockScope) String() string {
	return fmt.Sprintf("BlockScope{StartLine: %d, EndLine: %d}", b.StartLine, b.EndLine)
}

func (b *BlockScope) IsEmpty() bool {
	return b.StartLine == 0 && b.EndLine == 0
}

func (b *BlockScope) IsValid() bool {
	return b.StartLine < b.EndLine
}

func (b *BlockScope) Contains(line int) bool {
	return line > b.StartLine && line < b.EndLine
}

func (b *BlockScope) ContainsRange(start, end int) bool {
	return start > b.StartLine && end < b.EndLine
}

type BlockScopes []BlockScope

func (b BlockScopes) Search(line int) int {
	// find lastest scope of the line
	l, r := 0, len(b)-1
	idx := 0
	for l <= r {
		mid := l + (r-l)/2
		if b[mid].StartLine <= line {
			if b[mid].EndLine >= line {
				idx = mid
			}
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	return idx
}

func BlockScopesOfGoAST(filename string, content []byte) (BlockScopes, error) {

	fset := token.NewFileSet()
	blockScopes := []BlockScope{}
	astFile, err := parser.ParseFile(fset, filename, content, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	blockScopes = append(blockScopes, BlockScope{
		StartLine: 1,
		EndLine:   len(lines),
	})
	for _, decl := range astFile.Decls {
		if declFunc, ok := decl.(*ast.FuncDecl); ok {
			blockScopes = append(blockScopes, BlockScope{
				StartLine: fset.Position(declFunc.Pos()).Line,
				EndLine:   fset.Position(declFunc.End()).Line,
			})
			for _, stmt := range declFunc.Body.List {
				ast.Inspect(stmt, func(node ast.Node) bool {
					if node == nil {
						return false
					}
					switch stmt := node.(type) {
					case *ast.IfStmt:
						blockScopes = append(blockScopes, BlockScope{
							StartLine: fset.Position(stmt.Pos()).Line,
							EndLine:   fset.Position(stmt.End()).Line,
						})
						if stmt.Else != nil {
							switch stmt.Else.(type) {
							case *ast.BlockStmt:
								block := stmt.Else.(*ast.BlockStmt)
								blockScopes = append(blockScopes, BlockScope{
									StartLine: fset.Position(block.Pos()).Line,
									EndLine:   fset.Position(block.End()).Line,
								})
							}
						}
					case *ast.ForStmt:
						blockScopes = append(blockScopes, BlockScope{
							StartLine: fset.Position(stmt.Pos()).Line,
							EndLine:   fset.Position(stmt.End()).Line,
						})
					case *ast.RangeStmt:
						blockScopes = append(blockScopes, BlockScope{
							StartLine: fset.Position(stmt.Pos()).Line,
							EndLine:   fset.Position(stmt.End()).Line,
						})
					case *ast.CaseClause:
						if stmt.Body != nil {
							blockScopes = append(blockScopes, BlockScope{
								StartLine: fset.Position(stmt.Body[0].Pos()).Line - 1,
								EndLine:   fset.Position(stmt.Body[len(stmt.Body)-1].End()).Line + 1,
							})
						}

					case *ast.CommClause:
						if stmt.Body != nil {
							blockScopes = append(blockScopes, BlockScope{
								StartLine: fset.Position(stmt.Body[0].Pos()).Line - 1,
								EndLine:   fset.Position(stmt.Body[len(stmt.Body)-1].End()).Line + 1,
							})
						}
					case *ast.FuncLit:
						blockScopes = append(blockScopes, BlockScope{
							StartLine: fset.Position(stmt.Pos()).Line,
							EndLine:   fset.Position(stmt.End()).Line,
						})
					}
					return true
				})
			}
		}
	}
	return blockScopes, nil
}
