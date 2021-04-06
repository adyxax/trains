# Trains

Trains is a simple web app to display train timetables for specific lines at stations on France SNCF's network. It queries the SNCF official api by default but will work with any compatible Navitia api implementation and present the results in a minimal web page that loads fast (unlike the official sites with all their images and ads).

Api queries are cached for 60 seconds so that someone refreshing your instance cannot simply DOS the api and exhaust your request quota.

A personal instance runs at https://trains.adyxax.org/.

## Content

- [Dependencies](#dependencies)
- [Quick install](#quick-install)
- [Configuration](#configuration)
- [Usage](#usage)
- [Building](#building)
- [References](#references)

## Dependencies

go is required. Only go version >= 1.16 on linux amd64 (Gentoo and Ubuntu 20.04) and on OpenBSD amd64 has been tested.

## Quick Install

```
go install git.adyxax.org/adyxax/trains/cmd/trains-webui@latest
```

## Configuration

The default configuration file location is `$HOME/.config/trains/config.yaml`. It is a yaml configuration file that should look like the following :

```
address: 127.0.0.1
port: 8082
token: 12345678-9abc-def0-1234-56789abcdef0
trainStop: stop_area:SNCF:87723502
```

`address` can be any ipv4 or ipv6 address or a hostname that resolves to such address and defaults to `127.0.0.1`. `port` can be any valid tcp port number or service name and defaults to `8080`.

You can get a free token from the [official SNCF's website](https://www.digital.sncf.com/startup/api/token-developpeur) for up to 5000 requests per day.

For now you can get the trainStop codes by manually parsing queries like the following. It is quite horrible but with a bit of perseverance you will find the stop code you want :
```
for i in `seq 4`; do curl 'https://TOKEN@api.sncf.com/v1/coverage/sncf/stop_areas?count=1000&start_page='$i |jq > $i.json; done
```

## Usage

Launching the webui server is as simple as :
```
trains-webui -c /path/to/config/file.yaml
```

The server will then listen for requests on the specified hostname and port until interrupted or killed.

Please consider running it behind a reverse proxy, with https. Also even though the static assets are embedded in the program's binary and can be served from there, consider serving the static assets directly from the web server acting as the reverse proxy or a cdn.

## Building

To run tests, use :
```
go test -cover ./...
```

For a debug build, use :
```
go build ./cmd/trains-webui/
```

For a release build, use :
```
go build -ldflags="-s -w" ./cmd/trains-webui/
```

To cross-compile for another os or architecture, use :
```
GOOS=openbsd GOARCH=amd64 go build -ldflags="-s -w" ./cmd/trains-webui/
```

## References

- https://www.digital.sncf.com/startup/api
- http://doc.navitia.io/
