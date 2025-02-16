vm.execute("(head -c 4095 </dev/zero | tr '\0' 'X'; echo) 1>&2 && echo 'FOO_OOO'")
