@echo off
set PGPASSWORD=postgres
psql -h localhost -U postgres -d forum -f migrations/000001_init.up.sql
pause 