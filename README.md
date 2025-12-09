# mini-gateway

A simple Go HTTP server with dynamic configuration capabilities, inspired by Nginx.

## Features

- Dynamic configuration blocks
- Support for multiple server instances

## Installation

To install mini-gateway, use the following command:

```bash
go install github.com/Bean-jun/mini-gateway
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
- name: server1
  port: 7256
  protocol: http
- name: server2
  port: 7257
  protocol: tcp
- name: server3
  port: 7258
  protocol: http
```

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.
