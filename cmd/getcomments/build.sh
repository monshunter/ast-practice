#!/bin/bash

# 构建getcomments工具
echo "正在构建getcomments工具..."
go build -o getcomments main.go

# 检查构建结果
if [ $? -eq 0 ]; then
    echo "构建成功，可执行文件: getcomments"
    echo "用法: ./getcomments <文件路径或代码内容>"
    
    # 移动到系统路径（可选）
    # read -p "是否要安装到系统路径？(y/n): " answer
    # if [ "$answer" = "y" ]; then
    #     sudo mv getcomments /usr/local/bin/
    #     echo "已安装到 /usr/local/bin/getcomments"
    # fi
else
    echo "构建失败"
    exit 1
fi 