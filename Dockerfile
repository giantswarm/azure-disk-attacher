FROM alpine:3.13.5

COPY azure-disk-attacher /azure-disk-attacher

ENTRYPOINT ["/azure-disk-attacher"]
