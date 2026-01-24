# AideCMS 框架测试

## 测试结构

```
test/
├── unit/            # 单元测试
│   ├── app_test.go      # 应用核心测试
│   └── router_test.go   # 路由测试
└── integration/     # 集成测试
    └── database_test.go # 数据库集成测试
```

## 运行测试

```bash
# 运行所有测试
go test ./...

# 运行单元测试
go test ./test/unit/...

# 运行集成测试
go test ./test/integration/...
```

## 测试覆盖率

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out