# protoc-gen-rainhari

Build and install binary
```console
go build -o ~/go/bin/protoc-gen-rainhari ../main.go ../rainhari.go
```

Generate proto
```console
protoc --go_out=. --go_opt=paths=source_relative --rainhari_out=. --rainhari_opt=paths=source_relative -I . greeting.proto
```