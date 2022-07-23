# connection-killer

## Use

This was made for soloeing some raids in Destiny by causing your connection to "crash".

## Dependencies

To make this software run, you will need [Windivert](https://reqrypt.org/windivert.html). Put it next to the executable for it to work.

You will also need [Go](https://go.dev) to continue.

## Compiling

To compile, you will need to build it for Windows (since it's made for Windows..).

For 64bits :

```
go mod download
GOOS=windows GOARCH=amd64 go build -o bin/app.exe -ldflags -H=windowsgui cmd/connectionkiller/*.go
```

For 32bits :

```
go mod download
GOOS=windows GOARCH=386 go build -o bin/app.exe -ldflags -H=windowsgui cmd/connectionkiller/*.go
```