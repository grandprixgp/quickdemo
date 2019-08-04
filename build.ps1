# windows executable
go build -o "build/quickdemo-windows-amd64.exe"

#windows module
go build -buildmode=c-shared -o "build/quickdemo-windows-amd64.dll"

# linux executable
xgo --targets=linux/amd64 -out build/quickdemo .

# linux module
xgo -buildmode=c-shared --targets=linux/amd64 -out build/quickdemo .