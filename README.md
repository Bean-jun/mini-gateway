# g-server

A simple Go HTTP server with dynamic configuration capabilities, inspired by Nginx.

## Features

- Dynamic configuration blocks
- Support for multiple server instances

## Installation

To install g-server, use the following command:

```bash
go install github.com/Bean-jun/g-server
```

## Usage

To run the g-server, execute the following command:

```bash
go run main.go
```

## Configuration

The g-server configuration is defined in a YAML file. Below is an example configuration:

```yaml
configs:
  - name: server1
    port: 7257
  - name: server2
    port: 7258
  - name: server3
    port: 7259
```

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.
