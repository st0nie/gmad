# Malicious Online Rejection Traffic Interception System

> "I am fearless of being DDoSed."

Prevent DDoS attacks on your Garry's Mod server.

## Motivation
The reason for writing this is that `gmod-netprotect` does not support frp proxy.

## Installation
1. Build the Go program: `go build -o mortis main.go`
2. Ensure `iptables` and `ipset` are installed on your system.
3. Run the executable: `mortis -port=28080 -limit=5/sec -burst=5 -protected-port=27070`

## Usage
- Access `http://server-ip:28080/` to automatically add the visitor's IP to the whitelist.

## Acknowledgements
Special thanks to [OverlordAkise](https://github.com/OverlordAkise/gmod-netprotect/) for his work on `gmod-netprotect`, which inspired this project.