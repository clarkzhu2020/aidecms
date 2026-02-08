# 🚀 快速开始 - 场景测试演示

## 方式1: 运行Go演示脚本 (推荐)

### 1. 启动后端服务
```bash
cd backend
go run main.go
```

### 2. 在新终端运行演示
```bash
cd backend
go run run-scene-demo.go
```

**预期输出**:
```
========================================
智能客服系统 - 场景演示
========================================
场景: 高置信度问答
问题: 如何使用产品A？

Session ID: scene-demo-1738900000

步骤1: 发送问题...
-----------------------------------
✅ 请求成功

步骤2: AI回答...
-----------------------------------
产品A的使用方法如下：首先打开应用，然后点击设置按钮...

步骤3: 置信度评估...
-----------------------------------
置信度: 0.92
进度条: ██████████ 92% (绿色)

步骤4: 验证置信度...
-----------------------------------
✅ 置信度检查通过: 0.92 >= 0.7 (状态: 高)

步骤5: 知识来源...
-----------------------------------
✅ 找到 1 个知识来源:
  [1] 产品A使用指南 (相关度: 92%)

步骤6: 建议操作...
-----------------------------------
✅ 建议操作:
  [1] 查看详细文档
  [2] 联系技术支持

步骤7: 转接判断...
-----------------------------------
✅ 未触发转接 (符合预期)

步骤8: 获取历史记录...
-----------------------------------
✅ 历史记录: 共 2 条消息

========================================
✅ 场景演示完成!
========================================

测试结果:
- 置信度: 0.92 (>=0.7 ✓)
- 知识来源: 1 个文档 ✓
- 转人工: 未触发 ✓
- 消息保存: 成功 ✓

🎉 完整7步流程验证通过!
```

---

## 方式2: 运行PowerShell脚本 (Windows)

### 1. 启动后端服务
```powershell
cd backend
go run main.go
```

### 2. 在新终端运行演示
```powershell
cd backend
.\run-scene-demo.ps1
```

**如果遇到执行策略限制**:
```powershell
# 临时允许脚本执行
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope Process
.\run-scene-demo.ps1
```

---

## 方式3: 运行批处理脚本

### 1. 启动后端服务
```cmd
cd backend
go run main.go
```

### 2. 在新终端运行测试
```cmd
cd backend
test-scene-high-confidence.bat
```

---

## 方式4: 运行单元测试

```bash
cd backend
go test -v ./pkg/css -run SceneHighConfidenceTest
```

---

## 方式5: 手动curl测试

### 1. 启动后端服务
```bash
cd backend
go run main.go
```

### 2. 发送测试请求
```bash
curl -X POST http://localhost:8888/api/css/question \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-001",
    "question": "如何使用产品A？"
  }'
```

**预期响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "answer": "产品A的使用方法如下：首先打开应用，然后点击设置按钮，选择产品A，按照屏幕提示完成配置即可。如果需要帮助，可以查看快速入门指南。",
    "confidence": 0.92,
    "sources": [
      {
        "document_id": "doc-product-a",
        "title": "产品A使用指南",
        "relevance": 0.92
      }
    ],
    "suggested_actions": [
      "查看详细文档",
      "联系技术支持"
    ]
  }
}
```

### 3. 验证响应
- ✅ 置信度 >= 0.7
- ✅ 包含知识来源
- ✅ 不包含 transfer_to 字段

### 4. 查看历史记录
```bash
curl http://localhost:8888/api/css/history/test-001
```

---

## 方式6: 前端界面测试

### 1. 启动后端
```bash
cd backend
go run main.go
```

### 2. 启动前端
```bash
cd web
npm install
npm run dev
```

### 3. 访问演示页面
打开浏览器访问: http://localhost:5173/customer-service-demo

### 4. 输入问题
在聊天框中输入: "如何使用产品A？"

### 5. 观察结果
- ✅ 显示AI回答
- ✅ 显示绿色置信度进度条 (92%)
- ✅ 显示知识来源 "产品A使用指南"
- ✅ 显示建议操作按钮
- ✅ 无转人工提示

---

## 📊 验证清单

运行完任何一种方式后,验证以下项目:

### ✅ 流程验证
- [ ] 会话已创建
- [ ] 用户消息已保存
- [ ] RAG检索返回相关文档
- [ ] AI回答有内容
- [ ] 置信度 >= 0.7
- [ ] 显示知识来源
- [ ] 显示建议操作
- [ ] 不触发转接
- [ ] AI消息已保存

### ✅ 预期结果
```
输入: "如何使用产品A？"
流程: 完整7步 ✅
预期: 置信度≥0.7, 显示知识来源, 不转接
状态: ✅ 已实现

验证结果:
- 置信度: 0.92 (≥0.7 ✓)
- 知识来源: 1个文档 ✓
- 转人工: 未触发 ✓
- 消息保存: 成功 ✓
```

---

## 🔧 故障排查

### 问题1: 连接失败
**错误**: `connect: connection refused`

**解决**:
```bash
# 检查后端是否启动
curl http://localhost:8888/api/css/status

# 查看端口占用
netstat -an | grep 8888  # Linux/Mac
netstat -an | findstr 8888  # Windows
```

### 问题2: CSS系统未启用
**错误**: `CSS system not enabled`

**解决**:
```bash
# 检查 .env 文件
cat backend/.env | grep CSS_ENABLED

# 应该是:
CSS_ENABLED=true
```

### 问题3: 置信度总是很低
**原因**: 知识库为空

**解决**:
```sql
-- 添加测试文档到知识库
INSERT INTO kb_documents (title, content, category, created_at)
VALUES (
    '产品A使用指南',
    '产品A是一款智能办公软件。使用方法：1. 打开应用 2. 点击设置 3. 选择产品A 4. 按照提示完成配置',
    'user-guide',
    NOW()
);
```

### 问题4: PowerShell执行策略限制
**错误**: `无法加载文件，因为在此系统上禁止运行脚本`

**解决**:
```powershell
# 查看当前策略
Get-ExecutionPolicy

# 临时允许
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope Process

# 或以管理员身份运行
Set-ExecutionPolicy RemoteSigned
```

---

## 📚 相关文档

| 文档 | 说明 |
|------|------|
| `QUICK_START.md` | 本文档 |
| `SCENARIOS.md` | 场景说明 |
| `doc/CSS_TEST_SCENARIOS.md` | 完整测试场景 |
| `doc/CSS_DEPLOYMENT.md` | 部署指南 |

---

## 🎉 下一步

测试通过后,可以尝试其他场景:

1. **低置信度转接**: 测试"我要投诉这个产品"
2. **多轮对话**: 连续提问测试上下文
3. **WebSocket测试**: 使用前端界面测试实时通讯
4. **性能测试**: 批量测试响应时间

祝测试顺利! 🚀
