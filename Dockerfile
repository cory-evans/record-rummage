FROM --platform=$BUILDPLATFORM golang:alpine AS builder
ARG BUILDPLATFORM
ARG TARGETPLATFORM

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /server cmd/api/main.go

FROM alpine

EXPOSE 80

WORKDIR /app

COPY --from=builder /server /app/server

ENTRYPOINT [ "/app/server" ]