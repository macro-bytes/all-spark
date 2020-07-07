#!/usr/bin/python3
import sys

def write_env(vm_tags: str):
    env_variables = vm_tags.split(';')
    with open("/allspark/env.sh", "w") as fh:
        for variable in env_variables:
            k, v = variable.split(':', 1)
            fh.write(f"{k}={v}\n")

if __name__ == "__main__":
    write_env(sys.stdin.read().strip())