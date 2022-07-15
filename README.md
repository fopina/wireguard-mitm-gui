# wireguard-mitm-gui
WireGuard MITM GUI

Quickly redirect a predefined client (IP) through a transparent proxy (or disable it).

Possible use case (mine):

* Custom wireguard profile on mobile phone with dedicated IP (different than normal profile)
* Wireguard running on a server (such as home raspberry pi)
* Need to run `iptables ... -j REDIRECT MY_WORKSTATION_IP:8080` every time my laptop IP changes (or when using different machine for mitmproxy/Burp/whatever)

This small GUI allows to quick change `MY_WORKSTATION_IP` (and port) or disable redirection entirely.

### Usage

```
Usage of wgmitmgui:
 ...
```
