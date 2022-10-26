# Build stage
FROM golang:1.18.6-alpine3.16 as builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY db/migration ./db/migration
COPY app.env .
COPY start.sh .
COPY wait-for .

EXPOSE 8080

CMD [ "/app/main" ]

ENTRYPOINT [ "/app/start.sh" ]