FROM alpine:3.4 


RUN mkdir -p /go/bin/etc
ADD hipster-cache-proxy /go/bin/hipster-cache-proxy
RUN chmod +x /go/bin/hipster-cache-proxy
ADD application.json /go/bin/etc/application.json

ENTRYPOINT ["/go/bin/hipster-cache-proxy","-config=/go/bin/etc/application.json"]
#ENTRYPOINT (/go/bin/hipster-cache-proxy -config=/go/bin/etc/application.json 2>&1) > /var/log/hipster-cache.log
#ENTRYPOINT ["(/go/bin/hipster-cache-proxy","-config=/go/bin/etc/application.json","2>&1)",">", "/var/log/hipster-cache.log"]
#ENTRYPOINT ["sleep","1000"]
