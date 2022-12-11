# syntax=docker/dockerfile:1

FROM golang:1.18-alpine3.16
WORKDIR /app
# ARG DB_URL
# ENV DATABASE_URL ${DBURL}
COPY  . .

RUN go mod download


RUN go build -o  /main
EXPOSE 8000
CMD [ "/main" ]