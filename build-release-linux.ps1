xgo -buildmode=c-shared --targets=linux/amd64 -out build/release/module/quickdemo github.com/grandprixgp/quickdemo
$env:GOOS="linux"; go build -o build/release/quickdemo.bin