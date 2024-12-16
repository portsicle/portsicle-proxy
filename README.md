## Attorney Proxy Server and Toolkit

Attorney is a lightweight and flexible HTTP/HTTPS forward proxy server implemented in Go.

The Toolkit builds on top of the core Attorney server to provide additional features to enhance its functionality.

## Installation guide

1. Install the provided binary from latest release.

2. Give executeable permission to the binary `chmod +x ./attonery-toolkit`.

3. Use the CLI to run the server locally:

```zsh
attorney-toolkit run -p 8888
```

4. Use help to know more:

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

1. Run the executable binary with appropriate flags.
2. Add `http://127.0.0.1` and the port specified you've specified on your system's proxy configuration for both HTTP and HTTPS Proxy.
3. Make sure `127.0.0.1` is not listed in your machine's Ignored Hosts list.

---

### Features:

- Site Blocking:

  Blocks your access to the sites you add in your block list.

  Creates a **sqlite** database file `blocked_domains.db` in your directory to store your blocked domains.

  1. Add a site to blocklist:

  ```
  attorney-toolkit block --add https://example.com
  # or
  attorney-toolkit block --add someexample.com
  # or
  attorney-toolkit block -a other.example.com
  ```

  2. Remove a site from blocklist:

  ```
  attorney-toolkit block --remove https://example.com
  # or
  attorney-toolkit block -r someexample.com
  ```

---

### Future Scope:

- Ad / Content Blocker.
- Offline Browsing Mode.

## How It Works

1. **Proxying HTTPS Traffic**: The proxy listens for `CONNECT` requests and establishes a TCP tunnel to the target server. Once the tunnel is established, it enables bidirectional data transfer between the client and the destination.

2. **Request Handling**: Hijacks HTTP connections to bypass `HTTP/1.1` handling and directly interact with TCP sockets. Provides HTTP 200 response upon successful tunnel establishment.

## Contributing

Contributions are welcome! Feel free to open issues, submit pull requests, or suggest new features.
