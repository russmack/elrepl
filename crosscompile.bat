setlocal
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -o elrepl.bin elrepl.go command.go commandparser.go dispatcher.go httpclient.go handler.go handlealias.go handleclose.go handlecount.go handledir.go handledoc.go handleduplicatescount.go handleenv.go handleexit.go handleflush.go handleget.go handlehelp.go handleindex.go handleload.go handlelog.go handlemapping.go handleopen.go handleoptimize.go handlepost.go handleput.go handlerecovery.go handlerefresh.go handlereindex.go handlerun.go handlesegments.go handlesettings.go handlestats.go handlestatus.go handleversion.go
endlocal
pause
