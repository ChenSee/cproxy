CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags "-s -w" -o cproxy main.go
upx --brute cproxy