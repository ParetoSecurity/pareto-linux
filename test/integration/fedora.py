vm.wait_for_unit("multi-user.target")
print(vm.succeed("ls -all /mnt/package"))
vm.succeed("sudo dnf install -y /mnt/package/paretosecurity_amd64.rpm")
res = vm.succeed("auditor check --json")
fail_count = res.count("fail")
assert fail_count == 0, f"Found {fail_count} failed checkas"
