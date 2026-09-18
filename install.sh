#!/usr/bin/env bash
# ==============================================================================
# The Legend of Zelda: Breath of the Wild Multiplayer - Automated Installer
# Repository: https://github.com/CarlosEvCode/botw-multiplayer-linux
# ==============================================================================
set -e

REPO_OWNER="CarlosEvCode"
REPO_NAME="botw-multiplayer-linux"
RELEASE_TAG="v1.0.0"
RELEASE_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${RELEASE_TAG}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" 2>/dev/null && pwd || echo "$HOME")"

USER_HOME="$HOME"
PREFIX_DIR="$USER_HOME/.local/share/wineprefixes/botw-multiplayer"
DRIVE_C="$PREFIX_DIR/drive_c"
LAUNCHER_DIR="$USER_HOME/Zelda_BotW_Multiplayer"
TEMP_DL="/tmp/botw_mp_setup"
mkdir -p "$TEMP_DL" "$LAUNCHER_DIR" "$USER_HOME/.local/bin"

echo "========================================================================"
echo " Zelda: Breath of the Wild Multiplayer - Linux Installer"
echo "========================================================================"
echo ""

# 1. Dependency checks
echo "[*] Checking system dependencies..."
command -v wine >/dev/null 2>&1 || { echo "[ERROR] 'wine' is required. Please install wine from your package manager."; exit 1; }
command -v curl >/dev/null 2>&1 || { echo "[ERROR] 'curl' is required. Please install curl."; exit 1; }
command -v 7z >/dev/null 2>&1 || command -v 7za >/dev/null 2>&1 || { echo "[ERROR] '7z' (p7zip) is required."; exit 1; }

# Determine terminal for desktop launcher
TERMINAL_BIN="xterm"
for t in kitty foot ghostty alacritty konsole gnome-terminal xfce4-terminal; do
    if command -v "$t" >/dev/null 2>&1; then
        TERMINAL_BIN="$t"
        break
    fi
done

# 2. Game path (optional at install time, configurable anytime in TUI)
GAME_PATH=""
if [ -n "$1" ]; then
    GAME_PATH="$1"
elif [ -t 0 ]; then
    echo ""
    echo "Tip: You can specify your base game directory now, or configure it"
    echo "visually later using the TUI Manager (Pressing 'f' to browse)."
    echo ""
    if command -v zenity >/dev/null 2>&1 && [ -n "$DISPLAY$WAYLAND_DISPLAY" ]; then
        if zenity --question --title="Zelda BotW Multiplayer" --text="Would you like to select your Base Game folder now with the file picker?" --ok-label="Select Folder" --cancel-label="Configure Later in TUI" 2>/dev/null; then
            GAME_PATH="$(zenity --file-selection --directory --title="Select Zelda: Breath of the Wild Game Directory" 2>/dev/null || true)"
        fi
    fi

    if [ -z "$GAME_PATH" ]; then
        read -r -p "Base Game path (leave empty to configure later in TUI): " INPUT_PATH
        if [ -n "$INPUT_PATH" ] && [ -d "$INPUT_PATH" ]; then
            GAME_PATH="$INPUT_PATH"
        fi
    fi
fi

# 3. Create Wineprefix (64-bit)
echo ""
echo "[*] Setting up Wineprefix at: $PREFIX_DIR"
mkdir -p "$PREFIX_DIR"
WINEARCH=win64 WINEPREFIX="$PREFIX_DIR" WINEDEBUG=-all wineboot -u

# 4. Download & install runtimes (.NET 8 Desktop x64 and Visual C++ x64)
echo "[*] Installing Microsoft Visual C++ 2015-2022 Redistributable (x64)..."
VC_FILE="$SCRIPT_DIR/VC_redist.x64.exe"
if [ ! -f "$VC_FILE" ]; then
    VC_FILE="$TEMP_DL/VC_redist.x64.exe"
    if [ ! -f "$VC_FILE" ]; then
        curl -fSL "$RELEASE_URL/VC_redist.x64.exe" -o "$VC_FILE" || curl -fSL "https://aka.ms/vs/17/release/vc_redist.x64.exe" -o "$VC_FILE"
    fi
fi
WINEARCH=win64 WINEPREFIX="$PREFIX_DIR" WINEDEBUG=-all wine "$VC_FILE" /install /quiet /norestart || true

echo "[*] Installing Microsoft .NET 8.0 Desktop Runtime (x64)..."
DOTNET_FILE="$SCRIPT_DIR/windowsdesktop-runtime-8.0.31-win-x64.exe"
if [ ! -f "$DOTNET_FILE" ]; then
    DOTNET_FILE="$TEMP_DL/windowsdesktop-runtime-8.0.31-win-x64.exe"
    if [ ! -f "$DOTNET_FILE" ]; then
        curl -fSL "$RELEASE_URL/windowsdesktop-runtime-8.0.31-win-x64.exe" -o "$DOTNET_FILE" || curl -fSL "https://builds.dotnet.microsoft.com/dotnet/WindowsDesktop/8.0.31/windowsdesktop-runtime-8.0.31-win-x64.exe" -o "$DOTNET_FILE"
    fi
fi
WINEARCH=win64 WINEPREFIX="$PREFIX_DIR" WINEDEBUG=-all wine "$DOTNET_FILE" /install /quiet /norestart || true

# 5. Extract Cemu 1.26.2 and Milk Bar Launcher
echo "[*] Deploying Cemu 1.26.2 (pre-configured with merged multiplayer graphic packs & 32 Link models)..."
mkdir -p "$DRIVE_C/cemu_1.26.2" "$DRIVE_C/MilkBarLauncher" "$DRIVE_C/Games"

CEMU_ARCHIVE="$SCRIPT_DIR/cemu_1.26.2.7z"
if [ ! -f "$CEMU_ARCHIVE" ]; then
    CEMU_ARCHIVE="$TEMP_DL/cemu_1.26.2.7z"
    if [ ! -f "$CEMU_ARCHIVE" ]; then
        echo "    Downloading Cemu package from GitHub Release..."
        curl -fSL "$RELEASE_URL/cemu_1.26.2.7z" -o "$CEMU_ARCHIVE"
    fi
fi
7z x -y "$CEMU_ARCHIVE" -o"$DRIVE_C/" >/dev/null

echo "[*] Deploying Milk Bar Launcher & Dedicated Server..."
MBL_ARCHIVE="$SCRIPT_DIR/MilkBarLauncher.zip"
if [ ! -f "$MBL_ARCHIVE" ]; then
    MBL_ARCHIVE="$TEMP_DL/MilkBarLauncher.zip"
    if [ ! -f "$MBL_ARCHIVE" ]; then
        echo "    Downloading Milk Bar Launcher package from GitHub Release..."
        curl -fSL "$RELEASE_URL/MilkBarLauncher.zip" -o "$MBL_ARCHIVE"
    fi
fi
7z x -y "$MBL_ARCHIVE" -o"$DRIVE_C/MilkBarLauncher/" >/dev/null

# 6. Configure dedicated server defaults
SERVER_CONFIG="$DRIVE_C/MilkBarLauncher/DedicatedServer/ServerConfig.ini"
if [ -f "$SERVER_CONFIG" ]; then
    sed -i 's/^IP=localhost/IP=127.0.0.1/g' "$SERVER_CONFIG"
    sed -i 's/^DefaultGamemode=False/DefaultGamemode=True/g' "$SERVER_CONFIG"
fi

# 7. Configure BOTWM Roaming data
BOTWM_ROAMING="$DRIVE_C/users/$USER/AppData/Roaming/BOTWM"
mkdir -p "$BOTWM_ROAMING"
SERVER_RES="$DRIVE_C/MilkBarLauncher/DedicatedServer"
[ -f "$SERVER_RES/BOTWM.DedicatedServer.AppdataFiles.QuestFlagsNames.txt" ] && cp "$SERVER_RES/BOTWM.DedicatedServer.AppdataFiles.QuestFlagsNames.txt" "$BOTWM_ROAMING/QuestFlagsNames.txt"
[ -f "$SERVER_RES/BOTWM.DedicatedServer.AppdataFiles.ArmorMapping.txt" ] && cp "$SERVER_RES/BOTWM.DedicatedServer.AppdataFiles.ArmorMapping.txt" "$BOTWM_ROAMING/ArmorMapping.txt"

# 8. Configure BCML settings and compatibility merged link
BCML_LOCAL="$DRIVE_C/users/$USER/AppData/Local/bcml"
mkdir -p "$BCML_LOCAL" "$USER_HOME/.config/bcml" "$DRIVE_C/cemu_1.26.2/bcml"

GAME_FOLDER_NAME="The Legend of Zelda Breath of the Wild"
if [ -n "$GAME_PATH" ] && [ -d "$GAME_PATH" ]; then
    # Strip /content if path points inside
    CLEAN_PATH="$GAME_PATH"
    if [ "$(basename "$CLEAN_PATH")" = "content" ] || [ "$(basename "$CLEAN_PATH")" = "code" ]; then
        CLEAN_PATH="$(dirname "$CLEAN_PATH")"
    fi
    GAME_FOLDER_NAME="$(basename "$CLEAN_PATH")"
    ln -sfn "$CLEAN_PATH" "$DRIVE_C/Games/$GAME_FOLDER_NAME"
fi

cat <<EOF > "$BCML_LOCAL/settings.json"
{
  "cemu_dir": "C:/cemu_1.26.2",
  "game_dir": "C:/Games/$GAME_FOLDER_NAME/content",
  "game_dir_nx": "",
  "update_dir": "C:/cemu_1.26.2/mlc01/usr/title/0005000e/101c9400/content",
  "dlc_dir": "C:/cemu_1.26.2/mlc01/usr/title/0005000c/101c9400/content/0010",
  "dlc_dir_nx": "",
  "store_dir": "C:/users/$USER/AppData/Local/bcml",
  "export_dir": "C:/cemu_1.26.2/graphicPacks/BreathOfTheWild_BCML",
  "export_dir_nx": "",
  "load_reverse": false,
  "site_meta": "",
  "no_guess": false,
  "lang": "USen",
  "no_cemu": false,
  "wiiu": true,
  "no_hardlinks": true,
  "force_7z": false,
  "suppress_update": false,
  "loaded": true,
  "nsfw": false,
  "changelog": true,
  "strip_gfx": false,
  "auto_gb": true,
  "show_gb": false,
  "dark_theme": false,
  "last_version": "3.10.8"
}
EOF
cp "$BCML_LOCAL/settings.json" "$USER_HOME/.config/bcml/settings.json"
cp "$BCML_LOCAL/settings.json" "$DRIVE_C/cemu_1.26.2/bcml/settings.json"
ln -sfn "$DRIVE_C/cemu_1.26.2/graphicPacks/BreathOfTheWild_BCML" "$BCML_LOCAL/merged"

# 9. Configure Cemu settings.xml
CEMU_SETTINGS="$DRIVE_C/cemu_1.26.2/settings.xml"
if [ -f "$CEMU_SETTINGS" ]; then
    sed -i 's|<Entry>C:/Users/IK/Roms/WiiU/Base Games/The Legend of Zelda Breath of the Wild \[ALZE01\]/code/U-King.rpx</Entry>|<Entry>C:\\Games\\'"$GAME_FOLDER_NAME"'\\code\\U-King.rpx</Entry>|g' "$CEMU_SETTINGS"
    sed -i 's|<Entry>C:\\Users\\IK\\Roms\\WiiU\\Base Games</Entry>|<Entry>C:\\Games</Entry>|g' "$CEMU_SETTINGS"
    sed -i 's|<path>C:\\Users\\IK\\Roms\\WiiU\\Base Games\\The Legend of Zelda Breath of the Wild \[ALZE01\]\\code\\U-King.rpx</path>|<path>C:\\Games\\'"$GAME_FOLDER_NAME"'\\code\\U-King.rpx</path>|g' "$CEMU_SETTINGS"
    sed -i 's|<Entry filename="\.\.\\\.\.\\AppData\\Local\\bcml\\merged\\rules\.txt"/>|<Entry filename="graphicPacks\\BreathOfTheWild_BCML\\rules.txt"/>|g' "$CEMU_SETTINGS"
fi

# 10. Install / Build TUI Manager (botw-manager)
echo "[*] Installing TUI Manager (botw-manager)..."
MANAGER_INSTALLED=false

if command -v go >/dev/null 2>&1 && [ -f "$SCRIPT_DIR/main.go" ]; then
    echo "    Compiling botw-manager from Go source..."
    (cd "$SCRIPT_DIR" && go build -ldflags="-s -w" -o "$LAUNCHER_DIR/botw-manager" .) && MANAGER_INSTALLED=true || true
fi

if [ "$MANAGER_INSTALLED" = false ] && [ -f "$SCRIPT_DIR/botw-manager" ]; then
    cp -f "$SCRIPT_DIR/botw-manager" "$LAUNCHER_DIR/botw-manager"
    MANAGER_INSTALLED=true
fi

if [ "$MANAGER_INSTALLED" = false ]; then
    echo "    Downloading precompiled botw-manager binary from GitHub Releases..."
    MANAGER_TAR="$TEMP_DL/botw-manager-linux-amd64.tar.gz"
    if curl -fSL "$RELEASE_URL/botw-manager-linux-amd64.tar.gz" -o "$MANAGER_TAR"; then
        tar -xzf "$MANAGER_TAR" -C "$LAUNCHER_DIR/"
        MANAGER_INSTALLED=true
    fi
fi

if [ -f "$LAUNCHER_DIR/botw-manager" ]; then
    chmod +x "$LAUNCHER_DIR/botw-manager"
    install -m 755 "$LAUNCHER_DIR/botw-manager" "$USER_HOME/.local/bin/botw-manager"
fi

# 11. Generate launcher scripts
cat <<EOF > "$LAUNCHER_DIR/0_manager_tui.sh"
#!/usr/bin/env bash
if command -v botw-manager >/dev/null 2>&1; then
    exec botw-manager
elif [ -f "$LAUNCHER_DIR/botw-manager" ]; then
    exec "$LAUNCHER_DIR/botw-manager"
fi
EOF

cat <<EOF > "$LAUNCHER_DIR/1_iniciar_servidor.sh"
#!/usr/bin/env bash
export WINEARCH=win64
export WINEPREFIX="$PREFIX_DIR"
export WINEDEBUG=-all
cd "\$WINEPREFIX/drive_c/MilkBarLauncher/DedicatedServer" || exit 1
wine "\$WINEPREFIX/drive_c/MilkBarLauncher/DedicatedServer/MBL.DedicatedServer.exe"
EOF

cat <<EOF > "$LAUNCHER_DIR/2_iniciar_milkbar.sh"
#!/usr/bin/env bash
export WINEARCH=win64
export WINEPREFIX="$PREFIX_DIR"
export WINEDEBUG=-all
cd "\$WINEPREFIX/drive_c/MilkBarLauncher" || exit 1
wine "\$WINEPREFIX/drive_c/MilkBarLauncher/Milk Bar Launcher.exe"
EOF

cat <<EOF > "$LAUNCHER_DIR/3_iniciar_cemu.sh"
#!/usr/bin/env bash
export WINEARCH=win64
export WINEPREFIX="$PREFIX_DIR"
export WINEDEBUG=-all
cd "\$WINEPREFIX/drive_c/cemu_1.26.2" || exit 1
wine "\$WINEPREFIX/drive_c/cemu_1.26.2/Cemu.exe"
EOF

chmod +x "$LAUNCHER_DIR"/*.sh

# 12. Create Desktop Applications Menu entry
DESKTOP_DIR="$USER_HOME/.local/share/applications"
mkdir -p "$DESKTOP_DIR"

case "$TERMINAL_BIN" in
    foot)
        EXEC_CMD="foot -T \"Zelda BotW Multiplayer Manager\" $USER_HOME/.local/bin/botw-manager"
        ;;
    kitty)
        EXEC_CMD="kitty -T \"Zelda BotW Multiplayer Manager\" $USER_HOME/.local/bin/botw-manager"
        ;;
    ghostty)
        EXEC_CMD="ghostty --title=\"Zelda BotW Multiplayer Manager\" -e $USER_HOME/.local/bin/botw-manager"
        ;;
    alacritty)
        EXEC_CMD="alacritty -T \"Zelda BotW Multiplayer Manager\" -e $USER_HOME/.local/bin/botw-manager"
        ;;
    konsole)
        EXEC_CMD="konsole -p tabtitle=\"Zelda BotW Multiplayer Manager\" -e $USER_HOME/.local/bin/botw-manager"
        ;;
    gnome-terminal)
        EXEC_CMD="gnome-terminal --title=\"Zelda BotW Multiplayer Manager\" -- $USER_HOME/.local/bin/botw-manager"
        ;;
    xfce4-terminal)
        EXEC_CMD="xfce4-terminal --title=\"Zelda BotW Multiplayer Manager\" -e $USER_HOME/.local/bin/botw-manager"
        ;;
    *)
        EXEC_CMD="$TERMINAL_BIN -e $USER_HOME/.local/bin/botw-manager"
        ;;
esac

cat <<EOF > "$DESKTOP_DIR/zelda-botw-manager.desktop"
[Desktop Entry]
Name=Zelda BotW Multiplayer Manager
Comment=Interactive TUI Manager for Zelda BotW Multiplayer on Linux
Exec=$EXEC_CMD
Icon=utilities-terminal
Terminal=false
Type=Application
Categories=Game;
EOF

echo ""
echo "========================================================================"
echo " [SUCCESS] Installation completed successfully!"
echo "========================================================================"
echo " - TUI Manager command:  botw-manager"
echo " - Launcher Directory:   $LAUNCHER_DIR"
echo " - Desktop Entry:        $DESKTOP_DIR/zelda-botw-manager.desktop"
echo ""
echo " You can now start the manager by running:  botw-manager"
echo "========================================================================"

if [ -t 0 ]; then
    read -r -p "Would you like to open the TUI Manager now? [Y/n] " START_NOW
    case "$START_NOW" in
        [nN][oO]|[nN]) ;;
        *)
            exec "$USER_HOME/.local/bin/botw-manager"
            ;;
    esac
fi
