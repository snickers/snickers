<p align="center">
  <img src="https://cloud.githubusercontent.com/assets/244265/16903251/00c5f2f8-4c47-11e6-9f2c-9c86bb37f114.png">
</p>
<br><br>

[![Build Status](https://travis-ci.org/snickers/snickers.svg?branch=master)](https://travis-ci.org/snickers/snickers)
[![codecov](https://codecov.io/gh/snickers/snickers/branch/master/graph/badge.svg)](https://codecov.io/gh/snickers/snickers)
[![Go Report Card](https://goreportcard.com/badge/github.com/snickers/snickers)](https://goreportcard.com/report/github.com/snickers/snickers)
<br><br>
Snickers is an open source alternative to the existent cloud encoding services. It is a HTTP API that encode videos.

## Setting Up

First make sure you have [Go](https://golang.org/dl/) and [FFmpeg](http://ffmpeg.org/) with `--enable-shared` installed on your machine. If you don't know what this means, look at how the dependencies are being installed on our [Dockerfile](https://github.com/snickers/snickers-docker/blob/master/Dockerfile).

Download the dependencies:

```
$ make build
```

You can store presets and jobs on memory or [MongoDB](https://www.mongodb.com/). On your `config.json` file:

- For MongoDB, set `DATABASE_DRIVER: "mongo"` and `MONGODB_HOST: "your.mongo.host"`
- For memory, just set `DATABASE_DRIVER: "memory"` and you're good to go.

Please be aware that in case you use `memory`, Snickers will persist the data only while the application is running.

Run!

```
$ make run
```

## Running tests

Make sure you have [mediainfo](https://sourceforge.net/projects/mediainfo/) installed and a local instance of [MongoDB](https://github.com/mongodb/mongo) running.

```
$ make test
```

## Using the API

Check out the [Wiki](https://github.com/snickers/snickers/wiki/How-to-Use-the-API) to learn how to use the API.

## Contributing

1. Fork it
2. Create your feature branch: `git checkout -b my-awesome-new-feature`
3. Commit your changes: `git commit -m 'Add some awesome feature'`
4. Push to the branch: `git push origin my-awesome-new-feature`
5. Submit a pull request

## License

This code is under [Apache 2.0 License](https://github.com/snickers/snickers/blob/master/LICENSE).

