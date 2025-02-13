vm.execute(
    """
  # Record the install command's output and timings.
  script --log-timing /tmp/typescript.timing --log-out /tmp/typescript --command 'curl -sL pkg.paretosecurity.com/install.sh | sudo bash'
"""
)

# vm.execute(
#     """
#   # Upload the recorded output to termbin.com using netcat on port 9999.
#   cat /tmp/typescript | nc termbin.com 9999 > /tmp/paste_url.txt
#   # Upload the timing file to termbin.com.
#   cat /tmp/typescript.timing | nc termbin.com 9999 > /tmp/paste_timing_url.txt
#   # Print the exit code (this line is used by the driver).
#   echo $rc
#   """
# )

# # Later in the test, retrieve the paste URLs.
# log_url = vm.succeed("cat /tmp/paste_url.txt")
# timing_url = vm.succeed("cat /tmp/paste_timing_url.txt")
# print("Log URL:", log_url)
# print("Timing URL:", timing_url)

# vm.execute("scriptreplay /tmp/typescript.timing /tmp/typescript")
