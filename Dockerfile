FROM golang:latest

ADD ./src /server
WORKDIR /server
RUN echo "Asia/Ho_Chi_Minh" > /etc/timezone
RUN ls
CMD ["go", "run", "main.go"]
