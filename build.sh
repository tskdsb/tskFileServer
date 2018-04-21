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
    go build -o tskfs_${GOOS}_${GOARCH} main.go
  done
done

export GOARM=7
export GOOS=linux
export GOARCH=arm
go build -o tskfs_${GOOS}_${GOARCH}v${GOARM} main.go
