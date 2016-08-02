FROM flavioribeiro/snickers-docker:v3

# Download snickers
RUN go get -u github.com/snickers/snickers

# Run snickers!
RUN curl -O http://flv.io/snickers/config.json
RUN go install github.com/snickers/snickers
ENTRYPOINT snickers
EXPOSE 8000
