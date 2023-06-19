FROM golang:1.19 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o cfnc

FROM scratch as final

COPY --from=builder /app/cfnc /cfnc

ENTRYPOINT [ "/cfnc" ]
