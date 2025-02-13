vm.wait_for_unit("multi-user.target")
print(vm.succeed("ls -all /mnt/package"))
vm.execute(
    "DEBIAN_FRONTEND=noninteractive sudo dpkg -i /mnt/package/paretosecurity_amd64.deb > /tmp/dpkg-install.log 2>&1; echo $? > /tmp/dpkg-exit-code"
)
vm.succeed("cat /tmp/dpkg-install.log")
vm.succeed("cat /tmp/dpkg-exit-code")

res = vm.succeed("paretosecurity check --json")
fail_count = res.count("fail")
assert fail_count == 0, f"Found {fail_count} failed checks"
