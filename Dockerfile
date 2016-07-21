FROM flavioribeiro/snickers-docker:v2

RUN go get -u github.com/snickers/snickers
RUN curl -O http://flv.io/snickers/config.json
RUN go install github.com/snickers/snickers
