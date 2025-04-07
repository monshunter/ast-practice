package examples

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestExampleGetComments(t *testing.T) {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取工作目录失败: %v", err)
	}

	// 确定当前示例文件的路径
	exampleFile, err := filepath.Abs(filepath.Join(wd, "examples.go"))
	if err != nil {
		t.Fatalf("获取绝对路径失败: %v", err)
	}

	// 构建getcomments的绝对路径
	parentDir := filepath.Dir(wd)
	getcommentsPath := filepath.Join(parentDir, "getcomments")

	// 检查可执行文件是否存在
	if _, err := os.Stat(getcommentsPath); os.IsNotExist(err) {
		// 如果可执行文件不存在，则尝试用go run运行
		cmd := exec.Command("go", "run", filepath.Join(parentDir, "main.go"), exampleFile)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("运行getcomments失败: %v\n输出: %s", err, output)
		}

		processOutput(t, output)
	} else {
		// 直接运行可执行文件
		cmd := exec.Command(getcommentsPath, exampleFile)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("运行getcomments失败: %v\n输出: %s", err, output)
		}

		processOutput(t, output)
	}
}

func processOutput(t *testing.T, output []byte) {
	// 解析输出的JSON
	var result map[string][]string
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("解析输出JSON失败: %v\n输出: %s", err, output)
	}

	// 期望的输出结果，按行号和包含的注释内容
	expectedResults := map[string][]string{
		"examples.go:11": {
			"ExampleGetComments 示例",
			"输入:",
			"一段golang代码",
			"输出:",
			"文件名:行号: []string{comments...}",
			"函数名: 函数注释 ExampleGetComments",
		},
		"examples.go:13": {"定义两个变量"},
		"examples.go:15": {"打印变量"},
		"examples.go:17": {"判断变量大小", "if 分支"},
		"examples.go:20": {"否则进入 else 分支", "else 分支"},
		"examples.go:22": {"打印输出： x is less than y", "行尾注释打印输出： x is less than y"},
	}

	// 打印结果，便于调试
	fmt.Printf("获取到的注释结果:\n%s\n", output)

	// 检查是否包含所有预期的键
	for key := range expectedResults {
		if _, ok := result[key]; !ok {
			t.Errorf("未找到预期的键: %s", key)
		}
	}

	// 检查每个键的注释内容
	for key, expectedComments := range expectedResults {
		if comments, ok := result[key]; ok {
			for _, expectedContent := range expectedComments {
				found := false
				for _, actualComment := range comments {
					// 去除注释符号后比较内容
					actualContent := strings.TrimPrefix(actualComment, "//")
					actualContent = strings.TrimSpace(actualContent)
					if strings.Contains(actualContent, expectedContent) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("键 %s 的注释中应该包含 '%s'，但未找到", key, expectedContent)
				}
			}
		} else {
			t.Errorf("结果中没有键 %s", key)
		}
	}
}
