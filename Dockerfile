FROM golang:1.16
WORKDIR /go/src/app
COPY . .
RUN cd ./api && go build

FROM debian:10.8
COPY --from=0 /go/src/app/. .
# docker-compose-wait lets us wait for MongoDB to be responsive before launching the application with Docker Compose
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.3/wait /wait
RUN chmod +x ./wait && chmod +x ./api/api
EXPOSE 8080
CMD ["./api/api"]

