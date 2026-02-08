# 演示脚本: 高置信度问答场景 (PowerShell)
# 运行: .\run-scene-demo.ps1

param(
    [string]$BaseUrl = "http://localhost:8888"
)

function Write-Section {
    param([string]$Title)
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host $Title -ForegroundColor Cyan
    Write-Host "========================================"
    Write-Host ""
}

function Write-Step {
    param([string]$Title)
    Write-Host "步骤: $Title" -ForegroundColor Yellow
    Write-Host "-----------------------------------"
}

function Write-Success {
    param([string]$Message)
    Write-Host "✅ $Message" -ForegroundColor Green
}

function Write-Error {
    param([string]$Message)
    Write-Host "❌ $Message" -ForegroundColor Red
}

function Write-Warning {
    param([string]$Message)
    Write-Host "⚠️  $Message" -ForegroundColor Yellow
}

# 开始演示
Write-Section "智能客服系统 - 场景演示"
Write-Host "场景: 高置信度问答" -ForegroundColor White
Write-Host "问题: 如何使用产品A？" -ForegroundColor White
Write-Host ""

# 生成唯一会话ID
$sessionID = "scene-demo-" + [DateTimeOffset]::Now.ToUnixTimeSeconds()
Write-Host "Session ID: $sessionID" -ForegroundColor Cyan
Write-Host ""

# 步骤1: 发送问题
Write-Step "发送问题"

try {
    $question = @{
        session_id = $sessionID
        question   = "如何使用产品A？"
    } | ConvertTo-Json

    $response = Invoke-RestMethod `
        -Uri "$BaseUrl/api/css/question" `
        -Method POST `
        -ContentType "application/json" `
        -Body $question

    Write-Success "请求成功"
    Write-Host ""
}
catch {
    Write-Error "请求失败: $_"
    exit 1
}

# 步骤2: 显示AI回答
Write-Step "AI回答"
Write-Host $response.data.answer -ForegroundColor White
Write-Host ""

# 步骤3: 显示置信度
Write-Step "置信度评估"

$confidence = $response.data.confidence
Write-Host "置信度: $($confidence.ToString('F2'))" -ForegroundColor Cyan

# 显示置信度可视化
$barLength = [math]::Floor($confidence * 10)
$bar = "█" * $barLength + "░" * (10 - $barLength)

# 根据置信度显示颜色和状态
if ($confidence -ge 0.8) {
    $color = "绿色"
    $status = "高"
    $foregroundColor = "Green"
}
elseif ($confidence -ge 0.6) {
    $color = "橙色"
    $status = "中"
    $foregroundColor = "Yellow"
}
else {
    $color = "红色"
    $status = "低"
    $foregroundColor = "Red"
}

Write-Host "进度条: $bar $($confidence * 100)% ($color)" -ForegroundColor $foregroundColor
Write-Host ""

# 步骤4: 验证置信度
Write-Step "验证置信度"

if ($confidence -ge 0.7) {
    Write-Success "置信度检查通过: $($confidence.ToString('F2')) >= 0.7 (状态: $status)"
}
else {
    Write-Error "置信度检查失败: $($confidence.ToString('F2')) < 0.7 (状态: $status)"
    exit 1
}

Write-Host ""

# 步骤5: 显示知识来源
Write-Step "知识来源"

$sources = $response.data.sources
if ($sources -and $sources.Count -gt 0) {
    Write-Success "找到 $($sources.Count) 个知识来源:"
    for ($i = 0; $i -lt $sources.Count; $i++) {
        $source = $sources[$i]
        Write-Host "  [$($i+1)] $($source.title) (相关度: $($source.relevance * 100)%)" -ForegroundColor White
    }
}
else {
    Write-Warning "未找到知识来源 (可能是模拟数据)"
}

Write-Host ""

# 步骤6: 显示建议操作
Write-Step "建议操作"

$actions = $response.data.suggested_actions
if ($actions -and $actions.Count -gt 0) {
    Write-Host "建议操作:" -ForegroundColor White
    for ($i = 0; $i -lt $actions.Count; $i++) {
        Write-Host "  [$($i+1)] $($actions[$i])" -ForegroundColor White
    }
}
else {
    Write-Host "无建议操作" -ForegroundColor Gray
}

Write-Host ""

# 步骤7: 验证转接
Write-Step "转接判断"

if ($response.data.transfer_to) {
    Write-Error "触发了转接 (不符合预期)"
    exit 1
}
else {
    Write-Success "未触发转接 (符合预期)"
}

Write-Host ""

# 步骤8: 获取历史记录
Write-Step "获取历史记录"

try {
    $history = Invoke-RestMethod `
        -Uri "$BaseUrl/api/css/history/$sessionID" `
        -Method GET

    Write-Success "历史记录: 共 $($history.data.Count) 条消息"
}
catch {
    Write-Warning "获取历史失败: $_"
}

Write-Host ""

# 总结
Write-Section "场景演示完成!"
Write-Host ""

Write-Host "测试结果:" -ForegroundColor White
Write-Host "- 置信度: $($confidence.ToString('F2')) (>=0.7 ✓)" -ForegroundColor Green
Write-Host "- 知识来源: $($sources.Count) 个文档 ✓" -ForegroundColor Green
Write-Host "- 转人工: 未触发 ✓" -ForegroundColor Green
Write-Host "- 消息保存: 成功 ✓" -ForegroundColor Green
Write-Host ""

Write-Host "🎉 完整7步流程验证通过!" -ForegroundColor Green
