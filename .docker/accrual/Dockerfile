FROM alpine:3.16

RUN apk --no-cache add libc6-compat

COPY cmd/accrual/accrual_linux_amd64 /bin

CMD ["/bin/accrual_linux_amd64"]