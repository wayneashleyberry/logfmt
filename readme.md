> An opinionated log formatter that you probably don't want to use.

- `logfmt` reads from stdin and prints formatted logs
- `logfmt` expects structured json with certain keys
- `logfmt` formats output to be similar to Google Cloud Platform Logging

### Installation

You will need to have your `$PATH` setup to your go installation.

```
# Add this to your terminal .rc file
export PATH="$HOME/go/bin:$PATH"
```

```sh
go get -u github.com/overhq/logfmt
```

### Usage

`logfmt` reads from stdin, so pipe the output from your service into `logfmt`:

```sh
go run myapp.go | logfmt
```

### Example

```sh
cat testdata/test.json | go run main.go
```

<img width="1074" alt="screen shot 2018-05-27 at 12 27 47" src="https://user-images.githubusercontent.com/727262/40585374-67fe7a52-61a9-11e8-95a9-786df02f1913.png">
