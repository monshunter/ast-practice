package getcomments

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func TestExampleGetComments(t *testing.T) {
	// 获取示例文件的路径
	exampleFile := filepath.Join("testdata", "examples.go")

	// 直接调用ExtractComments函数
	commentsMap, err := ExtractComments(exampleFile)
	if err != nil {
		t.Fatalf("提取注释失败: %v", err)
	}

	// 输出结果以便调试
	output, err := json.MarshalIndent(commentsMap, "", "    ")
	if err != nil {
		t.Fatalf("序列化结果失败: %v", err)
	}

	fmt.Printf("获取到的注释结果:\n%s\n", output)

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

	// 检查是否包含所有预期的键
	for key := range expectedResults {
		if _, ok := commentsMap[key]; !ok {
			t.Errorf("未找到预期的键: %s", key)
		}
	}

	// 检查每个键的注释内容
	for key, expectedComments := range expectedResults {
		if comments, ok := commentsMap[key]; ok {
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
