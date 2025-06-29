@echo off
echo ========================================
echo    罗小黑出行 - 前端项目启动脚本
echo ========================================
echo.

echo 检查Node.js版本...
node --version
if %errorlevel% neq 0 (
    echo 错误: 请先安装Node.js
    pause
    exit /b 1
)

echo.
echo 检查npm版本...
npm --version
if %errorlevel% neq 0 (
    echo 错误: npm未正确安装
    pause
    exit /b 1
)

echo.
echo 检查依赖包是否已安装...
if not exist "node_modules" (
    echo 正在安装依赖包...
    npm install
    if %errorlevel% neq 0 (
        echo 错误: 依赖包安装失败
        pause
        exit /b 1
    )
) else (
    echo 依赖包已存在，跳过安装
)

echo.
echo ========================================
echo 启动信息:
echo - 前端地址: http://localhost:3001
echo - 后端地址: http://localhost:8899 (需要手动启动)
echo - API路径:  /v1/api/*
echo ========================================
echo.
echo 正在启动前端开发服务器...
echo 按 Ctrl+C 可停止服务器
echo.

npm run dev

pause 