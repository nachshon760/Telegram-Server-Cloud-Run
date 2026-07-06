# שלב 1: קומפילציה של הקוד ב-Go
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# שלב 2: הרצה על תמונה קלה ומאובטחת
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .

# הגדרת הפורט והרצה
EXPOSE 8080
CMD ["./server"]
