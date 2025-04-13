package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// MainEntry 表示一个main入口点
type MainEntry struct {
	Dir  string `json:"dir"`  // 相对于项目的目录
	File string `json:"file"` // 文件名
}

// 命令行参数
var (
	verbose    bool
	outputFile string
	help       bool
)

func init() {
	flag.BoolVar(&verbose, "v", false, "启用详细输出模式")
	flag.StringVar(&outputFile, "o", "", "输出文件路径（默认为标准输出）")
	flag.BoolVar(&help, "h", false, "显示帮助信息")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "用法: %s [选项] <项目路径>\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "选项:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n示例:\n")
	fmt.Fprintf(os.Stderr, "  %s /path/to/go/project\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -v -o result.json /path/to/go/project\n", os.Args[0])
}

func main() {
	flag.Parse()

	if help {
		usage()
		os.Exit(0)
	}

	// 检查命令行参数
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "错误: 必须提供项目路径")
		usage()
		os.Exit(1)
	}

	projectPath := args[0]

	// 检查项目路径是否存在
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "错误: 项目路径不存在: %s\n", projectPath)
		os.Exit(1)
	}

	// 获取绝对路径，用于后续计算相对路径
	absProjectPath, err := filepath.Abs(projectPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: 获取绝对路径失败: %v\n", err)
		os.Exit(1)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "分析项目: %s\n", absProjectPath)
	}

	// 查找所有main入口
	entries, err := findMainEntries(absProjectPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: 查找main入口失败: %v\n", err)
		os.Exit(1)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "找到 %d 个main入口\n", len(entries))
	}

	// 输出JSON格式结果
	result, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: JSON序列化失败: %v\n", err)
		os.Exit(1)
	}

	// 将结果写入输出
	var out io.Writer = os.Stdout
	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: 无法创建输出文件: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		out = file
		if verbose {
			fmt.Fprintf(os.Stderr, "结果已写入: %s\n", outputFile)
		}
	}

	fmt.Fprintln(out, string(result))
}

// findMainEntries 查找项目中所有的main入口
func findMainEntries(projectPath string) ([]MainEntry, error) {
	var entries []MainEntry

	// 遍历目录，查找所有Go文件
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// 跳过无法访问的目录/文件
			if verbose {
				fmt.Fprintf(os.Stderr, "警告: 跳过 %s: %v\n", path, err)
			}
			return nil
		}

		// 跳过目录和非Go文件
		if info.IsDir() {
			// 跳过常见的非Go代码目录
			base := filepath.Base(path)
			if base == "vendor" || base == ".git" || base == "node_modules" || base == "testdata" {
				if verbose {
					fmt.Fprintf(os.Stderr, "跳过目录: %s\n", path)
				}
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(info.Name(), ".go") || strings.HasSuffix(info.Name(), "_test.go") {
			return nil
		}

		// 检查文件是否包含main包和main函数
		isMainEntry, err := isMainEntryFile(path)
		if err != nil {
			if verbose {
				fmt.Fprintf(os.Stderr, "警告: 无法解析 %s: %v\n", path, err)
			}
			return nil
		}

		if isMainEntry {
			if verbose {
				fmt.Fprintf(os.Stderr, "找到main入口: %s\n", path)
			}

			// 计算相对于项目的目录
			relDir, err := filepath.Rel(projectPath, filepath.Dir(path))
			if err != nil {
				return err
			}

			// 确保路径分隔符统一（使用斜杠，符合JSON输出习惯）
			relDir = filepath.ToSlash(relDir)
			if relDir == "." {
				relDir = ""
			}

			// 添加到结果
			entries = append(entries, MainEntry{
				Dir:  relDir,
				File: info.Name(),
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return entries, nil
}

// isMainEntryFile 检查文件是否是main入口（包含package main和func main）
func isMainEntryFile(filePath string) (bool, error) {
	// 创建文件集合
	fset := token.NewFileSet()

	// 解析Go源文件
	file, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return false, err
	}

	// 检查是否是main包
	if file.Name.Name != "main" {
		return false, nil
	}

	// 检查是否包含main函数
	hasMainFunc := false
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if fn.Name.Name == "main" && fn.Recv == nil {
				hasMainFunc = true
				break
			}
		}
	}

	return hasMainFunc, nil
}
