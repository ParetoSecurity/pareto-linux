vm.wait_for_unit("multi-user.target")
print(vm.succeed("ls -all /mnt/package"))
vm.succeed("exit 0")
