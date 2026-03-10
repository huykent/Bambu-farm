@echo off
echo Removing Bambu Farm from Windows Startup...

set "STARTUP_FILE=%APPDATA%\Microsoft\Windows\Start Menu\Programs\Startup\BambuFarmStart.vbs"

if exist "%STARTUP_FILE%" (
    del "%STARTUP_FILE%"
    echo.
    echo Success! Bambu Farm has been removed from auto-startup.
) else (
    echo.
    echo Bambu Farm is not currently in the startup folder.
)

pause
