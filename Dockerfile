FROM golang:latest
ENV GO111MODULE=on

LABEL  maintainer="@chisty <chisty.sust@gmail.com>"

WORKDIR /app


COPY . .

RUN go mod download

EXPOSE 9000

RUN go build

RUN go test -v ./...

CMD ["./miniurl"]