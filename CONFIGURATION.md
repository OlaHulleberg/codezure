# Configuration Guide (Codzure)

`codzure` stores configuration in profiles at `~/.codzure/profiles/`.

## Profile Management

- Profiles let you switch between different Azure setups.
- Current profile tracked at `~/.codzure/current-profile.txt`.

### Interactive Configuration (First Run)

Run the interactive wizard:

```
codzure manage config
```

Choose an authentication method:
- Azure CLI (recommended): discovers subscriptions/resources/deployments via `az` and fetches keys at runtime.
- Keychain API Key: enter endpoint and deployment; API key is stored securely in the OS keychain.

Azure CLI prompts for:
- Subscription ID or name
- Resource group name
- Location (e.g., `eastus`)
- Azure OpenAI resource name
- Deployment name (e.g., `gpt-5`)
- Thinking level (optional: `low`, `medium`, `high`)

Keychain prompts for:
- Endpoint URL (e.g., `https://<resource>.openai.azure.com`)
- Deployment name (e.g., `gpt-5`)
- API key (masked, stored in OS keychain)

### List and Set

```
codzure manage config list
codzure manage config set <key> <value>
```

Keys: `auth` (`azure-cli` or `api-key`), `subscription`, `group`, `resource`, `location`, `endpoint`, `deployment`, `thinking`

## Migration from Old Config

If an old `~/.codzure/current.env` exists, it is migrated to `~/.codzure/profiles/default.env` on first run. The old file is backed up as `current.env.bak` and `default` becomes current.

## Env Variables

Exported when launching:
- `AZURE_OPENAI_API_KEY`
- `AZURE_OPENAI_ENDPOINT`
- `AZURE_OPENAI_DEPLOYMENT`
- `OPENAI_API_KEY`
- `OPENAI_BASE_URL` (uses v1 API endpoint)

Notes:
- In `api-key` mode, the API key is retrieved from the OS keychain per-profile.
- In `azure-cli` mode, keys are fetched via `az` at runtime and are never persisted.
