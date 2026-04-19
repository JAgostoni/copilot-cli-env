# Session 8: TUI Part 2 - Model Picker & Output Selection

## Goal
Complete the interactive flow by wiring up live network discovery and finalizing the output choices.

## Tooling & Library Versions
- Same as Session 7 (Bubble Tea ecosystem).

## Setup Instructions
```text
internal/
  ui/
    step_model.go    # Step 4: Model Picker
    step_output.go   # Step 5: Output Selection
    step_summary.go  # Step 6: Review & Apply
```

## Tasks

### 1. Build Step 4: Live Model Picker
This is the most complex UI component.
- When entering this step, trigger a `tea.Cmd` to fetch models asynchronously using the logic from Session 4.
- Show a loading spinner (`bubbles/spinner`) while fetching.
- Once fetched, display the models in a searchable `bubbles/list`.
- **Fallback:** Add a hotkey (e.g., `ctrl+m`) to switch from the list view to a raw `bubbles/textinput` to allow manual entry if the API fails or the user wants a model not listed.

### 2. Build Step 5: Output Mode Selection
Use a simple list to let the user choose: `Profile Update`, `.env file`, `Console Output`.

### 3. Build Step 6: Final Review
Display a Lipgloss-styled summary of what will happen:
```text
Target Shell: zsh
Provider: OpenRouter
Model: anthropic/claude-3-opus
API Key: sk-or-v1-****

Action: Will update ~/.zshrc
```
Require an explicit `y/N` confirmation before executing.

## Manual Testing
Run the full flow manually to verify:
1. Select OpenRouter.
2. Enter a valid key.
3. Observe the loading spinner and subsequent list of models.
4. Select a model.
5. Reach the summary screen.
6. Verify the API key is masked in the summary screen.
