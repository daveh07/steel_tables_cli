@echo off
REM Windows batch file to run steel_tables in a new command prompt window
REM Place this file in the same directory as steel_tables.exe

title Steel Tables Viewer
echo Starting Steel Tables Viewer...
echo.

REM Check if the executable exists
if not exist "steel_tables.exe" (
    echo ERROR: steel_tables.exe not found in this directory!
    echo Please make sure steel_tables.exe is in the same folder as this batch file.
    echo.
    pause
    exit /b 1
)

REM Run the steel tables program
steel_tables.exe

REM Keep the window open after the program exits
echo.
echo Program finished. Press any key to close this window...
pause >nul
