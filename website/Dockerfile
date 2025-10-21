FROM golang:1.25.3-alpine3.21

WORKDIR /app

COPY . .

RUN [ -f go.mod ] || go mod init website
RUN go get github.com/gin-gonic/gin@v1.10.0 github.com/go-sql-driver/mysql@v1.7.1
RUN go mod tidy

EXPOSE 8080

CMD ["go", "run", "./main.go"]
