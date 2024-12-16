## Attorney Proxy Server and Toolkit

Attorney is a lightweight and flexible HTTP/HTTPS forward proxy server implemented in Go.

The Toolkit builds on top of the core Attorney server to provide additional features to enhance its functionality.

#### Upcoming features :

- Site Access Management.
- Offline Browsing Mode.
- Ad / Content Blocker.

## Installation guide

1. Install the provided binary.

> recommended to use latest release

2. Use the CLI to run the server locally:

```zsh
attorney-toolkit run -p 8888
```

3. Use help to know more:

```
attorney-toolkit --help

Allows HTTP/HTTPS transparent proxying on your machine.

Usage:
  attorney-toolkit [command]

Available Commands:
  help        Help about any command
  run         Runs the proxy server

Use "attorney-toolkit [command] --help" for more information about a command.
```

## Usage

1. Run the source code or build.
2. Add `http://127.0.0.1` and the port specified on your system's proxy configuration for both HTTP and HTTPS Proxy.
3. Make sure `127.0.0.1` is not listed in your machine's Ignored Hosts list.

## How It Works

1. **Proxying HTTPS Traffic**: The proxy listens for `CONNECT` requests and establishes a TCP tunnel to the target server. Once the tunnel is established, it enables bidirectional data transfer between the client and the destination.

2. **Request Handlin**g: Hijacks HTTP connections to bypass `HTTP/1.1` handling and directly interact with TCP sockets. Provides HTTP 200 response upon successful tunnel establishment.

## Contributing

Contributions are welcome! Feel free to open issues, submit pull requests, or suggest new features.
