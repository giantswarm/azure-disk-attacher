FROM flynn/busybox

COPY azure-disk-attacher /azure-disk-attacher

ENTRYPOINT ["/azure-disk-attacher"]
