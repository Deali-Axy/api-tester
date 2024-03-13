# api-tester

接口测试工具

## build

build on linux, run on linux

```bash
CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build
```

build on windows, run on linux

```powershell
go env -w CGO_ENABLED=0
go env -w GOARCH=amd64
go env -w GOOS=linux
go build
```

## libraries

- https://github.com/codeskyblue/go-sh - 方便的在 go 里调用 shell

