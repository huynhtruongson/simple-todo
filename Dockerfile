# build project
FROM golang:1.21.4-alpine3.18 AS builder

WORKDIR /app

COPY . .

RUN go build -o main main.go

# copy binary file from built project
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/main .
COPY app.env .
COPY wait-for.sh .
COPY start.sh .
COPY migration/ ./migration/

EXPOSE 3000 3001

CMD [ "/app/main" ]

ENTRYPOINT [ "/app/start.sh"]

