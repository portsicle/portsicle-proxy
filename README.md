## Attorney Proxy Server

Attorney is a lightweight and flexible HTTP/HTTPS forward proxy server implemented in Go.

This project serves as a learning tool for understanding the intricacies of proxy servers and the HTTP CONNECT method.

## Installation

Clone the repository:

```
git clone https://github.com/amitsuthar69/attorney.git
cd attorney
```

Install **Air** for live reload:

```
go install github.com/air-verse/air@latest
```

### Run the server:

```
air
```

or

```
go build main.go
./main
```

## Usage

1. Run the source code or build.
2. Add `http://127.0.0.1` on port `:8888` on your system's proxy configuration for both HTTP and HTTPS Proxy.
3. Make sure `127.0.0.1` is not listed in your machine's Ignored Hosts list.

## How It Works

1. **Proxying HTTPS Traffic**: The proxy listens for `CONNECT` requests and establishes a TCP tunnel to the target server. Once the tunnel is established, it enables bidirectional data transfer between the client and the destination.

2. **Request Handlin**g: Hijacks HTTP connections to bypass `HTTP/1.1` handling and directly interact with TCP sockets. Provides HTTP 200 response upon successful tunnel establishment.

## Contributing

Contributions are welcome! Feel free to open issues, submit pull requests, or suggest new features.
