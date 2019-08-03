xgo -buildmode=c-shared --targets=linux/amd64,windows/amd64 -out build/release/module/quickdemo github.com/grandprixgp/quickdemo
$env:GOOS="windows";$env:GOARCH="amd64"; go build -o build/release/quickdemo.exe
$env:GOOS="linux";$env:GOARCH="amd64"; go build -o build/release/quickdemo.bin