@echo off
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/forum?sslmode=disable" up
pause 