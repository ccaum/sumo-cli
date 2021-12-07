FROM alpine:3.10

COPY sumo-linux-amd64 /sumo

ENTRYPOINT ["/sumo"]
