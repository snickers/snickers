<p align="center">
  <img src="https://cloud.githubusercontent.com/assets/244265/16285772/292ffa0e-38a6-11e6-922e-a96f98699c63.png">
</p>
<br><br>
[![Build Status](https://travis-ci.org/snickers/snickers.svg?branch=master)](https://travis-ci.org/snickers/snickers)
[![codecov](https://codecov.io/gh/snickers/snickers/branch/master/graph/badge.svg)](https://codecov.io/gh/snickers/snickers)
[![Go Report Card](https://goreportcard.com/badge/github.com/snickers/snickers)](https://goreportcard.com/report/github.com/snickers/snickers)
<br><br>
Snickers is an open source alternative to the existent cloud encoding services. It is a HTTP API that encode videos.

## Setting Up

First make sure you have [Go](https://golang.org/dl/) and [FFmpeg](http://ffmpeg.org/) with `--enable-shared` installed on your machine. If you don't know what this means, follow the [instructions](https://github.com/3d0c/gmf#build).

Download the dependencies:

```
$ make build
```

Run!

```
$ make run
```

## Running tests

```
$ make test
```

## Using the API

Check out the [Wiki](https://github.com/flavioribeiro/snickers/wiki/How-to-Use-the-API) to learn how to use the API.

## Contributing

1. Fork it
2. Create your feature branch: `git checkout -b my-awesome-new-feature`
3. Commit your changes: `git commit -m 'Add some awesome feature'`
4. Push to the branch: `git push origin my-awesome-new-feature`
5. Submit a pull request

## License

This code is under [Apache 2.0 License](https://github.com/flavioribeiro/snickers/blob/master/LICENSE). 

