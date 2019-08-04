# executables
xgo --targets=linux/amd64,windows/amd64 -out build/release/quickdemo .
# modules
xgo -buildmode=c-shared --targets=linux/amd64,windows/amd64 -out build/release/quickdemo .