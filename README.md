# tunnel-guard

Tunnel-guard is a tool to keep ssh tunnels alive and restore them if they fail

Tunnel-guard uses /etc/tunnel-guard/tuns.conf as the central configuration script and looks like this by default:

`/etc/tunnel-guard/tuns.conf:`
```
# syntax: [name] [server address] [local port] [remote port]
# example for Matrix Synapse: matrix 192.168.0.44 8008 8008 # (192.168.0.44:8008 -> 127.0.0.1:8008)
# example using non-standard ports: matrix 192.168.0.44 1278 8972 # (192.168.0.44:8972 -> 127.0.0.1:1278)
# begin user confs:

# Placeholders:
placeholder1 127.0.0.1 22 9876
placeholder2 127.0.0.1 22 7896
# Test tunnel activity by using "ssh -p 9876 localhost" and "ssh -p 7896 localhost"
```

# install instructions:

`git clone https://codeberg.org/firebadnofire/tunnel-guard`

`cd tunnel-guard`

`make install`

# Troubleshooting

On some systems, it may be needed to set `AllowTcpForwarding` to `yes` in your sshd configs. Your /etc/ssh/sshd_config should include the following somewhere in it:

`AllowTcpForwarding yes`

Be sure it is not commented out, then restart sshd and tunnel-guard
