go generate
go build
go build -ldflags "-H=windowsgui -s -w"  .