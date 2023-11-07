FROM golang:1.20.10-alpine3.17 AS build
WORKDIR /app
COPY . .
RUN go build -ldflags="-w -s" -o main cmd/server/main.go


FROM alpine:3.17
WORKDIR /app
COPY --from=build /app/main ./main
COPY configs/app.env ./configs/

EXPOSE 9090
EXPOSE 8080
CMD [ "/app/main" ]
