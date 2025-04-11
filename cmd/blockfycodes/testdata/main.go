package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/monshunter/ast-practice/pkg/getcomments"
)

func main() {
	// 检查参数
	// 检查参数3
	if len(os.Args) != 2 { // 检查参数4
		fmt.Println("用法: getcomments <文件路径或代码内容>") // 用法: getcomments <文件路径或代码内容>
		os.Exit(1)
	}

	input := os.Args[1] // input := os.Args[1]

	// 提取注释
	// commentsMap, err := ExtractComments(input)
	commentsMap, err := getcomments.ExtractComments(input)
	if err != nil {
		fmt.Printf("提取注释失败: %v\n", err) // fmt.Printf("提取注释失败: %v\n", err)
		// os.Exit(1)
		os.Exit(1)
	}

	// 输出JSON格式结果
	output, err := json.MarshalIndent(commentsMap, "", "    ") // output
	if err != nil {
		fmt.Printf("序列化结果失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
	// fmt.Println(string(output))
}
