<p align="center">
  <img src="https://raw.githubusercontent.com/filebrowser/logo/master/banner.png" width="550"/>
</p>

File Browser provides a file managing interface within a specified directory and it can be used to upload, delete, preview and edit your files. It is a **create-your-own-cloud**-kind of software where you can just install it on your server, direct it to a path and access your files through a nice web interface.

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

go build -ldflags="-w -s" -trimpath
```

Init
```
./filebrowser config set --address 0.0.0.0
./filebrowser config set --port 8080
./filebrowser config set --root /
```