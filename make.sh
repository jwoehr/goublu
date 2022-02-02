go build -ldflags "-X main.compileDate=`date -u +%Y-%m-%d.%H:%M:%S\(UTC\)` -X main.goubluVersion=v1.4.0" .
