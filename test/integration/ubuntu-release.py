# Capture installation logs
vm.execute(
    "curl -sL pkg.paretosecurity.com/install.sh | sudo bash > /tmp/install.log 2>&1; echo $? > /tmp/install-exit-code"
)

# Show logs after installation
vm.succeed("cat /tmp/install.log")
vm.succeed("cat /tmp/install-exit-code")

# Check systemd logs
vm.succeed("journalctl -xeu pareto-linux.socket --no-pager > /tmp/socket.log")
vm.succeed("cat /tmp/socket.log")

# Run systemctl status and print result
res = vm.succeed("sudo systemctl status pareto-linux.socket --no-pager")
print(res)

# Check paretosecurity results
res = vm.succeed("/usr/bin/paretosecurity check --json")
print(res)

fail_count = res.count("fail")
assert fail_count == 0, f"Found {fail_count} failed checks"
