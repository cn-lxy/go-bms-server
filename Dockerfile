FROM golang:alpine
WORKDIR /home/lxy/GoProject/book_manager_system
ADD . ./
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.io"
RUN go build -o bms .
EXPOSE 8080
ENTRYPOINT  ["./bms"]