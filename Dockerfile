FROM golang:latest

ADD . /server
WORKDIR /server/src
RUN echo "Asia/Ho_Chi_Minh" > /etc/timezone


CMD ["go", "run", "main.go"]
