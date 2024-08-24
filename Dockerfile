FROM golang:alpine

RUN apk update && \
    apk add --no-cache \
        nginx

COPY . /app
WORKDIR /app

ENV HOME='/root'
ENV GOROOT=/usr/local/go
ENV GOPATH=$HOME/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

RUN go mod vendor && \
    go mod tidy

RUN go build -o mvc ./cmd/main.go

COPY nginx.conf /etc/nginx/sites-available/default

EXPOSE 5050

CMD ["./mvc"]