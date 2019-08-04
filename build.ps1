# executables
xgo --targets=linux/amd64,windows/amd64 -out build/quickdemo .
# modules
xgo -buildmode=c-shared --targets=linux/amd64,windows/amd64 -out build/quickdemo .