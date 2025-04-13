#!/bin/bash

# 构建getcomments工具
echo "正在构建topomain工具..."
go build -o topomain main.go

# 检查构建结果
if [ $? -eq 0 ]; then
    echo "构建成功，可执行文件: topomain"
    echo "用法: ./topomain <项目路径>"
else
    echo "构建失败"
    exit 1
fi
