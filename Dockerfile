FROM golang:1.19.4 as builder
WORKDIR /src/
COPY . .
RUN go build -o /tmp/webapi ./cmd/webapi/


FROM debian:stable
COPY --from=builder /tmp/webapi /bin/webapi

CMD ["/bin/webapi"]