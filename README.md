# codezure

**Launch Codex with Azure OpenAI in one command.**

A lightweight CLI that configures Codex to use Azure OpenAI automatically, handling authentication and environment setup via Azure CLI or OS keychain.

---

## Installation

```bash
curl -fsSL https://raw.githubusercontent.com/OlaHulleberg/codezure/main/install.sh | bash
```

## Quick Start

```bash
# Interactive configuration (recommended)
codezure manage config

# Launch Codex
codezure
```

## Configuration

Settings are stored in profiles at `~/.codezure/profiles/`:

| Key | Description | Example |
|-----|-------------|---------|
| `subscription` | Azure subscription ID | `12345678-1234-1234-1234-123456789abc` |
| `group` | Resource group name | `my-openai-rg` |
| `resource` | Azure OpenAI resource name | `my-openai-resource` |
| `location` | Azure region | `eastus` |
| `deployment` | Model deployment name | `gpt-5` |
| `thinking` | Thinking level (optional) | `low`, `medium`, `high` |

```bash
codezure manage config                    # Interactive configuration
codezure manage config set <key> <value>  # Set a value
codezure manage config list               # View all settings
```

### Multiple Profiles

Manage multiple named profiles for different use cases:

```bash
codezure manage profiles                  # List all profiles
codezure manage config save work-dev      # Save current config as profile
codezure manage config switch personal    # Switch to different profile
codezure --codezure-profile work-prod     # Use specific profile for one run
```

## Usage

```bash
# Launch Codex
codezure                                        # Use current profile
codezure --codezure-profile work-dev             # Use specific profile

# Pass Codex CLI flags
codezure --resume                               # Resume last session
codezure --continue                             # Continue last session
codezure --debug                                # Debug mode
codezure --print "analyze this code"            # Non-interactive mode

# Combined (codezure config + Codex CLI passthrough)
codezure --codezure-profile work --resume --debug

# Configuration
codezure manage config                          # Interactive wizard
codezure manage config list                     # View settings
codezure manage config set deployment <value>   # Update setting

# Codex configuration (overrides, non-destructive)
# codezure passes -c model_provider="codezure" and provider details:
#   -c model_providers.codezure.name="Codezure"
#   -c model_providers.codezure.base_url="<endpoint>/openai/v1"
#   -c model_providers.codezure.env_key="CODEZURE_API_KEY"
#   -c model_providers.codezure.wire_api="responses"
# It also sets -c model="<deployment>" and optionally -c model_reasoning_effort="<low|medium|high>".

# Profiles
codezure manage profiles                        # List profiles
codezure manage config save my-profile          # Save as new profile
codezure manage config switch my-profile        # Switch profile
codezure manage config copy prod staging        # Copy profile
codezure manage config rename old new           # Rename profile
codezure manage config delete old-profile       # Delete profile

# Models
codezure manage models list                     # List available deployments
Note: Requires Azure CLI authentication.

# Updates
codezure manage update                          # Update to latest version
codezure manage version                         # Show version
```

### Override Flags

Override profile settings for a single run without changing your saved profile:

```bash
codezure --codezure-profile production
```

## What It Does

1. Loads your Azure OpenAI configuration from the current profile
2. Uses your chosen auth mode:
   - Azure CLI: fetches keys and endpoint via `az` at runtime
   - Keychain: retrieves API key from OS keychain; uses saved endpoint/deployment
3. Launches `codex` with the correct configuration overrides and `CODEZURE_API_KEY` in the child environment
4. Passes through any Codex CLI flags you provide

## Features

### 📋 Multiple Profiles
Save and switch between different Azure OpenAI configurations (work, personal, different projects).

### 🔍 Model Discovery
List all available model deployments in your Azure OpenAI resource with interactive selection.

### ⚡ Quick Overrides
Override profile settings for a single run using command-line flags.

### 🔄 Codex CLI Passthrough
Pass any Codex CLI flags and commands directly through codezure (e.g., `--resume`, `--debug`, `--print`).

### 🔐 Authentication Options
- Azure CLI (recommended): zero manual key management; keys fetched on-demand.
- OS Keychain (manual): store API key securely in the OS keychain; enter endpoint/deployment.

### 🔒 Privacy-First
- In Azure CLI mode, keys are fetched on-demand and never persisted.
- In Keychain mode, keys are stored securely in the OS keychain; never written to disk.

## Documentation

- **[Configuration Guide](CONFIGURATION.md)** - Detailed config options and profiles
- **[Pricing & Costs](PRICING.md)** - Azure OpenAI pricing information and cost estimates
- **[Troubleshooting](TROUBLESHOOTING.md)** - Common issues and solutions
- **[Contributing](CONTRIBUTING.md)** - Development guide and contribution guidelines

## Requirements

- [Codex](https://openai.com/codex/) (Codex CLI) installed
- One of:
  - [Azure CLI](https://aka.ms/azcli) configured (`az login`), or
  - OS keychain available (macOS Keychain, Windows Credential Manager, Linux Secret Service)
- Azure OpenAI access with model deployments

## License

MIT
