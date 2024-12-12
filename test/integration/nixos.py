assert "Pareto Security CLI" in machine.succeed("auditor --help")

res = machine.succeed("auditor check --json")
fail_count = res.count("fail")
assert fail_count == 0, f"Found {fail_count} failed checks"
