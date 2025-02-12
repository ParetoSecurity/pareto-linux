# Capture installation logs
  vm.succeed("""
    # Record the install command's output.
    script -q /tmp/typescript -c 'curl -sL pkg.paretosecurity.com/install.sh | sudo bash'
    rc=$?
    # Upload the recorded output to termbin.com using netcat on port 9999.
    cat /tmp/typescript | nc termbin.com 9999 > /tmp/paste_url.txt
    # Print the exit code (this line is used by the driver).
    echo $rc
  """)

# Later in the test, retrieve the paste URL.
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
