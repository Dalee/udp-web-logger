FROM alpine:3.6

CMD ["/udp-web-logger"]
COPY ./udp-web-logger /udp-web-logger
