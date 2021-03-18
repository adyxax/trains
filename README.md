# Trains

Trains is a simple web app to display train timetables for specific lines at stations on France SNCF's network. It records weekly passages and you can subscribe to certain train times to alert when schedules change or your stop is removed. It queries the SNCF official api and present the results in a minimal web page that loads fast (unlike the official sites with all their images and ads).

## Content

- [Dependencies](#dependencies)
- [Building](#building)
- [Usage](#usage)

## Dependencies

go is required. Only go version >= 1.16 on linux amd64 (Gentoo and Ubuntu 20.04) and on OpenBSD amd64 has been tested.

## Building

To run tests, use :
```
go test -cover ./...
```

For a debug build, use :
```
go build
```

For a release build, use :
```
go build -ldflags="-s -w"
```

## Usage

TODO
