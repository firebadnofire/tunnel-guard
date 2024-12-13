# tunnel-guard

Tunnel-guard is a tool to keep ssh tunnels alive and restore them if they fail

Tunnel-guard uses /etc/tunnel-guard/tuns.conf as the central configuration script and looks like this by default:

`/etc/tunnel-guard/tuns.conf:`
```
# syntax: [name] [local port] [remote port] [remote ip] (optional: local ip)
# example for Matrix Synapse:
# matrix 8008 8008 192.168.0.44
# example using non-standard ports:
# matrix 1278 8972 192.168.0.44 0.0.0.0
# begin user confs:
```

# Install instructions:

`git clone https://codeberg.org/firebadnofire/tunnel-guard`

`cd tunnel-guard`

`make install`

To uninstall, run `make uninstall`

# Alternative install:

Head to [https://archuser.org/tunnel-guard/builds/](https://archuser.org/tunnel-guard/builds/) and grab the latest deb package or a pre-compiled .tar.gz

```
wget -O tunnel-guard.tar.gz https://archuser.org/tunnel-guard/builds/tunnel-guard.{COMMITID}.V{VERSION}.tar.gz

tar -xavf tunnel-guard.tar.gz

cd tunnel-guard

sudo ./aio.sh [install/uninstall]

```

At this point, the tar and `tunnel-guard` directory created by tar are no longer needed and can be removed by `rm -r tunnel-guard/ tunnel-guard.tar.gz`

You can install the deb package with `sudo apt install -y /path/to/tunnel-guard.deb`

# Troubleshooting

On some systems, it may be needed to set `AllowTcpForwarding` and `GatewayPorts` to `yes` in your sshd configs. Your /etc/ssh/sshd_config should include the following somewhere in it:

```
AllowTcpForwarding yes
GatewayPorts yes
```

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

If you want to change the SystemD service check interval, run `sudo nano /etc/systemd/system/tunnel-guard.service` then change `ExecStart=/usr/bin/tunnel-guard` to `ExecStart=/usr/bin/tunnel-guard -m 0.1` or whatever interval value you want

This project is licensed under the GNU Affero General Public License, version 3 (AGPLv3) license. To view what rights are granted/limited, please see [the project license file](https://codeberg.org/firebadnofire/tunnel-guard/src/branch/main/LICENSE) or the [license rights file](https://codeberg.org/firebadnofire/tunnel-guard/src/branch/main/LICENSE-rights.md)

Included is a script called `tg-add-ssh-key` and it's syntax is as follows:

`sudo tg-add-ssh-key "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIDYOS9zxV7Qm9Qlnkfzj5ebLhtE/cdWELF0BIZiEnHWQ root@server"`

You may also substitute the key for `list` to list the contents of `/etc/tunnel-guard/.ssh/authorized_keys`

Another script called `tg-transfer-key` is available to copy the ssh key automatically to a remote host assuming the user you access as sudo privileges. It's syntax is as follows:

`tg-transfer-key USERNAME HOST`

The `USERNAME` should be set to the user you want to ssh as. This will likely be your admin user account and should NOT be the `ssh-tun` user. The script reads `/etc/tunnel-guard/ssh-tun.pub` then remotes into an admin user's sudo power to append the key to `/etc/tunnel-guard/.ssh/authorized_keys` as the `ssh-tun` user

