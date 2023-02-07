@echo off

::::::::::::::  stop
echo Start Sync ...

D:/www/rsync/cwRsync_5.4/rsync.exe -avzP  --port=873 --delete --no-super --password-file=/cygdrive/D/www/rsync/cwRsync_5.4/pass.txt --exclude=logs/* --exclude=.git/ --exclude=bin/ --exclude=.idea/ /cygdrive/D/www/goweb/jayden/ root@192.168.92.208::go_jayden

echo Success...
:: 延时
choice /t 1 /d y /n >nul
::pause
exit