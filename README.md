# mini-gateway

A simple Go HTTP server with dynamic configuration capabilities, inspired by Nginx.

## Features

- Dynamic configuration blocks
- Support for multiple server instances

## Installation

To install mini-gateway, use the following command:

```bash
go install github.com/Bean-jun/mini-gateway@latest
```

## Usage

To run the mini-gateway, execute the following command:

```bash
go run main.go
```

## Configuration

The mini-gateway configuration is defined in a YAML file. Below is an example configuration:

```yaml
server_blocks:
  - name: server # server name
    protocol: http  # protocol type
    port: 7256  # listen port
    max_body_size: 4 # max request body size unit MB
    ssl:  # ssl configuration
      cert_file: /path/to/cert.pem
      key_file: /path/to/key.pem
    locations:  # router configuration
      - path: /api/.+ # router path
        proxy_pass: # reverse proxy configuration
          - schema: http
            host: 127.0.0.1
            port: 5000
            weight: 1 # weight
          - schema: http
            host: 127.0.0.1
            port: 5001
            weight: 2
      - path: /.* # router path
        root: html  # static file root directory
```

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.
