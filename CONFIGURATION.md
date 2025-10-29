# Configuration Guide (Codezure)

`codezure` stores configuration in profiles at `~/.codezure/profiles/`.

## Profile Management

- Profiles let you switch between different Azure setups.
- Current profile tracked at `~/.codezure/current-profile.txt`.

### Interactive Configuration (First Run)

Run the interactive wizard:

```
codezure manage config
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
codezure manage config list
codezure manage config set <key> <value>
```

Keys: `auth` (`azure-cli` or `api-key`), `subscription`, `group`, `resource`, `location`, `endpoint`, `deployment`, `thinking`

## Migration from Old Config

Profiles are stored under `~/.codezure`. The installer offers migration from any legacy setup.

## Env Variables

Exported when launching:
- `AZURE_OPENAI_ENDPOINT`
- `AZURE_OPENAI_DEPLOYMENT`
- Child process env: `CODEZURE_API_KEY`

Notes:
- In `api-key` mode, the API key is retrieved from the OS keychain per-profile.
- In `azure-cli` mode, keys are fetched via `az` at runtime and are never persisted.

## Codex Configuration (Overrides)

Codex uses `~/.codex/config.toml` by default and supports runtime overrides via `--config/-c key=value`.

To avoid changing system configs, codezure configures Codex at launch using overrides and environment variables:

- Defines the Codezure provider and routes to the Azure Responses API:
  - `-c model_provider="codezure"`
  - `-c model_providers.codezure.name="Codezure"`
  - `-c model_providers.codezure.base_url="<endpoint>/openai/v1"`
  - `-c model_providers.codezure.env_key="CODEZURE_API_KEY"`
  - `-c model_providers.codezure.wire_api="responses"`
- Sets the model and optional reasoning effort:
  - `-c model="<deployment>"`
  - `-c model_reasoning_effort="<low|medium|high>"`

If you prefer to manage `~/.codex/config.toml` yourself, codezure respects any overrides you pass and simply provides the child environment.
