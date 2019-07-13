$env:GOOS="windows";$env:GOARCH="amd64"; go build -o quickdemo.exe
$env:GOOS="linux";$env:GOARCH="amd64"; go build -o quickdemo.bin