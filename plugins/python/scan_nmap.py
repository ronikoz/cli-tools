#!/usr/bin/env python3
import argparse
import subprocess
import sys


def main():
    parser = argparse.ArgumentParser(description="Run an nmap scan")
    parser.add_argument("target", help="target host or CIDR")
    parser.add_argument("--ports", default="", help="comma-separated ports")
    args, extra = parser.parse_known_args()

    cmd = ["nmap", "-sV", args.target]
    if args.ports:
        cmd.extend(["-p", args.ports])
    cmd.extend(extra)

    try:
        print(" ".join(cmd))
        return subprocess.call(cmd)
    except FileNotFoundError:
        print("nmap not found in PATH", file=sys.stderr)
        return 127


if __name__ == "__main__":
    sys.exit(main())


# Signed-off-by: ronikoz
