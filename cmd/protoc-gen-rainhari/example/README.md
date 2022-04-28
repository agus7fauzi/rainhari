# protoc-gen-rainhari

Build and install binary Linux
```console
go build -o ~/go/bin/protoc-gen-rainhari ../main.go
```

Build and install binary Windows Powershell
```console
go build -o $env:USERPROFILE\go\bin\protoc-gen-rainhari.exe ../main.go
```

Build and install binary Windows CMD
```console
go build -o %USERPROFILE%\go\bin\protoc-gen-rainhari.exe ../main.go
```

Install Protoc Gen Go
```console
go install google.golang.org/protobuf/cmd/protoc-gen-go
```

Generate proto
```console
protoc --go_out=. --go_opt=paths=source_relative --rainhari_out=. --rainhari_opt=paths=source_relative -I . greeting.proto
```