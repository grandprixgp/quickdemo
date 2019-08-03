$env:GOOS="windows";$env:GOARCH="amd64"; go build -o build/development/quickdemo.exe
$env:GOOS="linux";$env:GOARCH="amd64"; go build -o build/development/quickdemo.bin