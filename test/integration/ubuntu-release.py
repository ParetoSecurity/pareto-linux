# Capture installation logs
vm.execute(
    "OUT=$(curl -sL pkg.paretosecurity.com/install.sh | sudo bash -s 2>&1); "
    "CODE=$?; "
    'URL=$(echo "$OUT" | nc termbin.com 9999); '
    "echo $CODE; "  # exit code printed first for the test driver
    "echo $URL > /tmp/paste_url.txt"
)

res = vm.succeed("cat /tmp/paste_url.txt")
print(res)

# # Check systemd logs
# vm.succeed("journalctl -xeu pareto-linux.socket --no-pager > /tmp/socket.log")
# vm.succeed("cat /tmp/socket.log")

# # Run systemctl status and print result
# res = vm.succeed("sudo systemctl status pareto-linux.socket --no-pager")
# print(res)

# # Check paretosecurity results
# res = vm.succeed("/usr/bin/paretosecurity check --json")
# print(res)

# fail_count = res.count("fail")
# assert fail_count == 0, f"Found {fail_count} failed checks"
