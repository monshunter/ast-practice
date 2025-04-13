package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

// MainPackageInfo 表示一个main包的信息
type MainPackageInfo struct {
	MainDir  string   `json:"mainDir"`
	MainFile string   `json:"mainFile"`
	Imports  []string `json:"imports"`
}

// 全局变量
var (
	projectRoot  string
	modulePrefix string
)

func main() {
	flag.Parse()

	// 获取项目根目录的绝对路径
	var err error
	projectRoot, err = filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "无法获取绝对路径: %v\n", err)
		os.Exit(1)
	}

	// 尝试读取go.mod获取模块前缀
	modulePrefix = getModulePrefix(projectRoot)
	if modulePrefix == "" {
		fmt.Fprintf(os.Stderr, "警告: 未找到go.mod文件，无法确定模块前缀\n")
	}
	analyzeMainPackages()
}

// getModulePrefix 从go.mod文件中获取模块前缀
func getModulePrefix(root string) string {
	modFilePath := filepath.Join(root, "go.mod")
	content, err := os.ReadFile(modFilePath)
	if err != nil {
		return ""
	}
	modFile, err := modfile.Parse(modFilePath, content, nil)
	if err != nil {
		return ""
	}
	fmt.Println(modFile.Module.Mod.Path)
	return modFile.Module.Mod.Path
}

// analyzeMainPackages 分析所有main包
func analyzeMainPackages() {
	// 找到所有的Go文件
	var goFiles []string
	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			base := filepath.Base(path)
			if base == "vendor" || base == "testdata" || base == ".git" ||
				base == "node_modules" || base == ".cursor" || base == ".vscode" {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			goFiles = append(goFiles, path)
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "遍历目录错误: %v\n", err)
		os.Exit(1)
	}

	// 找到所有main包
	mainPackages, err := findMainPackages(goFiles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "找到main包错误: %v\n", err)
		os.Exit(1)
	}

	// fmt.Println(mainPackages)

	// 分析每个main包
	results := make([]MainPackageInfo, 0, len(mainPackages))
	for _, mainDir := range mainPackages {
		info := analyzeMainImports(mainDir)
		results = append(results, info)
	}

	// 输出JSON结果
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON编码错误: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonData))
}

func getRelativeDir(root, dir string) string {
	relDir, err := filepath.Rel(root, dir)
	if err != nil {
		return ""
	}
	return relDir
}

// findMainPackages 找到所有的main包
func findMainPackages(goFiles []string) ([]string, error) {
	visited := make(map[string]bool)
	results := make([]string, 0, len(goFiles))
	for _, file := range goFiles {
		dir := filepath.Dir(file)
		if visited[dir] {
			continue
		}
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, file, nil, parser.PackageClauseOnly)
		if err != nil {
			continue
		}
		if node.Name.Name == "main" {
			results = append(results, dir)
			visited[dir] = true
		}
	}
	return results, nil
}

// analyzeMainPackage 分析单个main包
func analyzeMainImports(mainDir string) MainPackageInfo {
	info := MainPackageInfo{
		MainDir:  mainDir,
		MainFile: findMainEntryFile(mainDir),
	}
	// 分析导入的包
	importedPkgs := make(map[string]bool)
	collectImports(mainDir, importedPkgs)

	// 过滤出项目内部包
	var internalPackages []string
	for pkg := range importedPkgs {
		internalPackages = append(internalPackages, strings.Join([]string{projectRoot, getRelativeDir(modulePrefix, pkg)}, "/"))
	}
	info.Imports = internalPackages
	return info
}

func findMainEntryFile(mainDir string) string {
	entries, err := os.ReadDir(mainDir)
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".go") && strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, filepath.Join(mainDir, entry.Name()), nil, parser.ParseComments)
		if err != nil {
			continue
		}
		for _, imp := range node.Decls {
			if node, ok := imp.(*ast.FuncDecl); ok && node.Name.Name == "main" && node.Recv == nil {
				return strings.Join([]string{mainDir, entry.Name()}, "/")
			}
		}
	}
	return ""
}

// collectImports 收集目录中所有Go文件的导入
func collectImports(dir string, importedPkgs map[string]bool) {
	// 获取目录中的所有Go文件
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}

		filePath := filepath.Join(dir, name)
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, filePath, nil, parser.ImportsOnly)
		if err != nil {
			continue
		}

		// 添加导入的包
		for _, imp := range node.Imports {
			importPath := strings.Trim(imp.Path.Value, "\"")
			if importedPkgs[importPath] {
				continue
			}
			// 递归处理内部包
			if isInternalPackage(importPath) {
				importedPkgs[importPath] = true
				pkgDir := getPackageDir(importPath)
				if pkgDir != "" && pkgDir != dir {
					collectImports(pkgDir, importedPkgs)
				}
			}
		}
	}
}

// isInternalPackage 判断是否为项目内部包
func isInternalPackage(importPath string) bool {
	// 如果有模块前缀，检查是否以模块前缀开头
	if modulePrefix != "" && strings.HasPrefix(importPath, modulePrefix) {
		return true
	}

	// 没有模块前缀，则检查包是否在项目目录下
	pkgDir := getPackageDir(importPath)
	return pkgDir != "" && strings.HasPrefix(pkgDir, projectRoot)
}

// getPackageDir 获取包对应的目录
func getPackageDir(importPath string) string {
	if modulePrefix != "" && strings.HasPrefix(importPath, modulePrefix) {
		// 去掉模块前缀，获取相对路径
		relPath := strings.TrimPrefix(importPath, modulePrefix)
		relPath = strings.TrimPrefix(relPath, "/")
		return filepath.Join(projectRoot, relPath)
	}

	// 尝试在GOPATH中查找包
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		pkgDir := filepath.Join(gopath, "src", importPath)
		if _, err := os.Stat(pkgDir); err == nil {
			return pkgDir
		}
	}

	return ""
}
