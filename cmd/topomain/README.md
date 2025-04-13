# golang main 内部拓扑

## 功能介绍
利用golang ast 分别分析golang project内部所有main package的 package拓扑纵深关系图(简化：仅限project的内部package，不包含外部或者标准std的package, 不需要构建拓扑关系图，只要求连通量)

## 输入与输出
### 输入:
- 一个golang的project目录

### 输出:
JSON格式的键值对:
```json
[{
  "mainDir": "",
  "mainFile": "",
  "functions": [
    "package1",
    "package2",
    "package3"
  ]
}]
```

## 使用方法

### 编译
```bash
cd cmd/topomain
go build
```

### 运行
```bash
# 分析当前目录
./topomain -pkg

# 分析指定目录
./topomain -pkg -dir=/path/to/golang/project

# 分析函数调用关系(旧功能)
./topomain -dir=/path/to/golang/project
```

### 输出说明
- `mainDir`: main包所在的目录
- `mainFile`: 包含main函数的文件路径
- `functions`: 由main函数直接或间接导入的项目内部包列表

## 实现说明
1. 使用go/ast包解析Go源码
2. 查找项目中所有的main包
3. 分析每个main包导入的所有内部包
4. 仅保留项目内部包，过滤掉外部包和标准库
5. 输出JSON格式的结果

## 特点
1. 能够准确分析包级别的依赖关系，不受接口和泛型影响
2. 自动识别项目内部包，排除标准库和第三方包
3. 支持go modules项目，自动识别模块前缀
4. 多个main包可同时分析

## 限制
1. 仅分析静态导入关系，无法分析运行时动态导入
2. 不包含通过反射使用的包依赖
3. 必须是正确编译的Go代码才能分析
