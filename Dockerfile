FROM golang:1.22-bullseye AS builder
COPY . .
RUN bash build.sh data-transmission /data-transmission .

FROM golang:1.22-bunllseye
COPY --from=builder /data-transmission /data-transmission
COPY data-transmission.toml /data-transmission.toml
COPY source.json /source.json
COPY model_config.json /model_config.json

# 调试页面静态页面
COPY ui/dist /ui/dist

COPY swagger /swagger
WORKDIR /
EXPOSE  80

ENTRYPOINT [ "/data-transmission" ]
CMD ["run"]
