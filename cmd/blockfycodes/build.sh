#!/bin/bash

# 构建getcomments工具
echo "正在构建getcomments工具..."
go build -o blockfycodes main.go

# 检查构建结果
if [ $? -eq 0 ]; then
    echo "构建成功，可执行文件: blockfycodes"
    echo "用法: ./blockfycodes <文件路径或代码内容>"
else
    echo "构建失败"
    exit 1
fi 