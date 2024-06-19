FROM --platform=$BUILDPLATFORM node:20.14-alpine3.20 as node-builder

WORKDIR /app

COPY frontend/package.json frontend/package-lock.json ./

RUN npm install

COPY frontend/ .

RUN npm run build

FROM --platform=$BUILDPLATFORM golang:alpine AS builder
ARG BUILDPLATFORM
ARG TARGETPLATFORM

WORKDIR /app

RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates

COPY go.mod go.sum ./

RUN go mod download

COPY cmd ./cmd
COPY pkg ./pkg
COPY internal ./internal

ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-w -s" -o /server cmd/api/main.go

FROM scratch

EXPOSE 80

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=node-builder /app/dist/frontend/browser /app/wwwroot

COPY --from=builder /server /app/server

ENTRYPOINT [ "/app/server" ]