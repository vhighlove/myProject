FROM golang:latest

COPY . .
EXPOSE 8080
RUN go build -o main .
CMD ["./main"]