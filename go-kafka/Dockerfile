FROM golang:alpine as build-stage

#Install dev tools
RUN apk add --update --no-cache alpine-sdk bash ca-certificates \
      libressl tar git openssh openssl  openssl-dev \
      yajl-dev zlib-dev cyrus-sasl-dev build-base coreutils

# Build librdkafka
RUN cd /root && \
    git clone https://github.com/edenhill/librdkafka.git && \
    cd /root/librdkafka && \
    ./configure && \
    make -j "$(getconf _NPROCESSORS_ONLN)" && \
    make install

# For golang applications
#RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

# Get confluent-kafka-go/kafka library
RUN go get -u github.com/confluentinc/confluent-kafka-go/kafka

# Install go-randomdata for fake profile generator
RUN go get github.com/Pallinder/go-randomdata

WORKDIR /

COPY consumer.go .
COPY producer.go .

RUN go build -o consumer consumer.go
RUN go build -o producer producer.go
