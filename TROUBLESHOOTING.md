# Troubleshooting

## "codex binary not found in PATH"

Codex CLI is not installed or not in your PATH.

Solution:
- Install Codex CLI and ensure it’s on PATH.
- Verify with:
  `which codex`

## "Azure CLI not found" or not logged in

`az` is not installed or you’re not authenticated.

Solutions:
- Install Azure CLI: https://aka.ms/azcli
- Run: `az login`
 - Or choose Keychain auth in `codezure manage config` to run without Azure CLI

## "missing subscription/group/resource" errors

Configuration incomplete.

Solutions:
- Run interactive setup: `codezure manage config`
- Or set values: `codezure manage config set <key> <value>`

## "could not list deployments"

Your Azure OpenAI resource has no deployments or you lack permissions.

Solutions:
- Check resource in Azure Portal
- Create deployments in Azure Portal
- Ensure you have Azure OpenAI access and proper RBAC

## "cannot update development build"

You’re running a development build (`version dev`).

Solution:
- Install from releases or tag and download a release

## Installation Issues

- No internet: verify GitHub access
- Unsupported platform: check releases for your OS/arch
- Permission denied: use `sudo` for `/usr/local/bin`

Debug:
```
curl -fsSL https://raw.githubusercontent.com/OlaHulleberg/codezure/main/install.sh -o install.sh
bash -x install.sh
```

## Profile Issues

- Profile not found: `codezure manage profiles`
- Cannot delete current profile: switch first, then delete
- Migration: legacy `current.env` migrated to `profiles/default.json` on first run

## Keychain Authentication Issues

### "failed to store API key in keychain" / "failed to retrieve API key from keychain"

Codzure uses the OS keychain via `go-keyring`.

Solutions:
- macOS: Ensure you are logged in and Keychain Access is available.
- Windows: Ensure Credential Manager is available; run as the same user.
- Linux: Install and run a Secret Service implementation (e.g., `gnome-keyring` or `libsecret`). Make sure your desktop session unlocks the keyring and `DBUS_SESSION_BUS_ADDRESS` is set.
- Try re-running `codezure manage config` and re-entering the API key.
- If your profile was renamed, the key is stored under the profile name; re-store the key after renaming.
