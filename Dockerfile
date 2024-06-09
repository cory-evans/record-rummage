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

COPY --from=builder /server /server
