# Malicious Online Rejection Traffic Interception System

> "I am fearless of being DDoSed."

Prevent DDoS attacks on your Garry's Mod server.

## Motivation
The reason for writing this is that `gmod-netprotect` does not support frp proxy.

## Installation
Assuming your gmod server is running on 27070/udp: 

1. Build the Go program: `go build -o mortis main.go`
2. Ensure `iptables` and `ipset` are installed on your system.
3. Run the executable: `mortis -port=28080 -limit=5/sec -burst=5 -protected-port=27070`

## Usage
- clone this repo into gmod addons directory
- modify `url` in [shared.lua](./lua/autorun/shared.lua)
- modify `sv_loadingurl` to `http://your.server.url/http://your.origin.loadingurl`

## Acknowledgements
Special thanks to [OverlordAkise](https://github.com/OverlordAkise/gmod-netprotect/) for his work on `gmod-netprotect`, which inspired this project.