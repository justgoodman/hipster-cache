FROM alpine:3.4 


RUN mkdir -p /go/bin/etc
ADD hipster-cache /go/bin/hipster-cache
RUN chmod +x /go/bin/hipster-cache
ADD application.json /go/bin/etc/application.json
# this will ideally be built by the ONBUILD below ;)
ENTRYPOINT ["/go/bin/hipster-cache","-config=/go/bin/etc/application.json"]
#ENTRYPOINT ["sleep","10000"]
#ONBUILD RUN go-wrapper download
#ONBUILD RUN go-wrapper install
