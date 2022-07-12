bin: build
	cp ./cfn-deploy ${GOPATH}/bin

build:
	go build -o cfn-deploy .

test:
	go test ./...