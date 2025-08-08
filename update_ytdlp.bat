@echo off
echo Updating yt-dlp to latest version...
curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe -o dependencies\yt-dlp.exe.tmp
if %ERRORLEVEL% EQU 0 (
    move dependencies\yt-dlp.exe.tmp dependencies\yt-dlp.exe
    echo yt-dlp updated successfully!
) else (
    echo Failed to download yt-dlp update
    if exist dependencies\yt-dlp.exe.tmp del dependencies\yt-dlp.exe.tmp
)
pause