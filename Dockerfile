FROM golang:1.18.3

ENV GOOS linux
ENV CGO_ENABLED 0

COPY . /app

WORKDIR /app

ENTRYPOINT ["go", "test", "-v", "./integration/"]

