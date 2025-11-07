<img width="282" height="69" alt="Group 10" src="https://github.com/user-attachments/assets/2758038f-d348-4267-82a8-4d25f200804c" />

A CLI tool to clone Git repositories and automatically install dependencies in a single shot.

## Installation

For macOS / Linux

```sh
curl -fsSL https://raw.githubusercontent.com/soft4dev/clonei/main/scripts/install.sh | bash
```

For Windows [PowerShell](https://learn.microsoft.com/en-us/powershell/)

```powershell
irm 'https://raw.githubusercontent.com/soft4dev/clonei/main/scripts/install.ps1' | iex
```

For Windows [cmd](https://learn.microsoft.com/en-us/windows-server/administration/windows-commands/windows-commands)

```cmd
powershell -c "irm 'https://raw.githubusercontent.com/soft4dev/clonei/main/scripts/install.ps1' | iex"
```

## Updating

```bash
clonei update
```

## Usage

### Basic usage (auto-detect project type)

```bash
clonei <repository-url>
```

Example:

```bash
clonei https://github.com/username/my-project.git
```

### Specify project type manually

```bash
clonei -p <project-type> <repository-url>
```

Example:

```bash
clonei -p npm https://github.com/username/my-project.git
```

## Remove

Just remove the binary at `~/.local/bin/clonei` and done.

## Supported Project Types

- **Node.js**: npm, pnpm
- rust
- golang
- php: composer
- java: maven
- more to be added...

## How it works

1. Clones the specified Git repository
2. Detects the project type (or uses the specified type)
3. Automatically installs dependencies based on the project type
4. Finally cd into the project

## License

See [LICENSE](LICENSE) for details.
