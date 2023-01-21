FROM golang:1.19.2

WORKDIR /app
ADD . /app

CMD ["go", "run", "main.go"]