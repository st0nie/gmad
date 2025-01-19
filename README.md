# Garry's Mod Anti DDoS
Prevent DDoS attacks on your Garry's Mod server

## Motivation
The reason for writing this is that `gm_protect` does not support frp proxy.

## Installation
1. Build the Go program: `go build -o gmad main.go`
2. Run the executable: `./gmad --port=8080 --protected-port=9090`
3. Ensure `iptables` and `ipset` are installed on your system.

## Usage
- Access `http://server-ip:8080/` to automatically add the visitor's IP to the whitelist.

## Acknowledgements
Special thanks to [OverlordAkise](https://github.com/OverlordAkise/gmod-netprotect/) for his work on `gmod-netprotect`, which inspired this project.