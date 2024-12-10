vm.wait_for_unit("multi-user.target")
print(vm.succeed("ls -all /mnt/package"))
vm.succeed("sudo dpkg -i /mnt/package/*amd64.deb")
vm.succeed("exit 0")
