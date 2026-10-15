@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

:: 设置变量
set "TARGET_DIR=C:\Program Files (x86)\jy_opcda"
set "EXE_NAME=opc_32bit.exe"
set "SERVICE_NAME=JY_AIOT_OPC_Proxy"

echo ==========================================
echo           服务部署工具
echo ==========================================
echo.

:: 1. 交互输入端口号
set "PORT="
echo [提示] 请输入服务运行端口 (直接回车将使用默认值: 11600)
set /p PORT="> "

:: 如果用户直接回车，PORT 变量会为空，此时赋值为 11600
if "!PORT!"=="" (
    set "PORT=11600"
)

echo.
echo [信息] 确认使用端口: !PORT!
echo.

echo [1/4] 正在创建安装目录...
if not exist "%TARGET_DIR%" (
    mkdir "%TARGET_DIR%"
)

echo [2/4] 正在拷贝程序文件...
:: 确保 exe 文件与此 bat 在同一目录下
copy /Y "%~dp0%EXE_NAME%" "%TARGET_DIR%\"

echo [3/4] 正在重置旧服务...
sc stop "%SERVICE_NAME%" >nul 2>&1
sc delete "%SERVICE_NAME%" >nul 2>&1
timeout /t 1 /nobreak >nul

echo [4/4] 正在创建 Windows 服务 (端口: !PORT!)...
:: binPath 的引号嵌套是确保带空格路径和参数生效的关键
sc create "%SERVICE_NAME%" binPath= "\"%TARGET_DIR%\%EXE_NAME%\" -p !PORT!" start= auto

echo [5/5] 正在配置服务自愈恢复功能...
:: 失败 1 分钟后重启，1 天后重置失败计数
sc failure "%SERVICE_NAME%" reset= 86400 actions= restart/60000/restart/60000/restart/60000

echo.
echo ==========================================
echo 部署完成！
echo 服务名称: %SERVICE_NAME%
echo 启动参数: -p !PORT!
echo 安装路径: %TARGET_DIR%
echo ==========================================
echo 正在尝试启动服务...
sc start "%SERVICE_NAME%"

echo.
pause