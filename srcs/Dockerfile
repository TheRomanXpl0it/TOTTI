FROM golang:1.22-alpine

COPY . /app/

WORKDIR /app

#RUN apk add --no-cache git make build-base libpcap-dev

RUN cd /app/ && go build ./

# PRODUCTION
ENTRYPOINT [ "/app/sub", "-c", "/app/config.yml" ]

# DEVELOPMENT
# ENTRYPOINT [ "tail", "-f", "/dev/null" ]
