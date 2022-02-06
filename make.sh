go build -ldflags "-X main.CompileDate=`date -u +%Y-%m-%d.%H:%M:%S\(UTC\)` -X main.GoubluVersion=v1.4.0" .
