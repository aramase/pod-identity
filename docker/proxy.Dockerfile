FROM us.gcr.io/k8s-artifacts-prod/build-image/debian-iptables-amd64:v12.1.2
COPY ./_output/proxy /bin/
RUN chmod a+x /bin/proxy

ENTRYPOINT [ "proxy" ]
EXPOSE 8000
