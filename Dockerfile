FROM scratch
ENTRYPOINT ["/heating-control-mqtt"]
COPY heating-control-mqtt config.yaml.example /
COPY config.yaml.example /config.yaml
