# devctl-tempo

A Go CLI for interacting with the Tempo API, using Cobra and Viper. Stores secrets securely on MacOS using the system keychain.

## Prerequisites
- Go 1.24.3+
- (For MacOS secrets) Keychain access

## Setup
Clone the repository and install dependencies:
```sh
git clone <repo-url>
cd devctl-tempo
go mod tidy
```

## Testing
Run all tests:
```sh
make test
```

## Building
Build for your current OS and architecture:
```sh
make build
```
The binary will be output to `bin/devctl-tempo-<os>-<arch>-<hash>`.

Build for all supported OS/arch:
```sh
make build-all
```
Binaries for Linux, MacOS (amd64/arm64) will be in `bin/`.

Build a standalone binary named `tempo` for your current system:
```sh
make build-standalone
```
The binary will be output to `bin/tempo`.

Build standalone binaries for all supported OS/arch:
```sh
make build-all-standalone
```
Standalone binaries (tempo-linux-amd64, tempo-linux-arm64, tempo-darwin-amd64, tempo-darwin-arm64) will be in `bin/`.

Cross-compile for a specific system:
```sh
GOOS=linux GOARCH=amd64 make build
```

## Quick Install (Online Script)

### Versioned Build (with git hash)
You can install the latest `devctl-tempo` binary automatically with a one-liner (Linux/macOS):

```sh
curl -sSL https://raw.githubusercontent.com/danlafeir/devctl-tempo/main/scripts/install.sh | sh
```

This script will detect your OS and architecture, download the latest versioned binary from the main branch, and install it to `~/.local/bin` (you may be prompted for your password).

### Standalone Build
You can install the standalone `tempo` binary with:

```sh
curl -sSL https://raw.githubusercontent.com/danlafeir/devctl-tempo/main/scripts/install-standalone.sh | sh
```

This script will detect your OS and architecture, download the standalone binary (`tempo-<os>-<arch>`) and install it as `tempo` to `~/.local/bin` (you may be prompted for your password).

**Security tip:** Always review install scripts before piping to `sh`.

## Running
After building, run the CLI:
```sh
./bin/devctl-tempo-<os>-<arch>-<hash> <command>
```

For example:
```sh
./bin/devctl-tempo-darwin-arm64-d58b1d5 <command>
```

## Configuration
To configure your Tempo API token, account ID, and default issue ID:
```sh
./bin/devctl-tempo configure --token <YOUR_TOKEN> --account-id <ACCOUNT_ID>
```
You can also omit flags to be prompted interactively.

- The API token is stored securely in the MacOS keychain.
- The account ID and default issue ID are stored in `$HOME/.devctl/config.yaml` under a `tempo` key.

**Note:** The account ID is associated with your profile in JIRA. You can find it by doing a "people search" in JIRA.

## Available Commands

### `timesheet`
Submit a Tempo timesheet for the current week (or a past week). This is the main command for logging time.

```sh
./bin/devctl-tempo timesheet
```

The command will:
1. Prompt you to confirm the week (defaults to current week, or you can specify weeks back)
2. Ask for time spent in three categories:
   - Development/design/testing (capitalizable time)
   - PTO (vacation or sick time)
   - Meeting time
3. Submit all time entries to Tempo via the API

### `configure`
Set up your Tempo API token, account ID, and default issue ID.

```sh
./bin/devctl-tempo configure
```

Options:
- `--token` - Tempo API token
- `--account-id` - Your Tempo account ID (from JIRA)

The issue ID will be automatically fetched from your most recent Tempo worklog entry. If the API call fails, you'll be prompted to enter it manually.

### Hidden Commands
- `get-week` — Fetch your current week's timecard from the Tempo API (for debugging)
- `completion` — Generate shell completion scripts (bash, zsh, fish, powershell)

## Notes
- Only MacOS is currently supported for secure secrets storage.
- For other OS support, contributions are welcome!

---

For more information, see the code and comments or open an issue.
