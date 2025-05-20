@echo off
go test -v -coverprofile=coverage.out ./internal/repository/...
go tool cover -func=coverage.out
go tool cover -html=coverage.out 