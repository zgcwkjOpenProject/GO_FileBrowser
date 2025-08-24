@echo off
rem by zgcwkj
echo build start
rem 根目录
set rootPath=%cd%
set outPath=%rootPath%/build/filebrowser_
rem 编译目标平台
set windows_archs=386 amd64 arm64
set linux_archs=386 amd64 arm64 mips mipsle mips64 mips64le
set darwin_archs=amd64 arm64
set freebsd_archs=386 amd64 arm64
rem 启用延迟变量扩展功能
setlocal EnableDelayedExpansion
rem 开始编译
for %%o in (windows linux darwin freebsd) do (
    for %%b in (!%%o_archs!) do (
            echo building for %%o/%%b
            rem 设置编译环境变量
            set CGO_ENABLED=0
            set GOOS=%%o
            set GOARCH=%%b
            rem 设置可执行文件后缀
            set exe_suffix=
            if "%%o"=="windows" set exe_suffix=.exe
            rem 编译程序
            cd %rootPath%
            go build -ldflags="-w -s" -trimpath -o %outPath%%%o_%%b!exe_suffix!
        )
    )
)
rem 回到根目录
cd %rootPath%
rem 结束
pause
