local-bin: build
	cp ./cfn-compose ${GOPATH}/bin

build:
	go build -o cfn-compose .

test:
	go test ./... -coverprofile coverage/coverage.out