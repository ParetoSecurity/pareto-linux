vm.wait_for_unit("multi-user.target")
print(vm.succeed("ls -all /mnt/package"))
vm.succeed("sudo apt install /mnt/package/*amd64.deb")
