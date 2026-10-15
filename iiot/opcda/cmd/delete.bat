@echo off
:: 切换到 ANSI 编码运行环境
chcp 65001 >nul
setlocal

set "TARGET_DIR=C:\Program Files (x86)\jy_opcda"
set "SERVICE_NAME=JY_AIOT_OPC_Proxy"
set "EXE_NAME=opc_32bit.exe"

echo [1/4] 正在停止服务...
sc stop "%SERVICE_NAME%" >nul 2>&1

echo [2/4] 强制结束残留进程...
:: 这一步非常重要，确保 exe 彻底退出释放文件锁
taskkill /f /im "%EXE_NAME%" >nul 2>&1
timeout /t 2 /nobreak >nul

echo [3/4] 正在删除服务登记...
sc delete "%SERVICE_NAME%" >nul 2>&1

echo [4/4] 正在清理目录...
:: 关键：先切出该目录，防止脚本自身占用文件夹
cd /d C:\
if exist "%TARGET_DIR%" (
    :: 使用 /s /q 强制删除
    rd /s /q "%TARGET_DIR%"
)

if exist "%TARGET_DIR%" (
    echo [错误] 目录仍未删除，请检查是否有文件被其他程序（如日志查看器）打开。
) else (
    echo [成功] 服务已卸载，目录已彻底清理。
)

pause