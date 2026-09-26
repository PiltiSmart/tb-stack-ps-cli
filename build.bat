@echo off
if not exist bin mkdir bin
go mod tidy
go build -ldflags="-s -w" -o bin\pilti.exe main.go
echo [SUCCESS] Built bin\pilti.exe
