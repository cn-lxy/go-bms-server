FROM golang:alpine
WORKDIR /code
ADD . ./
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.io"
RUN go build -o bms .
EXPOSE 8080
ENTRYPOINT  ["./bms"]