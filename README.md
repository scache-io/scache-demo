# SCache Demo

这是一个演示 SCache 功能的示例项目，展示了如何在实际项目中使用 SCache 缓存库。

## 项目结构

```
scache-demo/
├── go.mod              # Go 模块定义
├── go.sum              # 依赖校验
├── main.go             # 演示主程序
├── models/
│   └── user.go         # 数据模型
└── README.md           # 项目说明
```

## 功能演示

本项目包含以下演示场景：

### 1. 基础缓存操作
- 存储和读取结构体数据
- 检查键是否存在
- 获取缓存大小

### 2. TTL 过期机制
- 设置短期缓存（验证码等）
- 演示缓存过期
- 动态修改 TTL

### 3. 多数据类型
- String 类型缓存
- List 类型缓存
- Hash 类型缓存

### 4. 并发安全
- 多协程并发写入
- 多协程并发读取
- 线程安全保证

### 5. 性能测试
- 写入性能测试
- 读取性能测试
- 缓存命中率统计

## 快速开始

### 安装依赖

```bash
go mod tidy
```

### 运行演示

```bash
go run main.go
```

## 使用 scache 生成缓存代码

### 方式一：使用泛型版本（推荐）

```bash
# 生成泛型版本缓存代码
scache gen --generic

# 指定目录生成
scache gen --generic -dir ./models

# 只生成指定结构体
scache gen --generic -structs User,Product,Session
```

### 方式二：使用传统版本

```bash
# 生成传统版本缓存代码
scache gen

# 指定目录生成
scache gen -dir ./models
```

## 代码示例

### 基础使用

```go
package main

import (
    "time"
    "github.com/scache-io/scache"
    "github.com/scache-io/scache-demo/models"
)

func main() {
    // 存储数据
    user := models.User{ID: 1, Name: "张三"}
    scache.Store("user:1", &user, time.Hour)

    // 读取数据
    var loadedUser models.User
    err := scache.Load("user:1", &loadedUser)
    if err == nil {
        fmt.Printf("用户: %+v\n", loadedUser)
    }
}
```

### 使用生成的代码（泛型版本）

```go
package main

import (
    "time"
    "github.com/scache-io/scache-demo/models"
)

func main() {
    // 获取用户缓存实例
    userCache := models.GetUserScache()

    // 存储数据
    user := models.User{ID: 1, Name: "张三"}
    userCache.Store("user:1", user, time.Hour)

    // 读取数据
    loadedUser, err := userCache.Load("user:1")
    if err == nil {
        fmt.Printf("用户: %+v\n", loadedUser)
    }
}
```

## 运行结果

```
=== SCache 演示项目 ===

1. 基础缓存操作
-----------------
✓ 存储用户数据: user:1
✓ 读取用户数据: {ID:1 Name:张三 Email:zhangsan@example.com Age:25 CreatedAt:2026-03-04 09:48:00 +0800 CST}
✓ 用户缓存存在
✓ 缓存大小: 1

2. TTL 过期演示
---------------
✓ 设置验证码缓存（2秒过期）
✓ 验证码（立即读取）: 123456
⏳ 等待3秒...
✓ 验证码已过期
✓ 设置用户信息缓存（1小时过期）
✓ 修改过期时间为5分钟
✓ 剩余过期时间: 4m59.999s

3. 多数据类型演示
----------------
✓ 字符串: app:name = SCache Demo
  读取: SCache Demo
✓ 列表: tags:go = [Go 缓存 高性能 泛型]
  读取: [Go 缓存 高性能 泛型]
✓ 哈希: profile:1001 = {name: 李四, email: lisi@example.com, city: 北京}
  读取: map[city:北京 email:lisi@example.com gender:男 name:李四]

4. 并发安全演示
---------------
✓ 启动 10 个协程，每个写入 100 次
✓ 并发操作完成
  缓存大小: 1000
  写入次数: 1000
  命中次数: 1000
  未命中次数: 0

5. 性能对比演示
---------------
测试写入性能...
✓ 写入 1000 条数据，耗时: 15.234ms (65657.34 条/秒)
测试读取性能...
✓ 读取 1000 条数据，耗时: 8.567ms (116732.42 条/秒)
✓ 最终统计:
  命中率: 100.00%
  缓存大小: 1000
  内存使用: 125000 bytes

=== 演示完成 ===
```

## 特性说明

- ✅ **高性能**: 基于内存存储，读写性能优异
- ✅ **线程安全**: 内置锁机制，支持高并发
- ✅ **TTL 过期**: 灵活的缓存过期机制
- ✅ **LRU 淘汰**: 智能的缓存淘汰策略
- ✅ **多数据类型**: 支持 String、List、Hash、Struct
- ✅ **泛型支持**: 类型安全的泛型 API
- ✅ **自动生成**: 代码生成工具，减少重复代码

## 适用场景

- 用户会话管理
- 数据库查询结果缓存
- API 响应缓存
- 配置信息缓存
- 热点数据缓存
- 临时数据存储（验证码等）

## 许可证

MIT License

## 相关链接

- [SCache 主项目](https://github.com/scache-io/scache)
- [SCache 文档](https://github.com/scache-io/scache/blob/main/README.md)
