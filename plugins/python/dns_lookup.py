#!/usr/bin/env python3
import argparse
import subprocess
import sys


def parse_nslookup(lines, record_type):
    results = []
    key_map = {
        "A": "Address:",
        "AAAA": "Address:",
        "MX": "mail exchanger =",
        "NS": "nameserver =",
        "TXT": "text =",
    }
    key = key_map.get(record_type, "")
    for line in lines:
        line = line.strip()
        if not key:
            continue
        if key in line:
            value = line.split(key, 1)[1].strip()
            if value and value not in results:
                results.append(value)
    return results


def run_nslookup(domain, record_type):
    cmd = ["nslookup", f"-type={record_type}", domain]
    proc = subprocess.run(cmd, capture_output=True, text=True)
    if proc.returncode != 0:
        return [], proc.stderr.strip()
    lines = proc.stdout.splitlines()
    return parse_nslookup(lines, record_type), ""


def main():
    parser = argparse.ArgumentParser(description="DNS lookup utility")
    parser.add_argument("domain", help="domain name")
    parser.add_argument("--types", default="A,AAAA,MX,NS,TXT", help="comma-separated record types")
    args = parser.parse_args()

    record_types = [t.strip().upper() for t in args.types.split(",") if t.strip()]
    for record_type in record_types:
        values, err = run_nslookup(args.domain, record_type)
        if err:
            print(f"{record_type}: error: {err}", file=sys.stderr)
            continue
        if not values:
            print(f"{record_type}: (none)")
            continue
        for value in values:
            print(f"{record_type}: {value}")

    return 0


if __name__ == "__main__":
    sys.exit(main())
