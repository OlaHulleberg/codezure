# codzure

**Launch Codex with Azure OpenAI in one command.**

A lightweight CLI that configures Codex to use Azure OpenAI automatically, handling authentication and environment setup.

---

## Installation

```bash
curl -fsSL https://raw.githubusercontent.com/OlaHulleberg/codzure/main/install.sh | bash
```

## Quick Start

```bash
# Interactive configuration (recommended)
codzure manage config

# Or configure manually
codzure manage config set subscription <SUBSCRIPTION_ID>
codzure manage config set group <RESOURCE_GROUP>
codzure manage config set resource <OPENAI_RESOURCE>
codzure manage config set deployment <MODEL_DEPLOYMENT>

# Launch Codex
codzure
```

## Configuration

Settings are stored in profiles at `~/.codzure/profiles/`:

| Key | Description | Example |
|-----|-------------|---------|
| `subscription` | Azure subscription ID | `12345678-1234-1234-1234-123456789abc` |
| `group` | Resource group name | `my-openai-rg` |
| `resource` | Azure OpenAI resource name | `my-openai-resource` |
| `location` | Azure region | `eastus` |
| `deployment` | Model deployment name | `gpt-5` |
| `thinking` | Thinking level (optional) | `low`, `medium`, `high` |

```bash
codzure manage config                    # Interactive configuration
codzure manage config set <key> <value>  # Set a value
codzure manage config list               # View all settings
```

### Multiple Profiles

Manage multiple named profiles for different use cases:

```bash
codzure manage profiles                  # List all profiles
codzure manage config save work-dev      # Save current config as profile
codzure manage config switch personal    # Switch to different profile
codzure --codzure-profile work-prod      # Use specific profile for one run
```

## Usage

```bash
# Launch Codex
codzure                                        # Use current profile
codzure --codzure-profile work-dev             # Use specific profile

# Pass Codex CLI flags
codzure --resume                               # Resume last session
codzure --continue                             # Continue last session
codzure --debug                                # Debug mode
codzure --print "analyze this code"            # Non-interactive mode

# Combined (codzure config + Codex CLI passthrough)
codzure --codzure-profile work --resume --debug

# Configuration
codzure manage config                          # Interactive wizard
codzure manage config list                     # View settings
codzure manage config set deployment <value>   # Update setting

# Profiles
codzure manage profiles                        # List profiles
codzure manage config save my-profile          # Save as new profile
codzure manage config switch my-profile        # Switch profile
codzure manage config copy prod staging        # Copy profile
codzure manage config rename old new           # Rename profile
codzure manage config delete old-profile       # Delete profile

# Models
codzure manage models list                     # List available deployments

# Updates
codzure manage update                          # Update to latest version
codzure manage version                         # Show version
```

### Override Flags

Override profile settings for a single run without changing your saved profile:

```bash
codzure --codzure-profile production
```

## What It Does

1. Loads your Azure OpenAI configuration from the current profile
2. Fetches API keys and endpoint from Azure using `az` CLI
3. Launches `codex` with the correct environment variables set (`AZURE_OPENAI_*`, `OPENAI_*`)
4. Passes through any Codex CLI flags you provide

## Features

### 📋 Multiple Profiles
Save and switch between different Azure OpenAI configurations (work, personal, different projects).

### 🔍 Model Discovery
List all available model deployments in your Azure OpenAI resource with interactive selection.

### ⚡ Quick Overrides
Override profile settings for a single run using command-line flags.

### 🔄 Codex CLI Passthrough
Pass any Codex CLI flags and commands directly through codzure (e.g., `--resume`, `--debug`, `--print`).

### 🔐 Zero-Config Authentication
Automatically fetches Azure OpenAI keys and endpoints via Azure CLI—no manual key management.

### 🔒 Privacy-First
All configuration stored locally in `~/.codzure/`. Keys are fetched on-demand and never persisted.

## Documentation

- **[Configuration Guide](CONFIGURATION.md)** - Detailed config options and profiles
- **[Pricing & Costs](PRICING.md)** - Azure OpenAI pricing information and cost estimates
- **[Troubleshooting](TROUBLESHOOTING.md)** - Common issues and solutions
- **[Contributing](CONTRIBUTING.md)** - Development guide and contribution guidelines

## Requirements

- [Codex](https://openai.com/codex/) (Codex CLI) installed
- [Azure CLI](https://aka.ms/azcli) configured (`az login`)
- Azure OpenAI access with model deployments

## License

MIT
