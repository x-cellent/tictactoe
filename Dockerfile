FROM golang:1.15.3-buster as builder

ENV VERSION_PROTOC=3.13.0

# unzip
RUN apt update -qq \
 && apt install -y unzip

# protoc
RUN curl -fsSLO https://github.com/google/protobuf/releases/download/v${VERSION_PROTOC}/protoc-${VERSION_PROTOC}-linux-x86_64.zip \
 && unzip "protoc-${VERSION_PROTOC}-linux-x86_64.zip" -d protoc \
 && chmod -R o+rx protoc/ \
 && mv protoc/bin/* /usr/local/bin/ \
 && mv protoc/include/* /usr/local/include/ \
 && go get -u github.com/golang/protobuf/protoc-gen-go
