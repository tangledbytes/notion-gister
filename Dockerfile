FROM golang:1.19.3-alpine3.16 AS builder
WORKDIR /gister
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/gister .
RUN apk add -U --no-cache ca-certificates
RUN apk add -U --no-cache tzdata

FROM scratch
WORKDIR /gister
COPY --from=builder /gister/bin/gister ./gister
COPY --from=builder /gister/.gister.yaml ./.gister.yaml
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
ENTRYPOINT ["/gister/gister"]