FROM alpine:latest as alpine

COPY glok /

RUN apk add -U --no-cache ca-certificates
RUN chmod +x /

FROM scratch
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=alpine /glok /glok

EXPOSE 8888

ENTRYPOINT ["/glok"]
