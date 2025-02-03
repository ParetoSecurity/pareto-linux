vm.wait_for_unit("multi-user.target")
vm.succeed("curl -sL pkg.paretosecurity.com/install.sh | sudo bash")

res = vm.succeed("sudo systemctl status pareto-linux.socket --no-pager")
assert "active (running)" in res, "pareto-linux.socket is not running"

res = vm.succeed("/usr/bin/paretosecurity check --json")
fail_count = res.count("fail")
assert fail_count == 0, f"Found {fail_count} failed checks"
