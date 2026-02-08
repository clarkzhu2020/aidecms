@echo off
REM 场景测试: 高置信度问答
REM 测试问题: "如何使用产品A?"
REM 预期: 置信度>=0.7, 显示知识来源, 不转接

setlocal enabledelayedexpansion

if "%CSS_BASE_URL%"=="" set CSS_BASE_URL=http://localhost:8888
set SESSION_ID=test-scene-%RANDOM%

echo =========================================
echo 场景测试: 高置信度问答
echo =========================================
echo 问题: 如何使用产品A?
echo 预期: 置信度^>=0.7, 显示知识来源, 不转接
echo.

REM 步骤1: 发送问题
echo 步骤1: 发送问题...
echo -----------------------------------

curl -s -X POST "%CSS_BASE_URL%/api/css/question" ^
  -H "Content-Type: application/json" ^
  -d "{\"session_id\": \"%SESSION_ID%\", \"question\": \"如何使用产品A？\"}" > response.json

type response.json
echo.

REM 步骤2: 验证响应
echo 步骤2: 验证响应...
echo -----------------------------------

REM 这里使用 PowerShell 来解析 JSON
for /f "delims=" %%i in ('powershell -Command "(Get-Content response.json | ConvertFrom-Json).confidence"') do set CONFIDENCE=%%i
echo 置信度: %CONFIDENCE%

REM 验证置信度>=0.7
powershell -Command "if (%CONFIDENCE% -ge 0.7) { exit 0 } else { exit 1 }"
if %errorlevel% equ 0 (
    echo [OK] 置信度检查通过 (^>=0.7)
) else (
    echo [FAIL] 置信度检查失败 (^<0.7)
    exit /b 1
)

REM 检查是否有转接
powershell -Command "if ((Get-Content response.json | ConvertFrom-Json).transfer_to -eq $null) { exit 0 } else { exit 1 }"
if %errorlevel% equ 0 (
    echo [OK] 未触发转接
) else (
    echo [FAIL] 错误: 触发了转接
    exit /b 1
)

echo.

REM 步骤3: 获取对话历史
echo 步骤3: 获取对话历史...
echo -----------------------------------

curl -s "%CSS_BASE_URL%/api/css/history/%SESSION_ID%" > history.json
type history.json
echo.

REM 步骤4: 验证消息保存
echo 步骤4: 验证消息保存...
echo -----------------------------------

for /f "delims=" %%i in ('powershell -Command "(Get-Content history.json | ConvertFrom-Json | Where-Object {$_.role -eq 'user'} | Measure-Object).Count"') do set USER_COUNT=%%i
for /f "delims=" %%i in ('powershell -Command "(Get-Content history.json | ConvertFrom-Json | Where-Object {$_.role -eq 'assistant'} | Measure-Object).Count"') do set AI_COUNT=%%i

echo 用户消息数: %USER_COUNT%
echo AI回复数: %AI_COUNT%

if %USER_COUNT% gtr 0 if %AI_COUNT% gtr 0 (
    echo [OK] 消息已保存
) else (
    echo [FAIL] 消息保存失败
    exit /b 1
)

echo.
echo =========================================
echo [OK] 场景测试通过!
echo =========================================
echo.
echo 测试结果:
echo - 置信度: %CONFIDENCE% (^>=0.7 OK)
echo - 转人工: 未触发 OK
echo - 消息保存: 成功 OK
echo.

REM 清理临时文件
del response.json history.json
endlocal
