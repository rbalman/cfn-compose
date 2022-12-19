local-bin: build
	cp ./cfn-compose ${GOPATH}/bin

build:
	go build -o cfn-compose .

test:
	mkdir -p coverage
	go test ./... -coverprofile coverage/coverage.out
