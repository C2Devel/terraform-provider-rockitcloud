import os
import re
from datetime import datetime
from pathlib import Path
from zoneinfo import ZoneInfo

version = os.environ.get("RELEASE_VERSION", "").removeprefix("v")
if not re.fullmatch(r"(?:0|[1-9][0-9]*)(?:\.(?:0|[1-9][0-9]*)){2}", version):
    raise SystemExit("Set RELEASE_VERSION to a version such as 26.2.7 or v26.2.7.")

changelog = Path("CHANGELOG.md")
heading, newline, entries = changelog.read_bytes().partition(b"\n")
unreleased = re.fullmatch(rb"## [0-9]+\.[0-9]+\.[0-9]+ \(Unreleased\)(\r?)", heading)
if not unreleased:
    raise SystemExit("The first changelog heading must be marked (Unreleased).")

release_prefix = f"## {version} (".encode()
if any(line.startswith(release_prefix) for line in entries.splitlines()):
    raise SystemExit(f"Version {version} already has a changelog section.")

release_date = datetime.now(ZoneInfo("Europe/Moscow")).strftime("%B %d, %Y")
release_heading = f"## {version} ({release_date})"
changelog.write_bytes(release_heading.encode() + unreleased[1] + newline + entries)
print(release_heading)

if output_path := os.environ.get("GITHUB_OUTPUT"):
    with open(output_path, "a", encoding="utf-8") as output:
        print(f"version={version}", file=output)
