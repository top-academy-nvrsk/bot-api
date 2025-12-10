FROM golang:1.25-alpine AS builder 

RUN apk update && apk add --no-cache gcc musl-dev 

WORKDIR /app 

COPY go.mod go.sum ./ 

RUN go mod download 

COPY . . 

ENV CGO_ENABLED=1 

RUN go build -o my-api-app cmd/api/main.go 

FROM alpine:latest 

WORKDIR /app 

COPY --from=builder /app/my-api-app . 

RUN mkdir -p /app/cmd/api/db 
RUN chown -R 1000:1000 /app/cmd/api/db
USER 1000

EXPOSE 8080

CMD ["./my-api-app"]
