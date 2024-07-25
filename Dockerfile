# Stage 1: Build
FROM golang:1.22 AS builder

WORKDIR /workspace

COPY go.mod .

RUN go mod download \
    && CGO_ENABLED=1 GOOS=linux

COPY . .

RUN go build -v -gcflags="all=-N -l" -o app main.go

# Stage 2: Production
FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app

COPY --from=builder /workspace/app /app/app

EXPOSE 8080

CMD ["./app/app"]
