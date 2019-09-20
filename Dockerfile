FROM golang:1.11-alpine
WORKDIR /app
COPY . .

CMD ["app"]