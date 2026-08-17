#!/usr/bin/env bash
set -e

# Battery Remote/Local Installer Script
# Scaffolds Battery (Multi-Repository SDD Orchestrator), sets up Cooper & Troop foundations,
# installs the prebuilt Battery Go binary, and configures workspace topology.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/twoBoots/battery/main/install.sh | bash
#   or: ./install.sh [target_directory] [options]

RAW_BASE_URL="https://raw.githubusercontent.com/twoBoots/battery/main"
GITHUB_REPO="twoBoots/battery"
COOPER_RAW_BASE_URL="https://raw.githubusercontent.com/twoBoots/cooper/main"
TROOP_RAW_BASE_URL="https://raw.githubusercontent.com/twoBoots/troop/main"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}" 2>/dev/null)" && pwd || true)"
TARGET_DIR="$(pwd)"

# Parse flags
NON_INTERACTIVE=false
STRUCTURE_ARG=""

for arg in "$@"; do
    case "$arg" in
        --non-interactive|-y|--yes)
            NON_INTERACTIVE=true
            ;;
        --structure=*|-s=*)
            STRUCTURE_ARG="${arg#*=}"
            ;;
        -*)
            # Ignore other flags
            ;;
        *)
            TARGET_DIR="$arg"
            ;;
    esac
done

if [ "$CI" = "true" ]; then
    NON_INTERACTIVE=true
fi

cd "$TARGET_DIR"

if [ ! -d ".git" ] && ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    echo "Initializing Git repository in $(pwd)..."
    git init
fi

echo "🔋 Installing Battery (Multi-Repository SDD Orchestrator) into $(pwd)..."

# Helper function to fetch or copy a file from Battery repo
get_battery_file() {
    local filename="$1"
    local dest="${2:-$filename}"
    local dest_dir="$(dirname "$dest")"
    if [ "$dest_dir" != "." ] && [ ! -d "$dest_dir" ]; then
        mkdir -p "$dest_dir"
    fi

    if [ -n "$SCRIPT_DIR" ] && [ -f "$SCRIPT_DIR/$filename" ]; then
        cp "$SCRIPT_DIR/$filename" "$dest"
    elif command -v curl >/dev/null 2>&1; then
        curl -fsSL "$RAW_BASE_URL/$filename" -o "$dest"
    elif command -v wget >/dev/null 2>&1; then
        wget -qO "$dest" "$RAW_BASE_URL/$filename"
    else
        echo "Error: Neither curl nor wget found, and local $filename is missing."
        exit 1
    fi
}

# Helper function to fetch a file from Cooper repo if missing
get_cooper_file() {
    local filename="$1"
    local dest="${2:-$filename}"
    local dest_dir="$(dirname "$dest")"
    if [ "$dest_dir" != "." ] && [ ! -d "$dest_dir" ]; then
        mkdir -p "$dest_dir"
    fi

    if [ -n "$SCRIPT_DIR" ] && [ -f "$SCRIPT_DIR/../cooper/$filename" ]; then
        cp "$SCRIPT_DIR/../cooper/$filename" "$dest"
    elif command -v curl >/dev/null 2>&1; then
        curl -fsSL "$COOPER_RAW_BASE_URL/$filename" -o "$dest" 2>/dev/null || true
    elif command -v wget >/dev/null 2>&1; then
        wget -qO "$dest" "$COOPER_RAW_BASE_URL/$filename" 2>/dev/null || true
    fi
}

# 1. Setup Cooper & Troop foundation
echo "  [1/4] Setting up Cooper Hybrid SDD & Troop worktree foundation..."
if [ -n "$SCRIPT_DIR" ] && [ -f "$SCRIPT_DIR/../cooper/install.sh" ]; then
    bash "$SCRIPT_DIR/../cooper/install.sh" "$TARGET_DIR" >/dev/null 2>&1 || true
elif command -v curl >/dev/null 2>&1; then
    curl -fsSL "$COOPER_RAW_BASE_URL/install.sh" 2>/dev/null | bash -s "$TARGET_DIR" >/dev/null 2>&1 || true
fi

# Ensure base .cooper directory tree exists
mkdir -p .cooper/definition .cooper/code_styleguides .cooper/specs .cooper/active .cooper/archive

# Relocate TROOP.md into .cooper/ to keep project root clean if it exists in root
if [ -f "TROOP.md" ]; then
    mv "TROOP.md" ".cooper/TROOP.md"
fi

# Remove legacy root COOPER.md if present
rm -f "COOPER.md"

# Ensure .cooper/COOPER.md is present
if [ ! -s ".cooper/COOPER.md" ]; then
    get_cooper_file ".cooper/COOPER.md" ".cooper/COOPER.md"
fi

# 2. Install Battery specifications & configuration files
echo "  [2/4] Installing Battery specifications & guidelines..."
rm -f "strategy.md"
get_battery_file ".cooper/BATTERY.md" ".cooper/BATTERY.md"

# Setup AGENTS.md from template
if [ -f "AGENTS.md" ]; then
    if ! grep -qs "Battery Agent Rules" AGENTS.md; then
        TMP_TEMPLATE="$(mktemp)"
        if [ -n "$SCRIPT_DIR" ] && [ -f "$SCRIPT_DIR/AGENTS.template.md" ]; then
            cp "$SCRIPT_DIR/AGENTS.template.md" "$TMP_TEMPLATE"
        elif command -v curl >/dev/null 2>&1; then
            curl -fsSL "$RAW_BASE_URL/AGENTS.template.md" -o "$TMP_TEMPLATE"
        elif command -v wget >/dev/null 2>&1; then
            wget -qO "$TMP_TEMPLATE" "$RAW_BASE_URL/AGENTS.template.md"
        fi
        echo -e "\n" >> AGENTS.md
        cat "$TMP_TEMPLATE" >> AGENTS.md
        rm -f "$TMP_TEMPLATE"
        echo "  [✓] Appended Battery rules to existing AGENTS.md"
    fi
else
    get_battery_file "AGENTS.template.md" "AGENTS.md"
    echo "  [✓] Created AGENTS.md from Battery template"
fi

# Ensure .gitignore has .batteryrc.local, worktrees, and bin
if [ -f ".gitignore" ]; then
    if ! grep -qs "\.batteryrc\.local" .gitignore; then
        echo -e "\n# Battery Local Overrides\n.batteryrc.local\n.batteryrc.*.local\nbin/\ncoverage.out" >> .gitignore
    fi
else
    echo -e "# Troop Worktrees\n.worktrees/\n\n# Battery Local Overrides\n.batteryrc.local\n.batteryrc.*.local\nbin/\ncoverage.out" > .gitignore
fi
echo "  [✓] Updated .gitignore and architectural guides"

# 3. Install or Compile Battery CLI Binary
echo "  [3/4] Registering Battery CLI binary..."
CLI_INSTALLED=false

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64|amd64) ARCH="x86_64" ;;
    aarch64|arm64) ARCH="aarch64" ;;
esac

RELEASE_BINARY="battery-${OS}-${ARCH}"
INSTALL_BIN_DIR="${HOME}/.local/bin"
[ -d "/usr/local/bin" ] && [ -w "/usr/local/bin" ] && INSTALL_BIN_DIR="/usr/local/bin"
mkdir -p "$INSTALL_BIN_DIR"

RELEASE_URL="https://github.com/${GITHUB_REPO}/releases/latest/download/${RELEASE_BINARY}"
if command -v curl >/dev/null 2>&1; then
    if curl -fsSL "$RELEASE_URL" -o "${INSTALL_BIN_DIR}/battery" 2>/dev/null; then
        chmod +x "${INSTALL_BIN_DIR}/battery"
        CLI_INSTALLED=true
        echo "  [✓] Downloaded prebuilt binary from GitHub Releases to ${INSTALL_BIN_DIR}/battery"
    fi
fi

# Tier 2: Build locally if Go is available and source is present
if [ "$CLI_INSTALLED" = false ]; then
    if command -v go >/dev/null 2>&1; then
        if [ -n "$SCRIPT_DIR" ] && [ -f "$SCRIPT_DIR/main.go" ]; then
            (cd "$SCRIPT_DIR" && go build -ldflags="-s -w" -o "${INSTALL_BIN_DIR}/battery" .) >/dev/null 2>&1 || true
            if [ -f "${INSTALL_BIN_DIR}/battery" ]; then
                chmod +x "${INSTALL_BIN_DIR}/battery"
                CLI_INSTALLED=true
                echo "  [✓] Compiled and registered CLI globally with Go (${INSTALL_BIN_DIR}/battery)"
            fi
        fi
    fi
fi

# 4. Initialize .batteryrc with project structure prompt
echo "  [4/4] Configuring workspace topology..."
INIT_ARGS=()
if [ "$NON_INTERACTIVE" = true ]; then
    INIT_ARGS+=("--non-interactive")
fi
if [ -n "$STRUCTURE_ARG" ]; then
    INIT_ARGS+=("--structure" "$STRUCTURE_ARG")
fi

USE_TTY=false
if [ "$NON_INTERACTIVE" = false ] && [ -t 1 ] && [ -c /dev/tty ]; then
    USE_TTY=true
fi

BATTERY_BIN="battery"
if [ -f "${INSTALL_BIN_DIR}/battery" ]; then
    BATTERY_BIN="${INSTALL_BIN_DIR}/battery"
elif [ -f "./bin/battery" ]; then
    BATTERY_BIN="./bin/battery"
fi

if command -v "$BATTERY_BIN" >/dev/null 2>&1 || [ -x "$BATTERY_BIN" ]; then
    if [ "$USE_TTY" = true ]; then
        "$BATTERY_BIN" init "${INIT_ARGS[@]}" < /dev/tty
    else
        "$BATTERY_BIN" init "${INIT_ARGS[@]}"
    fi
else
    # Fallback default .batteryrc if binary is missing
    if [ ! -f ".batteryrc" ]; then
        echo '{"version":"1.0.0","structure":"multi-repo","barrels":[]}' > .batteryrc
        echo "  [✓] Created default .batteryrc"
    fi
fi

echo ""
echo "🔋 Battery successfully installed!"
echo "Available CLI commands:"
echo "  battery status               - View workspace structure & barrel connectivity"
echo "  battery barrel list          - List all configured barrels & Cooper tech stacks"
echo "  battery barrel add <path>    - Add a barrel to .batteryrc (or --local)"
echo "  battery barrel remove <name> - Remove a barrel from .batteryrc"
echo "  battery init                 - Reconfigure project structure & discovered barrels"
