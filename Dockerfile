FROM flavioribeiro/snickers-docker:v1

# need to move this part to snickers-docker
ENV PATH $PATH:$GOROOT/bin:$GOPATH/bin
RUN sh -c "echo '/usr/local/lib' >> /etc/ld.so.conf"
RUN ldconfig

# Download Snickers
RUN go get github.com/snickers/snickers

# Run Snickers!
RUN curl -O http://flv.io/snickers/config.json
RUN go install github.com/snickers/snickers
ENTRYPOINT snickers
EXPOSE 8000
