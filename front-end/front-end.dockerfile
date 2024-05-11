FROM alpine:latest

RUN mkdir /app

COPY frontEndApp /app
COPY cmd/web/templates /templates

CMD [ "/app/frontEndApp" ]
