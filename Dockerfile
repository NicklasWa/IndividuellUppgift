FROM golang:alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY *.yml ./
RUN go mod download
COPY database/gorm.sqlite database/sqlite
COPY database/gorm.sqlite /database/sqlite

COPY *.go ./
COPY data/*.go ./data/
COPY templates/* ./templates/

RUN go build -o /applikationen
RUN mkdir -p /database


EXPOSE 8080

CMD [ "/applikationen" ]