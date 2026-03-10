@echo off
echo Installing Bambu Farm to Windows Startup...

set "VBS_SCRIPT=%~dp0start_hidden.vbs"
set "STARTUP_FOLDER=%APPDATA%\Microsoft\Windows\Start Menu\Programs\Startup"

if not exist "%VBS_SCRIPT%" (
    echo Error: Could not find start_hidden.vbs at %VBS_SCRIPT%
    pause
    exit /b 1
)

copy "%VBS_SCRIPT%" "%STARTUP_FOLDER%\BambuFarmStart.vbs"

echo.
echo Success! Bambu Farm will now automatically start hidden in the background when your computer boots up.
pause
