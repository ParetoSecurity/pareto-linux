vm.wait_for_unit("multi-user.target")
print(vm.succeed("ls -all /mnt/package"))
vm.succeed("sudo apt install /mnt/package/*amd64.deb")

res = vm.succeed("auditor check --json")
fail_count = res.count("fail")
assert fail_count == 0, f"Found {fail_count} failed checks"
