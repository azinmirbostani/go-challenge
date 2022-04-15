deploy:
	export GOARCH=amd64
	export GOOS=linux
	go build -ldflags="-s -w" -o build/main cmd/main.go
	build-lambda-zip -o build/main.zip build/main
	sls deploy

