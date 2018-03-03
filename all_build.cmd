SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -ldflags "-s -w" .\main\deploy.go

move  /Y deploy .\release\linux-64

SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -ldflags "-s -w" .\main\deploy.go

move  /Y deploy.exe .\release\windows-64

SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -ldflags "-s -w" .\main\deploy.go

move  /Y deploy .\release\darwin-64