vm.wait_for_unit("multi-user.target")
print(vm.succeed("ls -all /mnt/package"))
vm.succeed(
    "DEBIAN_FRONTEND=noninteractive sudo dplg -i /mnt/package/paretosecurity_amd64.deb -y"
)

res = vm.succeed("paretosecurity check --json")
fail_count = res.count("fail")
assert fail_count == 0, f"Found {fail_count} failed checks"
