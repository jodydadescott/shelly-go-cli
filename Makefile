linux:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -trimpath -o shelly-cli shelly-cli.go
