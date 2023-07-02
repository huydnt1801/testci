FROM golang:1.19 AS builder
WORKDIR /app

RUN curl -LO https://release.ariga.io/atlas/atlas-linux-amd64-latest && \
        mv atlas-linux-amd64-latest atlas && chmod +x atlas

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/server cmd/app/main.go

FROM ubuntu:22.04
RUN apt update && apt install -y ca-certificates \
        && rm -rf /var/lib/apt/lists/*
EXPOSE 5000
WORKDIR /app
COPY --from=builder /app/server /app/server
CMD ["/app/server"]
