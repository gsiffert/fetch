# Fetch

Small command lines which can be used to fetch pages from the web.

## Usage

Build with Golang, it requires GCC to be installed:
```bash
$ make build-go
```

Download a web page:
```bash
$ ./fetch https://www.google.com
```

Retrieve the last metadata for a web page:
```bash
$ ./fetch --metadata https://www.google.com
```

## Usage with Docker

Build with Docker:
```bash
$ make build-docker
```

Download a web page:
```bash
$ docker compose run fetch https://www.google.com
```

Retrieve the last metadata for a web page:
```bash
$ docker compose run fetch --metadata https://www.google.com
```

