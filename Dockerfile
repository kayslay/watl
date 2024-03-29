FROM golang:latest

WORKDIR /usr/watl
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app cmd/http/main.go

FROM alpine:latest as runner
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /usr/watl/app .
COPY --from=0 /usr/watl/public public
RUN mkdir -p /data/img
EXPOSE 8000
CMD ["./app"]  