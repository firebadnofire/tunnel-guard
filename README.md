# tunnel-guard

Tunnel-guard is a tool to keep ssh tunnels alive and restore them if they fail

Tunnel-guard uses /etc/tunnel-guard/tuns.conf as the central configuration script and looks like this by default:

`/etc/tunnel-guard/tuns.conf:`
```
# syntax: [name] [server address] [local port] [remote port]
# example for Matrix Synapse: matrix 192.168.0.44 8008 8008 # (192.168.0.44:8008 -> 127.0.0.1:8008)
# example using non-standard ports: matrix 192.168.0.44 1278 8972 # (192.168.0.44:8972 -> 127.0.0.1:1278)
# begin user confs:
```

# install instructions:

`git clone https://codeberg.org/firebadnofire/tunnel-guard`

`cd tunnel-guard`

`make install`

To uninstall, run `make uninstall`

# Troubleshooting

On some systems, it may be needed to set `AllowTcpForwarding` to `yes` in your sshd configs. Your /etc/ssh/sshd_config should include the following somewhere in it:

`AllowTcpForwarding yes`

Be sure it is not commented out, then restart sshd and tunnel-guard

# Additional information:

This project has been tested with the following:

```
x86_64 Debian 12.7

ARM64 Debian 12.7

x86_64 Ubuntu 22.04

x86_64 Alma Linux 9
```

An installation of [Go](https://go.dev/dl/) is **required**

By default, the program will check SSH tunnels every 15 minutes. You can change this check interval with the `-m` flag. Eg: `tunnel-guard -m 1` for every minute. Decimals are also accepted.

This project is licensed under the GNU Affero General Public License, version 3 (AGPLv3) license. To view what rights are granted/limited, please see [the project license file](https://codeberg.org/firebadnofire/tunnel-guard/src/branch/main/LICENSE) or the [license rights file](https://codeberg.org/firebadnofire/tunnel-guard/src/branch/main/LICENSE-rights.md)

Included is a script called `tg-add-ssh-key` and it's syntax is as follows:

`sudo tg-add-ssh-key "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIDYOS9zxV7Qm9Qlnkfzj5ebLhtE/cdWELF0BIZiEnHWQ root@server"`

This will insert the provided SSH public key into `/etc/tunnel-guard/.ssh/authorized_keys` giving that key access to the ssh-tun user for incoming tunnel connections.
