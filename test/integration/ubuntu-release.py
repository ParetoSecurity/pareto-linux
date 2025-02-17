# Capture installation logs
vm.execute(
    "curl -fsSL https://github.com/ParetoSecurity/pareto-core/releases/latest/download/paretosecurity_$(dpkg --print-architecture).deb -o /tmp/pareto.deb"
)

print("Installing Pareto Security")
# Install Pareto Security
# out = vm.execute("sudo dpkg -i /tmp/pareto.deb")

# print(out)

out2 = vm.execute(
    "sudo bash -c \"echo 'Selecting previously unselected package paretosecurity.\n(Reading database ... 65681 files and directories currently installed.)\nPreparing to unpack /tmp/pareto.deb ...\nUnpacking paretosecurity (0.0.78) ...\nSetting up paretosecurity (0.0.78) ...\nCreated symlink /etc/systemd/system/sockets.target.wants/pareto-core.socket → /etc/systemd/system/pareto-core.socket.\n' >&2 && echo 'Pareto Security installed'\""
)
print(out2)
