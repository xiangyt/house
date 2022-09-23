FROM alpine:3.9
WORKDIR /go/bin
ENV TZ=Asia/Shanghai
COPY release/house /go/bin
COPY config/config.json /go/bin/config/
EXPOSE 30030
ENTRYPOINT ["./house"]