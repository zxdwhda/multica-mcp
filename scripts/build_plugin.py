#!/usr/bin/env python3
"""Build the existing private app-bound plugin and immutable MCP resources."""

import argparse
import hashlib
import json
import re
import zipfile
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SOURCE = ROOT / "plugin-src"
PACKAGE = ROOT / "plugins/multica-mcp"
BUNDLE = ROOT / "internal/skillbundle/bundle.json"


def encode(value):
    return (json.dumps(value, ensure_ascii=False, indent=2, sort_keys=True) + "\n").encode()


def outputs():
    legacy = json.loads((SOURCE / ".codex-plugin/plugin.json").read_text())
    binding = json.loads((SOURCE / ".app.json").read_text())
    if legacy["name"] != "dev-6aa8f019b2008191a824208a7681da62":
        raise ValueError("Existing plugin identity changed")
    if binding != {"apps": {legacy["name"]: {"id": "asdk_app_6aa8f019b2008191a824208a7681da62"}}}:
        raise ValueError("Existing app binding changed")
    if len(legacy["interface"]["shortDescription"]) > 30:
        raise ValueError("Plugin subtitle exceeds 30 characters")
    manifest = {
        "$schema": "https://agent-plugins.org/schemas/1.0.0/plugin.schema.json",
        **{k: legacy[k] for k in ["name", "version", "description", "author"]},
        "extensions": {"com.openai": {"apps": "./.app.json", "interface": legacy["interface"]}},
    }
    files = {"plugin.json": encode(manifest), ".app.json": encode(binding), ".codex-plugin/plugin.json": encode(legacy)}
    resources = {}
    skills = []
    for folder in sorted((SOURCE / "skills").iterdir()):
        if not folder.is_dir() or folder.is_symlink():
            raise ValueError("Skill sources must be real directories")
        entry = (folder / "SKILL.md").read_text()
        match = re.match(r"---\nname: ([^\n]+)\ndescription: ([^\n]+)\n---\n", entry)
        if not match or match[1] != folder.name:
            raise ValueError(f"Invalid frontmatter in {folder.name}")
        local_files = {"SKILL.md": entry.encode()}
        local_files["references/execution.md"] = (SOURCE / "references/execution.md").read_bytes()
        if "references/comments.md" in entry:
            local_files["references/comments.md"] = (SOURCE / "references/comments.md").read_bytes()
        skill_resources = []
        for relative, data in sorted(local_files.items()):
            text = data.decode()
            for link in re.findall(r"\]\(([^)]+)\)", text):
                if ":" not in link and not link.startswith("#"):
                    resolved = (Path(relative).parent / link).as_posix()
                    if resolved not in local_files:
                        raise ValueError(f"Missing reference {folder.name}/{resolved}")
            uri = f"skill://multica-mcp/{folder.name}/{relative}"
            files[f"skills/{folder.name}/{relative}"] = data
            resources[uri] = {"name": f"{folder.name}/{relative}", "mimeType": "text/markdown", "text": text}
            skill_resources.append({"uri": uri, "digest": "sha256:" + hashlib.sha256(data).hexdigest()})
        skills.append({"uri": f"skill://multica-mcp/{folder.name}/SKILL.md", "name": match[1], "description": match[2], "resources": skill_resources})
    if len(skills) != 5:
        raise ValueError("Expected all five workflow skills")
    index = {"version": manifest["version"], "skills": skills}
    resources["skill://multica-mcp/index.json"] = {"name": "Multica Skills index", "mimeType": "application/json", "text": encode(index).decode()}
    return files, encode(resources), manifest["name"]


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--check", action="store_true", help="Validate generated files without modifying them")
    parser.add_argument("--archive", type=Path, help="Write an upload ZIP outside the package directory")
    args = parser.parse_args()
    files, bundle, package_name = outputs()
    expected = {PACKAGE / name: data for name, data in files.items()}
    expected[BUNDLE] = bundle
    if args.check:
        for path, data in expected.items():
            if not path.exists() or path.read_bytes() != data:
                raise ValueError(f"Stale generated file: {path.relative_to(ROOT)}")
    else:
        for path, data in expected.items():
            path.parent.mkdir(parents=True, exist_ok=True)
            path.write_bytes(data)
    actual = {p for p in PACKAGE.rglob("*") if p.is_file()}
    if actual != {PACKAGE / name for name in files}:
        raise ValueError("Unexpected file in plugin package")
    if args.archive:
        if args.archive.resolve().is_relative_to(PACKAGE.resolve()):
            raise ValueError("Archive must be outside package directory")
        args.archive.parent.mkdir(parents=True, exist_ok=True)
        with zipfile.ZipFile(args.archive, "w", zipfile.ZIP_DEFLATED) as archive:
            for name, data in sorted(files.items()):
                info = zipfile.ZipInfo(f"{package_name}/{name}", (2026, 10, 1, 0, 0, 0))
                info.compress_type = zipfile.ZIP_DEFLATED
                info.external_attr = 0o100644 << 16
                archive.writestr(info, data)
    print(f"Validated 5 Skills, {len(json.loads(bundle))} MCP resources, app binding preserved")


if __name__ == "__main__":
    main()
