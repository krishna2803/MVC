FROM alpine:latest

RUN [ "env" ]

RUN apk update && \
    apk add --no-cache \
        su-exec \
        curl \
        build-base \
        postgresql \
        nginx

# yahan pe dockerignore se kuchh packages ignore karne hain
COPY . /app
WORKDIR /app

RUN wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz && \
    tar xvf go1.23.0.linux-amd64.tar.gz && \
    mv go /usr/local

# fir environment setup
ENV HOME='/root'
ENV GOROOT=/usr/local/go
ENV GOPATH=$HOME/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

# fir go packages install
RUN go mod vendor && \
    go mod tidy

# postgresql setup
RUN chmod +x ./psql-setup.sh && ./psql-setup.sh

# build 
RUN go build -o mvc ./cmd/main.go

# nginx vhost setup
COPY nginx.conf /etc/nginx/sites-available/default

EXPOSE 5050

# start nginx
CMD ["ash", "-c", "su-exec \
      postgres pg_ctl start -D /var/lib/postgresql/data -l /var/lib/postgresql/logfile.log && \
      /bin/ash"]
