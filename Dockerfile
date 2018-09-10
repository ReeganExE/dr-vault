FROM golang:1.10 AS builder

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && \
  chmod +x /usr/local/bin/dep

WORKDIR /go/src/github.com/reeganexe/dr-vault
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dr-vault .

FROM vault:0.11.0
RUN apk add supervisor

COPY bundle/supervisord.conf /etc/
COPY bundle/kv1.sh /
COPY --from=builder /go/src/github.com/reeganexe/dr-vault/dr-vault /

VOLUME ["/var/source"]

USER vault

ENTRYPOINT [ "/usr/bin/supervisord" ]
