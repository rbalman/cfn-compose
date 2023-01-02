local-bin: build
	cp ./cfnc ${GOPATH}/bin

build:
	go build -o cfnc .

test:
	mkdir -p coverage
	go test ./... -coverprofile coverage/coverage.out
