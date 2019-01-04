FROM alpine:latest as alpine

RUN apk add -U --no-cache ca-certificates

FROM scratch
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY glok /

USER 7200:7200
EXPOSE 8888

ENTRYPOINT ["/glok"]
