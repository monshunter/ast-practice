package getcomments

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// 测试文件生成函数
func createTestFile(b *testing.B, fileName string, lineCount int, commentRatio float64) string {
	// 确保testdata目录存在
	dir := "testdata"
	if err := os.MkdirAll(dir, 0755); err != nil {
		b.Fatalf("创建测试目录失败: %v", err)
	}

	// 创建文件
	filePath := filepath.Join(dir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		b.Fatalf("创建文件失败: %v", err)
	}
	defer file.Close()

	// 生成测试文件内容
	generateTestFileContent(file, lineCount, commentRatio)

	return filePath
}

// 生成测试文件内容
func generateTestFileContent(file *os.File, lineCount int, commentRatio float64) {
	// 写入包声明和导入
	file.WriteString("package testdata\n\n")
	file.WriteString("import (\n")
	file.WriteString("\t\"context\"\n")
	file.WriteString("\t\"errors\"\n")
	file.WriteString("\t\"fmt\"\n")
	file.WriteString("\t\"bytes\"\n")
	file.WriteString("\t\"time\"\n")
	file.WriteString(")\n\n")

	// 写入一些常量
	file.WriteString("const (\n")
	for i := 0; i < 5; i++ {
		file.WriteString(fmt.Sprintf("\tConst%d = %d // 常量值%d\n", i, i*100, i))
	}
	file.WriteString(")\n\n")

	// 写入一些变量
	file.WriteString("var (\n")
	for i := 0; i < 5; i++ {
		file.WriteString(fmt.Sprintf("\tVar%d = \"值%d\" // 变量%d\n", i, i, i))
	}
	file.WriteString(")\n\n")

	// 写入一些结构体
	for i := 0; i < 3; i++ {
		file.WriteString(fmt.Sprintf("// TestStruct%d 是测试结构体\n", i))
		file.WriteString(fmt.Sprintf("// 用于测试注释提取功能\n"))
		file.WriteString(fmt.Sprintf("type TestStruct%d struct {\n", i))
		file.WriteString("\tName string // 名称\n")
		file.WriteString("\tValue int // 值\n")
		file.WriteString("\tEnabled bool // 是否启用\n")
		file.WriteString("}\n\n")
	}

	// 写入一些函数
	for i := 0; i < 5; i++ {
		file.WriteString(fmt.Sprintf("// TestFunc%d 是测试函数\n", i))
		file.WriteString("// 该函数有以下特点:\n")
		file.WriteString("// 1. 有多行注释\n")
		file.WriteString("// 2. 包含参数验证\n")
		file.WriteString("// 3. 有返回值处理\n")
		file.WriteString(fmt.Sprintf("func TestFunc%d(ctx context.Context, input string) (string, error) {\n", i))
		file.WriteString("\t// 验证上下文\n")
		file.WriteString("\tif ctx == nil {\n")
		if rand.Float64() < commentRatio {
			file.WriteString("\t\treturn \"\", errors.New(\"上下文不能为空\") // 参数验证\n")
		} else {
			file.WriteString("\t\treturn \"\", errors.New(\"上下文不能为空\")\n")
		}
		file.WriteString("\t}\n\n")

		file.WriteString("\t// 验证输入\n")
		file.WriteString("\tif input == \"\" {\n")
		if rand.Float64() < commentRatio {
			file.WriteString("\t\treturn \"\", errors.New(\"输入不能为空\") // 参数验证\n")
		} else {
			file.WriteString("\t\treturn \"\", errors.New(\"输入不能为空\")\n")
		}
		file.WriteString("\t}\n\n")

		file.WriteString("\t// 处理逻辑\n")
		file.WriteString("\tresult := fmt.Sprintf(\"处理结果: %s\", input)\n")
		if rand.Float64() < commentRatio {
			file.WriteString("\treturn result, nil // 返回结果\n")
		} else {
			file.WriteString("\treturn result, nil\n")
		}
		file.WriteString("}\n\n")
	}

	// 写入一些方法
	for i := 0; i < 3; i++ {
		file.WriteString(fmt.Sprintf("// Process 处理请求\n"))
		file.WriteString(fmt.Sprintf("// 该方法实现了处理接口\n"))
		file.WriteString(fmt.Sprintf("func (s *TestStruct%d) Process(ctx context.Context, input string) (string, error) {\n", i))

		file.WriteString("\t// 验证参数\n")
		file.WriteString("\tif input == \"\" {\n")
		if rand.Float64() < commentRatio {
			file.WriteString("\t\treturn \"\", errors.New(\"输入不能为空\") // 参数校验\n")
		} else {
			file.WriteString("\t\treturn \"\", errors.New(\"输入不能为空\")\n")
		}
		file.WriteString("\t}\n\n")

		file.WriteString("\t// 设置超时\n")
		file.WriteString("\ttimeout := time.Second * 5\n")
		file.WriteString("\tctx, cancel := context.WithTimeout(ctx, timeout)\n")
		file.WriteString("\tdefer cancel()\n\n")

		file.WriteString("\t// 记录开始时间\n")
		file.WriteString("\tstartTime := time.Now()\n\n")

		file.WriteString("\t// 处理逻辑\n")
		file.WriteString("\tif s.Enabled {\n")
		if rand.Float64() < commentRatio {
			file.WriteString("\t\ts.Value++ // 增加计数\n")
		} else {
			file.WriteString("\t\ts.Value++\n")
		}
		file.WriteString("\t}\n\n")

		file.WriteString("\t// 构建结果\n")
		file.WriteString("\tresult := fmt.Sprintf(\"处理结果: %s (耗时: %v)\", input, time.Since(startTime))\n")
		if rand.Float64() < commentRatio {
			file.WriteString("\treturn result, nil // 返回结果\n")
		} else {
			file.WriteString("\treturn result, nil\n")
		}
		file.WriteString("}\n\n")
	}

	// 补充更多内容以达到所需行数
	currentLineCount := 50 + 3*10 + 5*20 + 3*30 // 粗略估计已写入的行数
	for i := 0; currentLineCount < lineCount; i++ {
		file.WriteString(fmt.Sprintf("// HelperFunc%d 是辅助函数\n", i))
		file.WriteString(fmt.Sprintf("// 用于辅助测试\n"))
		file.WriteString(fmt.Sprintf("func HelperFunc%d(value string) string {\n", i))

		// 写入约10行函数内容
		for j := 0; j < 10; j++ {
			if rand.Float64() < 0.3 {
				if rand.Float64() < commentRatio {
					file.WriteString(fmt.Sprintf("\t// 这是第%d行代码的注释\n", j))
				}
			}

			// 使用格式正确的代码行
			codeLines := []string{
				"result := \"processed\" + value",
				"if len(value) == 0 { return \"empty\" }",
				"fmt.Println(value)",
				"time.Sleep(time.Millisecond)",
				"var buffer bytes.Buffer",
				"buffer.WriteString(value)",
			}

			// 有些代码行需要多行才能表达完整，单独处理
			complexLines := [][]string{
				{"for i := 0; i < len(value); i++ {",
					"\tfmt.Println(i)",
					"}"},
				{"ctx, cancel := context.WithTimeout(context.Background(), time.Second)",
					"defer cancel()"},
				{"if value == \"test\" {",
					"\treturn \"测试值\"",
					"}"},
			}

			// 70%概率使用简单代码行，30%概率使用复杂代码行
			if rand.Float64() < 0.7 {
				code := codeLines[rand.Intn(len(codeLines))]
				if rand.Float64() < commentRatio {
					file.WriteString(fmt.Sprintf("\t%s // 行尾注释\n", code))
				} else {
					file.WriteString(fmt.Sprintf("\t%s\n", code))
				}
			} else {
				// 使用复杂多行代码
				lines := complexLines[rand.Intn(len(complexLines))]
				for lineIdx, line := range lines {
					if lineIdx == 0 && rand.Float64() < commentRatio {
						file.WriteString(fmt.Sprintf("\t%s // 复杂代码块\n", line))
					} else {
						file.WriteString(fmt.Sprintf("\t%s\n", line))
					}
				}
				// 复杂代码行会增加行数，所以减少j的计数
				j += len(lines) - 1
			}
		}

		file.WriteString("\treturn value\n")
		file.WriteString("}\n\n")

		currentLineCount += 15 // 每个函数大约15行
	}
}

func TestMain(m *testing.M) {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 运行测试
	os.Exit(m.Run())
}

// 在测试前准备文件
func prepareTestFiles(b *testing.B) {
	// 生成不同规模的测试文件
	testFiles := []struct {
		name         string
		lineCount    int
		commentRatio float64
	}{
		{"realistic_small.go", 200, 0.3},
		{"realistic_medium.go", 500, 0.3},
		{"realistic_large.go", 2000, 0.3},
		{"realistic_high_comment.go", 1000, 0.5},
	}

	// 创建所有测试文件
	for _, f := range testFiles {
		filePath := filepath.Join("testdata", f.name)
		// 检查文件是否已存在，避免重复创建
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			createTestFile(b, f.name, f.lineCount, f.commentRatio)
		}
	}
}

// 小规模文件测试 (约200行)
func BenchmarkExtractComments_RealisticSmall(b *testing.B) {
	prepareTestFiles(b)
	filename := "testdata/realistic_small.go"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ExtractComments(filename)
		if err != nil {
			b.Fatalf("提取注释失败: %v", err)
		}
	}
}

// 中等规模文件测试 (约500行)
func BenchmarkExtractComments_RealisticMedium(b *testing.B) {
	prepareTestFiles(b)
	filename := "testdata/realistic_medium.go"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ExtractComments(filename)
		if err != nil {
			b.Fatalf("提取注释失败: %v", err)
		}
	}
}

// 大规模文件测试 (约2000行)
func BenchmarkExtractComments_RealisticLarge(b *testing.B) {
	prepareTestFiles(b)
	filename := "testdata/realistic_large.go"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ExtractComments(filename)
		if err != nil {
			b.Fatalf("提取注释失败: %v", err)
		}
	}
}

// 高注释密度文件测试 (约1000行，注释密度50%)
func BenchmarkExtractComments_RealisticHighComment(b *testing.B) {
	prepareTestFiles(b)
	filename := "testdata/realistic_high_comment.go"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ExtractComments(filename)
		if err != nil {
			b.Fatalf("提取注释失败: %v", err)
		}
	}
}

// 代码字符串测试 (不读取文件，直接处理字符串)
func BenchmarkExtractComments_RealisticCodeString(b *testing.B) {
	prepareTestFiles(b)
	// 读取示例文件内容作为代码字符串
	content, err := os.ReadFile("testdata/realistic_small.go")
	if err != nil {
		b.Fatalf("读取文件失败: %v", err)
	}

	codeString := string(content)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ExtractComments(codeString)
		if err != nil {
			b.Fatalf("提取注释失败: %v", err)
		}
	}
}

// 优化版测试

// 小规模文件测试 (约200行)
func BenchmarkExtractCommentsOptimized_RealisticSmall(b *testing.B) {
	prepareTestFiles(b)
	filename := "testdata/realistic_small.go"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ExtractCommentsOptimized(filename)
		if err != nil {
			b.Fatalf("提取注释失败: %v", err)
		}
	}
}

// 中等规模文件测试 (约500行)
func BenchmarkExtractCommentsOptimized_RealisticMedium(b *testing.B) {
	prepareTestFiles(b)
	filename := "testdata/realistic_medium.go"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ExtractCommentsOptimized(filename)
		if err != nil {
			b.Fatalf("提取注释失败: %v", err)
		}
	}
}

// 大规模文件测试 (约2000行)
func BenchmarkExtractCommentsOptimized_RealisticLarge(b *testing.B) {
	prepareTestFiles(b)
	filename := "testdata/realistic_large.go"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ExtractCommentsOptimized(filename)
		if err != nil {
			b.Fatalf("提取注释失败: %v", err)
		}
	}
}

// 高注释密度文件测试 (约1000行，注释密度50%)
func BenchmarkExtractCommentsOptimized_RealisticHighComment(b *testing.B) {
	prepareTestFiles(b)
	filename := "testdata/realistic_high_comment.go"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ExtractCommentsOptimized(filename)
		if err != nil {
			b.Fatalf("提取注释失败: %v", err)
		}
	}
}

// 代码字符串测试 (不读取文件，直接处理字符串)
func BenchmarkExtractCommentsOptimized_RealisticCodeString(b *testing.B) {
	prepareTestFiles(b)
	// 读取示例文件内容作为代码字符串
	content, err := os.ReadFile("testdata/realistic_small.go")
	if err != nil {
		b.Fatalf("读取文件失败: %v", err)
	}

	codeString := string(content)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ExtractCommentsOptimized(codeString)
		if err != nil {
			b.Fatalf("提取注释失败: %v", err)
		}
	}
}

// 缓存效率测试
func BenchmarkExtractCommentsOptimized_RealisticCacheEfficiency(b *testing.B) {
	prepareTestFiles(b)
	filename := "testdata/realistic_medium.go"

	// 预热缓存
	_, err := ExtractCommentsOptimized(filename)
	if err != nil {
		b.Fatalf("提取注释失败: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ExtractCommentsOptimized(filename)
		if err != nil {
			b.Fatalf("提取注释失败: %v", err)
		}
	}
}
