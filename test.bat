@echo off
echo Running tests...

go test -v -cover -coverprofile=c.out
go tool cover -html=c.out -o coverage.html