FROM golang:1.18-bullseye as base

RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid 65532 \
  small-user

RUN go install github.com/Mic92/ssh-to-age/cmd/ssh-to-age@latest && cp $GOPATH/bin/ssh-to-age /ssh-to-age && chmod +x /ssh-to-age

FROM debian:bullseye-slim

COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group

COPY --from=base /ssh-to-age .

USER small-user:small-user

ENTRYPOINT ["./ssh-to-age"]
