out2 = vm.execute("sudo bash -c \"echo 'Created foo → bar.\n' >&2 && echo 'foo' \"")
print(out2)
