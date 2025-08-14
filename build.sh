#!/bin/bash
# by zgcwkj
echo "build start"
# 根目录
rootPath=$(pwd)
outPath="$rootPath/build/filebrowser_"
# 编译目标平台
windows_archs="386 amd64 arm64"
linux_archs="386 amd64 arm64 mips mipsle mips64 mips64le"
darwin_archs="amd64 arm64"
freebsd_archs="386 amd64 arm64"
# 开始编译
for os in windows linux darwin freebsd; do
    archs="${os}_archs"
    IFS=" " # 设置分割字符
    for arch in ${!archs}; do
        echo "building for $os/$arch"
        # 设置编译环境变量
        export CGO_ENABLED=0
        export GOOS=$os
        export GOARCH=$arch
        # 设置可执行文件后缀
        exe_suffix=""
        if [ "$os" = "windows" ]; then
            exe_suffix=".exe"
        fi
        # 编译程序
        cd "$rootPath"
        outputFile="${outPath}${os}_${arch}${exe_suffix}"
        go build -ldflags="-w -s" -trimpath -o "$outputFile"
        # 压缩文件
        zip -j "${outputFile}.zip" "$outputFile"
        # 压缩后删除原文件
        rm "$outputFile"
    done
done
# 回到根目录
cd "$rootPath"
# 结束
read -p "Press any key to continue..." -n1 -s
