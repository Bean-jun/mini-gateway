set EXECUTE=server.exe

TASKLIST | FINDSTR %EXECUTE% && TASKKILL /F /IM %EXECUTE% 
IF EXIST %EXECUTE% del %EXECUTE%
go build -ldflags "-w -s" -trimpath
cmd /C "start cmd /K %EXECUTE% -port 5000"
cmd /C "start cmd /K %EXECUTE% -port 5001"