FROM BASE_IMAGE

WORKDIR /opt/file-transfer

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone && \
      mkdir -p /var/log/file-transfer

COPY file-transfer /opt/file-transfer/bin/

ENTRYPOINT ["/opt/file-transfer/bin/file-transfer"]
CMD ["-c", "/etc/file-transfer/file-transfer.yaml"]
