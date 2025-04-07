package getcomments

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// 创建不同规模的测试文件
func createLargeTestFile(b *testing.B, lines int, commentsPercentage float64) string {
	dir := "testdata"
	// 确保目录存在
	if err := os.MkdirAll(dir, 0755); err != nil {
		b.Fatalf("创建测试目录失败: %v", err)
	}

	filename := filepath.Join(dir, fmt.Sprintf("large_test_%d.go", lines))

	// 检查文件是否已存在
	if _, err := os.Stat(filename); err == nil {
		return filename // 文件已存在，直接返回文件名
	}

	// 创建文件
	file, err := os.Create(filename)
	if err != nil {
		b.Fatalf("创建测试文件失败: %v", err)
	}
	defer file.Close()

	// 写入包声明
	file.WriteString("package testdata\n\n")
	file.WriteString("// 导入标准库\n")
	file.WriteString("import (\n\t\"fmt\"\n\t\"math\"\n)\n\n")

	// 生成具有注释的函数
	commentLines := int(float64(lines) * commentsPercentage)
	normalLines := lines - commentLines - 10 // 减去包声明、导入等行数

	// 写入主函数
	file.WriteString("// TestFunction 是一个测试函数\n")
	file.WriteString("// 它包含了很多注释\n")
	file.WriteString("// 用于测试注释提取工具的性能\n")
	file.WriteString("func TestFunction() {\n")

	// 写入函数体
	for i := 0; i < normalLines; i++ {
		if i%5 == 0 && commentLines > 0 {
			// 每5行添加一个注释
			file.WriteString(fmt.Sprintf("\t// 这是第%d行的注释\n", i))
			commentLines--
		}
		file.WriteString(fmt.Sprintf("\tfmt.Println(\"line %d\")\n", i))
	}

	// 写入额外的注释行（如果还有剩余的注释配额）
	for commentLines > 0 {
		file.WriteString(fmt.Sprintf("\t// 额外的注释行 %d\n", commentLines))
		commentLines--
	}

	file.WriteString("}\n")
	return filename
}

// 小规模文件测试 (约100行)
func BenchmarkExtractComments_Small(b *testing.B) {
	filename := "testdata/examples.go"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ExtractComments(filename)
		if err != nil {
			b.Fatalf("提取注释失败: %v", err)
		}
	}
}

// 中等规模文件测试 (约500行)
func BenchmarkExtractComments_Medium(b *testing.B) {
	filename := createLargeTestFile(b, 500, 0.3) // 30%的行包含注释
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ExtractComments(filename)
		if err != nil {
			b.Fatalf("提取注释失败: %v", err)
		}
	}
}

// 大规模文件测试 (约2000行)
func BenchmarkExtractComments_Large(b *testing.B) {
	filename := createLargeTestFile(b, 2000, 0.3) // 30%的行包含注释
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ExtractComments(filename)
		if err != nil {
			b.Fatalf("提取注释失败: %v", err)
		}
	}
}

// 大规模文件测试 + 高密度注释 (约2000行，50%注释)
func BenchmarkExtractComments_Large_HighDensity(b *testing.B) {
	filename := createLargeTestFile(b, 2000, 0.5) // 50%的行包含注释
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ExtractComments(filename)
		if err != nil {
			b.Fatalf("提取注释失败: %v", err)
		}
	}
}

// 多文件组合测试
func BenchmarkExtractComments_MultiFiles(b *testing.B) {
	// 准备几个不同大小的文件
	smallFile := "testdata/examples.go"
	mediumFile := createLargeTestFile(b, 500, 0.3)
	largeFile := createLargeTestFile(b, 1000, 0.3)

	files := []string{smallFile, mediumFile, largeFile}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// 循环处理所有文件
		for _, file := range files {
			_, err := ExtractComments(file)
			if err != nil {
				b.Fatalf("提取注释失败 (文件: %s): %v", file, err)
			}
		}
	}
}

// 代码字符串测试 (不读取文件，直接处理字符串)
func BenchmarkExtractComments_CodeString(b *testing.B) {
	// 读取示例文件内容作为代码字符串
	content, err := os.ReadFile("testdata/examples.go")
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
