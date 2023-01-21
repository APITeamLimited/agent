@REM Remove build directory if it exists
IF EXIST build-windows rmdir /s /q build-windows 

@REM Create build directory
mkdir build-windows

@REM Build agent
go build -o build-windows/apiteam-agent.exe -tags windows -ldflags -H=windowsgui
go-msi make --path targets/windows/wix.json --msi build-windows/apiteam-agent.msi --version 0.1.20 --license LICENSE.md

@REM Cleanup
del build-windows\apiteam-agent.exe
