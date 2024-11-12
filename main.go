package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "os/exec"
    "os/user"
    "path/filepath"
    "strings"
    "syscall"
    "time"
)

type TunnelConfig struct {
    Name       string
    ServerAddr string
    LocalPort  string
    RemotePort string
}

func main() {
    interval := flag.Float64("m", 15.0, "Interval in minutes to check tunnels (accepts decimals up to 5 digits)")
    flag.Parse()

    checkRoot()

    tunsConfPath := "/etc/tunnel-guard/tuns.conf"

    if _, err := os.Stat(tunsConfPath); os.IsNotExist(err) {
        fmt.Println("First init detected, setting up...")
        setupTunnelGuard()
    }

    tunnels, err := readTunnelsConfig(tunsConfPath)
    if err != nil {
        log.Fatalf("Error reading tuns.conf: %v", err)
    }

    tunnelProcesses := make(map[string]*exec.Cmd)
    for _, tunnel := range tunnels {
        cmd, err := startTunnel(tunnel)
        if err != nil {
            log.Printf("Error starting tunnel %s: %v", tunnel.Name, err)
            continue
        }
        tunnelProcesses[tunnel.Name] = cmd
    }

    defer func() {
        log.Println("Program exiting. Killing all SSH tunnels...")
        for name, cmd := range tunnelProcesses {
            if cmd != nil && cmd.Process != nil {
                cmd.Process.Kill()
                log.Printf("Killed tunnel %s", name)
            }
        }
    }()

    intervalDuration := time.Duration(*interval * float64(time.Minute))

    for {
        time.Sleep(intervalDuration)
        for _, tunnel := range tunnels {

            if !checkPortInUse(tunnel.LocalPort) {
                log.Printf("Tunnel %s is not active on port %s, restarting...", tunnel.Name, tunnel.LocalPort)
                newCmd, err := startTunnel(tunnel)
                if err != nil {
                    log.Printf("Error restarting tunnel %s: %v", tunnel.Name, err)
                    continue
                }
                tunnelProcesses[tunnel.Name] = newCmd
            } else {
                log.Printf("Tunnel %s on port %s is active.", tunnel.Name, tunnel.LocalPort)
            }
        }
    }
}

func checkRoot() {
    if os.Geteuid() != 0 {
        log.Fatal("This script must be run as root.")
    }
}

func setupTunnelGuard() {
    tunsConfPath := "/etc/tunnel-guard/tuns.conf"
    tunnelGuardDir := "/etc/tunnel-guard/"
    sshDir := filepath.Join(tunnelGuardDir, ".ssh")

    os.MkdirAll(tunnelGuardDir, 0755)

    sampleConf := `# syntax: [name] [server address] [local port] [remote port]
# example for Matrix Synapse: 
# matrix 192.168.0.44 8008 8008
# example using non-standard ports: 
# matrix 192.168.0.44 1278 8972
# begin user confs:
`
    ioutil.WriteFile(tunsConfPath, []byte(sampleConf), 0644)

    _, err := user.Lookup("ssh-tun")
    if err != nil {
        cmd := exec.Command("useradd", "-d", tunnelGuardDir, "ssh-tun")
        if err := cmd.Run(); err != nil {
            log.Fatalf("Error creating user ssh-tun: %v", err)
        }
    }

    os.MkdirAll(sshDir, 0700)
    keyPath := filepath.Join(sshDir, "id_ed25519")
    if _, err := os.Stat(keyPath); os.IsNotExist(err) {
        cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-N", "", "-f", keyPath)
        if err := cmd.Run(); err != nil {
            log.Fatalf("Error generating SSH key: %v", err)
        }
    }
    os.Chmod(keyPath, 0600)
    pubKeyPath := keyPath + ".pub"
    symlinkPath := filepath.Join(tunnelGuardDir, "ssh-tun.pub")
    if _, err := os.Stat(symlinkPath); os.IsNotExist(err) {
        os.Symlink(pubKeyPath, symlinkPath)
    }
    os.Chmod(symlinkPath, 0644)
    authorizedKeysPath := filepath.Join(sshDir, "authorized_keys")
    pubKeyData, err := ioutil.ReadFile(pubKeyPath)
    if err != nil {
        log.Fatalf("Error reading public key: %v", err)
    }
    ioutil.WriteFile(authorizedKeysPath, pubKeyData, 0600)

    chownCmd := exec.Command("chown", "-R", "ssh-tun:ssh-tun", sshDir)
    if err = chownCmd.Run(); err != nil {
        log.Fatalf("Error changing ownership of .ssh directory: %v", err)
    }

    fmt.Println("Setup complete. Public key generated at:", symlinkPath)
    fmt.Println("Please install this public key on the remote machine(s) in the ssh-tun user's authorized_keys file.")
    fmt.Printf("You can use the following command:\n")
    fmt.Printf("sudo ssh-copy-id -f -i %s ssh-tun@dest\n", pubKeyPath)
}

func readTunnelsConfig(confPath string) ([]TunnelConfig, error) {
    data, err := ioutil.ReadFile(confPath)
    if err != nil {
        return nil, err
    }
    lines := strings.Split(string(data), "\n")
    var tunnels []TunnelConfig
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        fields := strings.Fields(line)
        if len(fields) < 4 {
            log.Printf("Skipping invalid line in config: %s", line)
            continue
        }
        tunnel := TunnelConfig{
            Name:       fields[0],
            ServerAddr: fields[1],
            LocalPort:  fields[2],
            RemotePort: fields[3],
        }
        tunnels = append(tunnels, tunnel)
    }
    return tunnels, nil
}

func startTunnel(tunnel TunnelConfig) (*exec.Cmd, error) {
    keyPath := "/etc/tunnel-guard/.ssh/id_ed25519"
    sshCmd := "ssh"
    localForward := fmt.Sprintf("%s:localhost:%s", tunnel.LocalPort, tunnel.RemotePort)
    args := []string{
        "-i", keyPath,
        "-o", "StrictHostKeyChecking=no",
        "-o", "UserKnownHostsFile=/dev/null",
        "-N",
        "-L", localForward,
        "ssh-tun@" + tunnel.ServerAddr,
    }
    cmd := exec.Command(sshCmd, args...)
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Setsid: true,
    }
    err := cmd.Start()
    if err != nil {
        return nil, err
    }
    log.Printf("Started tunnel %s: %s", tunnel.Name, strings.Join(cmd.Args, " "))
    return cmd, nil
}

func checkPortInUse(port string) bool {
    cmd := exec.Command("lsof", "-i", fmt.Sprintf(":%s", port))
    output, err := cmd.Output()
    if err != nil {
        log.Printf("Error checking port %s: %v", port, err)
        return false
    }
    return strings.Contains(string(output), "root")
}
