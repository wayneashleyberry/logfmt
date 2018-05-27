> An opinionated log formatter that you probably don't want to use.

[![Go Report Card](https://goreportcard.com/badge/github.com/wayneashleyberry/logfmt)](https://goreportcard.com/report/github.com/wayneashleyberry/logfmt)

- `logfmt` reads from stdin and prints formatted logs
- `logfmt` expects structured json with certain keys
- `logfmt` formats output to be similar to Google Cloud Platform Logging


### Installation

```sh
go get -u github.com/wayneashleyberry/logfmt
```

### Usage

```sh
go run app.go | logfmt
```

### Example

<img width="1074" alt="screen shot 2018-05-27 at 11 55 58" src="https://user-images.githubusercontent.com/727262/40585113-ef77e266-61a4-11e8-9c1d-60307614b04c.png">
