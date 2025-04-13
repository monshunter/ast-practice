package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsMainEntryFile(t *testing.T) {
	// 创建临时测试目录
	tempDir, err := os.MkdirTemp("", "findmain-test")
	if err != nil {
		t.Fatalf("无法创建临时目录: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 测试用例1: 包含main包和main函数的文件
	validMainContent := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`
	validMainPath := filepath.Join(tempDir, "valid_main.go")
	if err := os.WriteFile(validMainPath, []byte(validMainContent), 0644); err != nil {
		t.Fatalf("无法写入测试文件: %v", err)
	}

	// 测试用例2: 包含main包但没有main函数的文件
	noMainFuncContent := `package main

import "fmt"

func otherFunc() {
	fmt.Println("Not a main function")
}
`
	noMainFuncPath := filepath.Join(tempDir, "no_main_func.go")
	if err := os.WriteFile(noMainFuncPath, []byte(noMainFuncContent), 0644); err != nil {
		t.Fatalf("无法写入测试文件: %v", err)
	}

	// 测试用例3: 不是main包的文件
	notMainPackageContent := `package other

import "fmt"

func main() {
	fmt.Println("Not a main package")
}
`
	notMainPackagePath := filepath.Join(tempDir, "not_main_package.go")
	if err := os.WriteFile(notMainPackagePath, []byte(notMainPackageContent), 0644); err != nil {
		t.Fatalf("无法写入测试文件: %v", err)
	}

	// 运行测试
	testCases := []struct {
		name     string
		path     string
		expected bool
	}{
		{"有效的main入口", validMainPath, true},
		{"没有main函数", noMainFuncPath, false},
		{"不是main包", notMainPackagePath, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := isMainEntryFile(tc.path)
			if err != nil {
				t.Fatalf("isMainEntryFile出错: %v", err)
			}
			if result != tc.expected {
				t.Errorf("isMainEntryFile(%s) = %v, 期望 %v", tc.path, result, tc.expected)
			}
		})
	}
}

func TestFindMainEntries(t *testing.T) {
	// 创建临时项目结构用于测试
	tempDir, err := os.MkdirTemp("", "findmain-project-test")
	if err != nil {
		t.Fatalf("无法创建临时目录: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建项目结构
	dirs := []string{
		filepath.Join(tempDir, "cmd", "app1"),
		filepath.Join(tempDir, "cmd", "app2"),
		filepath.Join(tempDir, "internal", "pkg"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("无法创建目录 %s: %v", dir, err)
		}
	}

	// 创建测试文件
	files := map[string]string{
		// 有效的main入口1
		filepath.Join(tempDir, "cmd", "app1", "main.go"): `package main

func main() {}`,
		// 有效的main入口2
		filepath.Join(tempDir, "cmd", "app2", "app.go"): `package main

func main() {}`,
		// 非main包
		filepath.Join(tempDir, "internal", "pkg", "util.go"): `package pkg

func Util() {}`,
		// main包但无main函数
		filepath.Join(tempDir, "cmd", "app1", "helper.go"): `package main

func helper() {}`,
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("无法写入文件 %s: %v", path, err)
		}
	}

	// 运行测试
	entries, err := findMainEntries(tempDir)
	if err != nil {
		t.Fatalf("findMainEntries出错: %v", err)
	}

	// 验证结果
	if len(entries) != 2 {
		t.Errorf("找到 %d 个入口点, 期望 2", len(entries))
	}

	// 验证特定入口点
	expectedEntries := map[string]string{
		"cmd/app1": "main.go",
		"cmd/app2": "app.go",
	}

	for _, entry := range entries {
		expectedFile, ok := expectedEntries[entry.Dir]
		if !ok {
			t.Errorf("未预期的入口点目录: %s", entry.Dir)
			continue
		}
		if entry.File != expectedFile {
			t.Errorf("目录 %s 中的文件名错误, 得到 %s, 期望 %s", entry.Dir, entry.File, expectedFile)
		}
		delete(expectedEntries, entry.Dir)
	}

	if len(expectedEntries) > 0 {
		for dir, file := range expectedEntries {
			t.Errorf("未找到预期的入口点: %s/%s", dir, file)
		}
	}
}
