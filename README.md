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
Usage of ./wgmitmgui_0.1.0_linux_arm64:
  -b, --bind string                address:port to bind webserver (default "127.0.0.1:8081")
  -c, --client-ip string           Client IP that should be redirected (default "192.168.0.222")
  -x, --ip-header string           header for user IP, such as X-Real-IP or X-Forwarded-For - this is NOT used for security, it's only for displaying remote IP in the UI
  -i, --iptables-bin string        Path to iptables (default "/sbin/iptables")
  -s, --iptables-save-bin string   Path to iptables-save (default "/sbin/iptables-save")
      --version                    display version
```

![Jul-16-2022 00-48-26](https://user-images.githubusercontent.com/636320/179325250-4b8f2779-05ff-450d-acbf-2dc3a8a51b55.gif)
