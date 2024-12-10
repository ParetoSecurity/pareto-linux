vm.wait_for_unit("multi-user.target")
print(vm.succeed("ls -all /mnt/package"))
vm.succeed("sudo dnf install -y /mnt/package/*x86_64.rpm")
vm.succeed("exit 0")
