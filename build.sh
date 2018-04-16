export CGO_ENABLED=0

osS='
linux
windows
'

archS='
amd64
'

for os in ${osS}; do
  for arch in ${archS}; do
    export GOOS=${os}
    export GOARCH=${arch}
    go build -o tskfs_${os}_${arch} main.go
    chmod a+x tskfs_${os}_${arch}
  done
done
