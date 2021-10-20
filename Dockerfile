FROM golang:latest

RUN apt-get -y update && apt-get -y install git

RUN mkdir /app

COPY ./app /app

WORKDIR /app

RUN go mod vendor
RUN go mod tidy

CMD bash -c /run.sh
