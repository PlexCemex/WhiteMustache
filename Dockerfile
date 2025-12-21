FROM golang:latest
WORKDIR /app
COPY main main
EXPOSE 8080
ENTRYPOINT [ "./main" ]