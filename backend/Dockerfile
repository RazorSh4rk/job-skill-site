FROM golang

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o app .
EXPOSE 8080

CMD [ "./app" ]