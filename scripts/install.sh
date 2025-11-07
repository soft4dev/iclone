#!/usr/bin/env bash

set -euo pipefail

REPO="soft4dev/clonei"
BIN_NAME="clonei"
UPDATE_MODE=false

# Check for update argument
if [ "${1:-}" = "update" ]; then
    UPDATE_MODE=true
fi

# Colors (only if stdout is a terminal)
if [ -t 1 ]; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[1;33m'
    NC='\033[0m'
else
    RED=''
    GREEN=''
    YELLOW=''
    NC=''
fi

# --- Helper functions ---
get_os() {
    os=$(uname -s)
    case "$os" in
        Darwin) echo "Darwin" ;;
        Linux) echo "Linux" ;;
        *)
            printf "${RED}Unsupported OS: %s${NC}\n" "$os" >&2
            printf "This script supports macOS and Linux only.\n" >&2
            printf "Windows: irm https://raw.githubusercontent.com/%s/main/install.ps1 | iex\n" "$REPO" >&2
            exit 1
            ;;
    esac
}

get_arch() {
    arch=$(uname -m)
    case "$arch" in
        x86_64|amd64) echo "x86_64" ;;
        aarch64|arm64) echo "arm64" ;;
        i386|i686) echo "i386" ;;
        armv7l) echo "armv7" ;;
        *)
            printf "${RED}Unsupported architecture: %s${NC}\n" "$arch" >&2
            exit 1
            ;;
    esac
}

get_latest_version() {
    curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" |
        grep -E '"tag_name"|"name"' |
        head -n1 |
        sed -E 's/.*"([^"]+)".*/\1/'
}

download_file() {
    url="$1"
    dest="$2"
    if command -v curl >/dev/null 2>&1; then
        curl -fL --retry 3 --progress-bar "$url" -o "$dest"
    elif command -v wget >/dev/null 2>&1; then
        wget --tries=3 -q "$url" -O "$dest"
    else
        printf "${RED}Error: curl or wget is required${NC}\n" >&2
        exit 1
    fi
}

# --- Main installer ---
install_binary() {
    os=$(get_os)
    arch=$(get_arch)
    version=$(get_latest_version)

    if [ -z "$version" ]; then
        printf "${RED}Error: Could not determine latest release version${NC}\n" >&2
        exit 1
    fi

    if [ "$UPDATE_MODE" = true ]; then
        printf "${GREEN}Updating %s to %s...${NC}\n" "$BIN_NAME" "$version"
    else
        printf "${GREEN}Installing %s %s...${NC}\n" "$BIN_NAME" "$version"
    fi

    archive_name="${BIN_NAME}_${os}_${arch}.tar.gz"
    download_url="https://github.com/${REPO}/releases/download/${version}/${archive_name}"

    tmp_dir=$(mktemp -d)
    trap "rm -rf '$tmp_dir'" EXIT

    printf "${YELLOW}Downloading from %s...${NC}\n" "$download_url"
    download_file "$download_url" "$tmp_dir/$archive_name"

    printf "${YELLOW}Extracting archive...${NC}\n"
    tar -xzf "$tmp_dir/$archive_name" -C "$tmp_dir"

    # Choose install directory
    if [ -n "${BIN_DIR:-}" ]; then
        bin_dir="$BIN_DIR"
    elif [ -w /usr/local/bin ]; then
        bin_dir="/usr/local/bin"
    else
        bin_dir="$HOME/.local/bin"
    fi

    mkdir -p "$bin_dir"

    # Move binary safely
    target="$bin_dir/$BIN_NAME"
    if [ -f "$target" ]; then
        if [ "$UPDATE_MODE" = true ]; then
            printf "${YELLOW}Updating existing installation at %s...${NC}\n" "$target"
        else
            printf "${YELLOW}Warning: %s already exists. Overwrite? [y/N]: ${NC}" "$target"
            read -r reply </dev/tty
            case "$reply" in
                [yY]*) ;;
                *) printf "${RED}Installation aborted.${NC}\n"; exit 1 ;;
            esac
        fi
    fi

    mv -f "$tmp_dir/$BIN_NAME" "$target"
    chmod +x "$target"

    if [ "$UPDATE_MODE" = true ]; then
        printf "${GREEN}✓ %s updated successfully at %s${NC}\n" "$BIN_NAME" "$target"
    else
        printf "${GREEN}✓ %s installed to %s${NC}\n" "$BIN_NAME" "$target"
    fi

    # PATH setup
    case ":$PATH:" in
        *":$bin_dir:"*) 
            printf "${GREEN}PATH already includes %s${NC}\n" "$bin_dir"
            ;;
        *)
            printf "${YELLOW}Adding %s to PATH...${NC}\n" "$bin_dir"
            shell_profile=""
            # Check user's actual shell from SHELL env var, not the script's runtime shell
            case "${SHELL:-}" in
                */bash) 
                    shell_profile="$HOME/.bashrc"
                    [ "$(uname -s)" = "Darwin" ] && shell_profile="$HOME/.bash_profile"
                    ;;
                */zsh)  
                    shell_profile="$HOME/.zshrc"
                    ;;
                */fish) 
                    shell_profile="$HOME/.config/fish/config.fish"
                    ;;
                *)      
                    shell_profile="$HOME/.profile"
                    ;;
            esac

            if [ -n "$shell_profile" ]; then
                mkdir -p "$(dirname "$shell_profile")"
                if ! grep -qs "$bin_dir" "$shell_profile"; then
                    {
                        echo ""
                        echo "# Added by ${BIN_NAME} installer"
                        echo "export PATH=\"$bin_dir:\$PATH\""
                    } >> "$shell_profile"
                    printf "${GREEN}✓ Added to %s${NC}\n" "$shell_profile"
                else
                    printf "${YELLOW}PATH already configured in %s${NC}\n" "$shell_profile"
                fi
            fi

            export PATH="$bin_dir:$PATH"
            ;;
    esac

    if [ "$UPDATE_MODE" = true ]; then
        printf "\n${GREEN}✅ Update complete!${NC}\n"
        printf "Run: ${YELLOW}%s --help${NC}\n" "$BIN_NAME"
    else
        printf "\n${GREEN}✅ Installation complete!${NC}\n"
        printf "Run: ${YELLOW}%s --help${NC}\n" "$BIN_NAME"
        
        if [ -n "${shell_profile:-}" ] && [ -f "$shell_profile" ]; then
            printf "\n${YELLOW}To use %s immediately, run:${NC}\n" "$BIN_NAME"
            printf "  ${GREEN}source %s${NC}\n" "$shell_profile"
            printf "Or restart your terminal.\n"
        else
            printf "You may need to restart your terminal.\n"
        fi
    fi
}

install_binary
