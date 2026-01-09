#!/usr/bin/env python3
import argparse
import json
import subprocess
import sys
import urllib.error
import urllib.request


def run_whois(domain):
    try:
        proc = subprocess.run(["whois", domain], capture_output=True, text=True, timeout=20)
    except FileNotFoundError:
        return "", "whois not found"
    except subprocess.TimeoutExpired:
        return "", "whois timeout"
    if proc.returncode != 0:
        return "", proc.stderr.strip()
    return proc.stdout.strip(), ""


def fetch_crtsh(domain):
    url = f"https://crt.sh/?q=%25.{domain}&output=json"
    req = urllib.request.Request(url, headers={"User-Agent": "ct-osint"})
    with urllib.request.urlopen(req, timeout=20) as resp:
        data = resp.read().decode("utf-8")
    return json.loads(data)


def extract_subdomains(entries, domain):
    subdomains = set()
    suffix = f".{domain}"
    for entry in entries:
        name_value = entry.get("name_value", "")
        for name in name_value.splitlines():
            name = name.strip().lower()
            if name.endswith(suffix) or name == domain:
                subdomains.add(name)
    return sorted(subdomains)


def main():
    parser = argparse.ArgumentParser(description="OSINT enrichment for a domain")
    parser.add_argument("domain", help="domain name")
    args = parser.parse_args()

    whois_text, whois_err = run_whois(args.domain)
    if whois_err:
        print(f"whois error: {whois_err}", file=sys.stderr)

    subdomains = []
    try:
        entries = fetch_crtsh(args.domain)
        subdomains = extract_subdomains(entries, args.domain)
    except (urllib.error.URLError, json.JSONDecodeError) as exc:
        print(f"crt.sh query failed: {exc}", file=sys.stderr)

    if whois_text:
        print("WHOIS:")
        print(whois_text)
        print("")

    if subdomains:
        print("SUBDOMAINS:")
        for name in subdomains:
            print(name)

    return 0


if __name__ == "__main__":
    sys.exit(main())


# Signed-off-by: ronikoz
