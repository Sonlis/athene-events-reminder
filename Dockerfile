FROM alpine as packages
RUN apk --update add ca-certificates
RUN apk add tzdata

FROM scratch
WORKDIR /opt/app

COPY --from=packages /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=packages /usr/share/zoneinfo /usr/share/zoneinfo
COPY server server

USER 1001
ENTRYPOINT ["./server"]
