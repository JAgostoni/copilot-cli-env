# Copilot Env Manager Product Brief

## Overview

Copilot Env Manager is a cross-platform Go CLI/TUI that helps users configure GitHub Copilot CLI for BYOK providers and OpenAI-compatible endpoints by walking them through provider selection, live model discovery, and shell-aware environment output generation.[1][2]

GitHub Copilot CLI documents BYOK configuration through environment variables including `COPILOT_PROVIDER_BASE_URL`, `COPILOT_MODEL`, optional `COPILOT_PROVIDER_TYPE`, and optional `COPILOT_PROVIDER_API_KEY`, which makes an environment-focused configuration tool a natural fit.[1]

## Problem

Users who want GitHub Copilot CLI to run against providers such as Anthropic, Azure OpenAI, or OpenAI-compatible services like OpenRouter must set multiple environment variables before launching Copilot CLI, and the exact persistence mechanism varies by shell, OS, and terminal host.[1]

That setup is easy to get wrong across zsh, bash, fish, PowerShell, CMD, Git Bash, and mixed Windows terminal contexts, especially when users want a persistent configuration rather than a one-off shell session.[3][1]

## Goals

- Eliminate the need for users to memorize Copilot CLI BYOK environment variables.[1]
- Make provider and model selection interactive and discoverable.[1]
- Support both safe ephemeral workflows and persistent setup workflows across macOS, Linux, and Windows.[3][1]
- Handle shell and OS differences without hiding them from advanced users.[3]
- Make API-driven model discovery a first-class workflow, while always allowing manual model entry as an override.[1]

## Non-goals

- Replacing GitHub Copilot CLI itself.[1]
- Managing long-running Copilot sessions or chat state.[1]
- Deep secret-vault integration in v1, since Copilot ultimately consumes environment variables and v1 should stay aligned with that model.[1]
- Windows registry mutation in v1.[3]

## Product summary

The product should launch into an onboarding/configuration flow that detects the user’s current shell and terminal context, presents that detection as the default, and always allows override before any output is generated.[3]

A second step in onboarding should let the user choose the desired output mode, such as profile persistence, `.env` generation, console copy/paste, or shell-specific command generation.[3]

## Primary users

- Developers already using GitHub Copilot CLI who want to switch to BYOK providers or local models without hand-editing shell profile files.[1]
- Power users working across multiple terminals, shells, and operating systems who need repeatable and correct environment setup.[3]
- Teams experimenting with provider routing, cost control, or local/offline Copilot CLI operation.[1]

## Key product decisions

### Implementation language

Go is the recommended implementation language for v1 because it offers strong cross-platform support, straightforward filesystem and subprocess handling, single-binary distribution, and a mature TUI ecosystem anchored by Bubble Tea.[2]

Rust remains a viable alternative, but Go is better aligned with shipping a polished MVP quickly because this project is terminal-UX-heavy rather than CPU-bound.[2]

### Provider model

The UI should present provider choices in user-friendly language while mapping back to GitHub Copilot CLI’s documented provider model of `openai`, `azure`, and `anthropic`.[1]

OpenRouter should be presented as an OpenAI-compatible preset rather than as a distinct first-class Copilot provider type, because Copilot’s documented provider model treats OpenAI-compatible endpoints under the `openai` path.[1]

### Secret handling

Version 1 should remain environment-oriented and avoid inventing a secret-storage abstraction that Copilot CLI cannot directly consume.[1]

Future profile switching can introduce a local encrypted file or database under the user profile folder, but activation should still materialize usable environment variables for Copilot CLI.[1]

### Windows scope

Windows support in v1 should avoid registry writes and instead focus on PowerShell profile handling, Git Bash profile handling, session-only commands, `.env` generation, and copy/paste output.[3]

Because Windows Terminal is a host rather than the shell itself, the product should distinguish between terminal host, current shell, and persistence target in its UI and logic.[3]

## Functional requirements

### Environment and shell detection

The tool must detect the operating system, current shell, and terminal host when possible, including macOS, Linux, Windows, bash, zsh, fish, PowerShell, CMD, Git Bash, and WSL-adjacent shell contexts where relevant.[3]

Detection must be advisory rather than absolute, with the detected environment presented as a default that the user can override during onboarding.[3]

### Provider selection

The tool must present a list of provider presets supported by GitHub Copilot CLI’s BYOK model, including OpenAI-compatible, OpenRouter, OpenAI, Azure OpenAI, Anthropic, Ollama/local-compatible, and custom endpoint options.[1]

Each provider preset should map to the correct resulting Copilot environment-variable shape, including required base URL, provider type when applicable, model identifier, and API key when needed.[1]

### Dynamic model discovery

Model discovery is a primary feature and should occur through provider APIs after the user supplies the minimum viable credentials and endpoint information required for that provider.[1]

The model picker should be searchable, selectable, and type-ahead-driven, while always allowing arbitrary manual entry if discovery fails, the provider does not expose a suitable endpoint, or the user wants to override the discovered catalog.[1]

### Output modes

The tool must support at least four output modes in v1: shell profile persistence, `.env` file generation, direct console dump for copy/paste, and shell-specific command generation.[3][1]

These output modes are necessary because Copilot CLI’s documented BYOK flow depends on environment variables being present before launch, but users vary widely in whether they want persistent shell mutation or ephemeral session activation.[1]

### Env generation

The tool must generate valid GitHub Copilot CLI environment variables using the documented BYOK interface, including `COPILOT_PROVIDER_BASE_URL`, `COPILOT_MODEL`, optional `COPILOT_PROVIDER_TYPE`, optional `COPILOT_PROVIDER_API_KEY`, and `COPILOT_OFFLINE=true` when relevant.[1]

GitHub documents `COPILOT_PROVIDER_BASE_URL` and `COPILOT_MODEL` as required, with provider type optional and API key optional for providers that do not require authentication, such as some local endpoints.[1]

### Safe profile persistence

If the user chooses profile persistence, the tool should preview the target profile file, show the exact managed block that will be written, and require explicit confirmation before making changes.[3]

Profile writes should be idempotent where possible by using a clearly marked managed block that can be updated or replaced instead of appended repeatedly.[3]

## User experience requirements

### Onboarding flow

The canonical UX should be a guided onboarding/configuration flow that begins with detection, then override confirmation, then provider selection, then model discovery, and finally output-mode choice and preview.[3]

The workflow should use progressive disclosure so users are only asked for the minimum information needed to unlock the next step, rather than presenting a large up-front configuration form.[3][1]

### Terminal-native interaction

The application should be keyboard-first, fast, and usable in ANSI terminals with a clean TUI that degrades gracefully to less-rich text output when color or styling is unavailable.[2]

Searchable pickers, type-ahead inputs, and clear confirmation screens should be treated as core interaction patterns rather than optional polish.[2]

### Transparency

The product should always show both the friendly abstraction, such as “OpenRouter preset,” and the raw resulting Copilot environment variables so advanced users can verify and trust the output.[1]

## Suggested command surface

A practical command surface for the CLI could include interactive entry points such as `init`, `configure`, `detect`, `models`, `render`, and `apply`, with flags to support partially non-interactive onboarding when configuration is supplied as parameters.[2]

This command design would support both a guided first-run experience and more scriptable use cases for advanced users who want to seed the flow with provider, shell, or output-mode defaults.[2]

## Technical architecture

The system should be organized into clear modules for shell detection, provider adapters, model discovery, output rendering, environment mapping, UI flow, and safety features such as masking and preview.[2][1]

Provider adapters should define required fields, base URL logic, authentication rules, model-list fetching, env mapping, and validation behavior so presets like OpenRouter and Ollama can behave as opinionated OpenAI-compatible flavors while still mapping back to Copilot’s documented model.[1]

## Security model

In v1, secrets should be entered interactively, masked in the UI, redacted in previews by default, and never logged.[1]

Whenever the user chooses an output mode that writes secrets to plaintext locations such as shell profile files or `.env` files, the product should provide an explicit warning before completion.[3]

## Future roadmap

A future version can add named saved profiles such as `openrouter-work`, `anthropic-personal`, or `ollama-local`, stored in a local encrypted file or database under the user profile directory and unlocked using a user-entered PIN as part of the key derivation flow.[1]

That future activation flow should allow rapid provider switching while still rendering the final active configuration as environment variables that Copilot CLI can consume.[1]

## Risks and edge cases

Model discovery will not be uniform across providers, so the design should assume provider-specific adapters and graceful fallback to manual entry rather than one universal catalog strategy.[1]

Shell persistence is inherently messy because startup file selection differs across shells, operating systems, and user dotfile setups, which is why detection should remain suggestive and user-overridable rather than fully automatic and silent.[3]

Writing keys to persistent files is convenient but risky, so the product must make plaintext-secret storage visible and intentional rather than implicit.[3]

## Acceptance criteria

A v1 build is successful if a user on macOS, Linux, or Windows can launch the tool, accept or override a detected shell target, choose a provider preset, fetch models live from the provider API, optionally override the model manually, and emit valid Copilot CLI environment variables in one of the supported output modes.[3][1]

The product should also support PowerShell, Git Bash, and CMD-oriented outputs on Windows without touching the registry, while maintaining a preview-first workflow for any persistent profile mutation.[3]