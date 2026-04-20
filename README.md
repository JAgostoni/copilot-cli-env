# Copilot Env Manager

Copilot Env Manager (`copilot-cli-env`) is an interactive terminal and command-line utility designed to help you quickly, safely, and securely manage custom AI provider environments for the **GitHub Copilot CLI**.

By leveraging the "Bring Your Own Key" (BYOK) capabilities of Copilot CLI, this tool allows you to easily switch between providers like OpenAI, Anthropic, OpenRouter, Azure, and local Ollama instances, dynamically mapping models, token limits, and wire APIs.

## Features

- **Interactive TUI:** A keyboard-driven, beautiful terminal UI to guide you through provider and model selection.
- **Live Model Discovery:** Automatically connects to provider APIs to fetch the latest available models and injects their context windows and token limits into your configuration.
- **Smart Formatting:** Generates the exact `export` or `$env:` syntax needed for your specific shell (Bash, Zsh, PowerShell, CMD) and intelligently writes it directly to your profile.
- **Non-Destructive `.env` Management:** Safely appends and removes Copilot-specific variables from your project's `.env` files without disrupting existing keys.
- **Automated Fallbacks:** Automatically configures internal API parameters (e.g. `COPILOT_PROVIDER_WIRE_API=responses` for GPT-5 models) to ensure seamless compatibility.

## Installation & Build

Ensure you have Go 1.26.x installed.

```bash
git clone https://github.com/jagostoni/copilot-cli-env.git
cd copilot-cli-env

# Build using the standard Go toolchain
go build -o copilot-cli-env main.go

# Alternatively, build cross-platform binaries using the provided Makefile
make build

# The Makefile binaries will be available in the bin/ directory:
# ./bin/copilot-cli-env-linux
# ./bin/copilot-cli-env-mac
# ./bin/copilot-cli-env.exe
```

## Usage

### Interactive Mode (Recommended)

Simply run the tool without any arguments to launch the interactive onboarding flow:

```bash
copilot-cli-env
# or
copilot-cli-env init
```

Use your arrow keys, `Tab`, and `Enter` to navigate. During the model selection phase, simply **start typing** to filter the live model list using a strict substring search, or press `Ctrl+T` to manually type in a custom model name.

### Non-Interactive Mode (Automation)

You can bypass the TUI entirely and configure your environment in a single command using the `configure` subcommand.

**Example: Write to your Shell Profile**
```bash
copilot-cli-env configure \
  --provider anthropic \
  --model claude-3-opus-20240229 \
  --api-key sk-ant-api03... \
  --output profile
```

**Example: Generate a `.env` file**
```bash
copilot-cli-env configure \
  --provider openrouter \
  --model google/gemma-7b-it \
  --api-key sk-or-v1-... \
  --output env
```

**Example: Print to Console**
```bash
copilot-cli-env configure \
  --provider openai \
  --model gpt-4o \
  --api-key sk-... \
  --output console
```

### Removing Configuration

To safely remove the injected Copilot configurations from your shell profile or `.env` file, use the `reset` provider:

```bash
# Interactive
copilot-cli-env

# Non-Interactive
copilot-cli-env configure --provider reset --output profile
```

## Supported Providers

- **OpenAI** (Requires API Key)
- **Anthropic** (Requires API Key)
- **Azure OpenAI** (Requires API Key & Custom Resource URL)
- **OpenRouter** (Requires API Key)
- **Ollama** (Local, Offline Mode)

## Security Note

When outputting configurations to a `.env` file or injecting them directly into your shell profile, `copilot-cli-env` writes your API key in plaintext. The tool will warn you before performing this action. Always ensure your `.env` files are excluded from source control (e.g., via `.gitignore`).
