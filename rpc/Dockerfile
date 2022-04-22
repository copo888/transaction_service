FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
COPY bin/transaction /app/transaction
COPY rpc/etc /app/etc

CMD ["./transaction", "-f", "etc/transaction.yaml"]
