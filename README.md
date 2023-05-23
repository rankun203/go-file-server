# go-file-server [![Go](https://github.com/rankun203/go-file-server/actions/workflows/go.yml/badge.svg)](https://github.com/rankun203/go-file-server/actions/workflows/go.yml) [![docker pull rankun203/openconnect-proxy](https://img.shields.io/docker/v/rankun203/go-file-server?label=docker%20pull%20rankun203/go-file-server)](https://hub.docker.com/r/rankun203/go-file-server)

You can POST / GET / list files to this server.

## Features

1. upload file: POST http://localhost:3000/upload
2. download file: GET http://localhost:3000/download/file.ext
3. list files: GET http://localhost:3000/files

## User manual

Try to upload a file using the following command:

upload: 

```bash
curl -v -F "file=@/path/to/file.ext" http://localhost:3000/upload
curl -v -F "file=@/path/to/some.zip;filename=fancy.zip" http://localhost:3000/upload
```

download:

```bash
curl -O http://localhost:3000/download/file.ext
curl -O http://localhost:3000/download/fancy.zip
```

list files:

```bash
curl http://localhost:3000/files
```

Compress files:

```bash
tar -cf compressed.tar /path/to/directory
tar -cf - /path/to/directory | pigz > compressed.tar.gz
tar -czf compressed.tar.gz /path/to/directory
```

Decompress files:

```bash
tar -xf compressed.tar
tar -xzf compressed.tar.gz
```

## Docker

Build docker image

```bash
docker build -t go-file-server .
```

Run docker image

```bash
docker run -p 3000:3000 go-file-server

# or use existing server
docker run -p 3000:3000 rankun203/go-file-server
```

Visit http://localhost:3000.
