FROM golang:1.15-alpine as builder

WORKDIR /app

ENV CGO_ENABLED=0
ARG VERSION="undefined"

ARG PAT
RUN if [ -n "$PAT" ] ; then echo "[url \"https://$PAT:x-oauth-basic@github.com/\"] insteadOf = https://github.com/" >> /root/.gitconfig ; fi
COPY . .

RUN apk add make git tzdata && \
    apk --update add ca-certificates

RUN VERSION=${VERSION} make build

FROM golang:1.15-alpine
LABEL service="system-design"

RUN apk add ghostscript

COPY --from=builder /app/build/application /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

ENTRYPOINT ["/application"]