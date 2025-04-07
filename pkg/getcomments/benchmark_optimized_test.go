package getcomments

import (
	"os"
	"testing"
)

// 小规模文件测试 (约100行)
func BenchmarkExtractCommentsOptimized_Small(b *testing.B) {
	filename := "testdata/examples.go"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ExtractCommentsOptimized(filename)
		if err != nil {
			b.Fatalf("提取注释失败: %v", err)
		}
	}
}

// 中等规模文件测试 (约500行)
func BenchmarkExtractCommentsOptimized_Medium(b *testing.B) {
	filename := createLargeTestFile(b, 500, 0.3) // 30%的行包含注释
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ExtractCommentsOptimized(filename)
		if err != nil {
			b.Fatalf("提取注释失败: %v", err)
		}
	}
}

// 大规模文件测试 (约2000行)
func BenchmarkExtractCommentsOptimized_Large(b *testing.B) {
	filename := createLargeTestFile(b, 2000, 0.3) // 30%的行包含注释
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ExtractCommentsOptimized(filename)
		if err != nil {
			b.Fatalf("提取注释失败: %v", err)
		}
	}
}

// 大规模文件测试 + 高密度注释 (约2000行，50%注释)
func BenchmarkExtractCommentsOptimized_Large_HighDensity(b *testing.B) {
	filename := createLargeTestFile(b, 2000, 0.5) // 50%的行包含注释
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ExtractCommentsOptimized(filename)
		if err != nil {
			b.Fatalf("提取注释失败: %v", err)
		}
	}
}

// 多文件组合测试
func BenchmarkExtractCommentsOptimized_MultiFiles(b *testing.B) {
	// 准备几个不同大小的文件
	smallFile := "testdata/examples.go"
	mediumFile := createLargeTestFile(b, 500, 0.3)
	largeFile := createLargeTestFile(b, 1000, 0.3)

	files := []string{smallFile, mediumFile, largeFile}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// 循环处理所有文件
		for _, file := range files {
			_, err := ExtractCommentsOptimized(file)
			if err != nil {
				b.Fatalf("提取注释失败 (文件: %s): %v", file, err)
			}
		}
	}
}

// 代码字符串测试 (不读取文件，直接处理字符串)
func BenchmarkExtractCommentsOptimized_CodeString(b *testing.B) {
	// 读取示例文件内容作为代码字符串
	content, err := readFileToString("testdata/examples.go")
	if err != nil {
		b.Fatalf("读取文件失败: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ExtractCommentsOptimized(content)
		if err != nil {
			b.Fatalf("提取注释失败: %v", err)
		}
	}
}

// 缓存优化测试 - 重复处理同一文件
func BenchmarkExtractCommentsOptimized_CacheEfficiency(b *testing.B) {
	filename := "testdata/examples.go"

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

// 辅助函数，读取文件内容为字符串
func readFileToString(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
