FROM golang:alpine AS builder

WORKDIR /app

COPY * ./

RUN go build -o hello-world .

FROM scratch

WORKDIR /app

COPY --from=builder /app .

CMD [ "./hello-world" ]