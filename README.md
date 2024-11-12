<p align="center">
  <img src="https://raw.githubusercontent.com/filebrowser/logo/master/banner.png" width="550"/>
</p>

[![Build](https://github.com/filebrowser/filebrowser/actions/workflows/main.yaml/badge.svg)](https://github.com/filebrowser/filebrowser/actions/workflows/main.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/filebrowser/filebrowser?style=flat-square)](https://goreportcard.com/report/github.com/filebrowser/filebrowser)
[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/filebrowser/filebrowser)
[![Version](https://img.shields.io/github/release/filebrowser/filebrowser.svg?style=flat-square)](https://github.com/filebrowser/filebrowser/releases/latest)
[![Chat IRC](https://img.shields.io/badge/freenode-%23filebrowser-blue.svg?style=flat-square)](http://webchat.freenode.net/?channels=%23filebrowser)

filebrowser provides a file managing interface within a specified directory and it can be used to upload, delete, preview, rename and edit your files. It allows the creation of multiple users and each user can have its own directory. It can be used as a standalone app.

# run & build

## Vue

Env
```
export NODE_OPTIONS=--openssl-legacy-provider
```

Build
```
cd frontend

npm install

npm run build
```

## Golang

Env
```
export CGO_ENABLED=0
export GOOS=linux 
export GOARCH=mipsle

export CGO_ENABLED=0
export GOOS=windows
export GOARCH=amd64

export CGO_ENABLED=0
export GOOS=darwin
export GOARCH=amd64
```

Build
```
go run main.go

go build
```

Init
```
./filebrowser config set --address 0.0.0.0
./filebrowser config set --port 8080
./filebrowser config set --root /
```