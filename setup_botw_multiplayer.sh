#!/usr/bin/env bash
# ==============================================================================
# Script de Instalacion y Configuracion: Zelda BotW Multiplayer en Linux
# Repositorio: https://github.com/CarlosEvCode/botw-multiplayer-linux
# ==============================================================================
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
USER_HOME="$HOME"
PREFIX_DIR="$USER_HOME/.local/share/wineprefixes/botw-multiplayer"
DRIVE_C="$PREFIX_DIR/drive_c"
LAUNCHER_DIR="$USER_HOME/Zelda_BotW_Multiplayer"
RELEASE_URL="https://github.com/CarlosEvCode/botw-multiplayer-linux/releases/latest/download"

echo "========================================================================"
echo " Instalacion de The Legend of Zelda: Breath of the Wild Multiplayer (Linux)"
echo "========================================================================"

# Verificacion de herramientas del sistema
command -v wine >/dev/null 2>&1 || { echo "[ERROR] 'wine' no se encuentra instalado. Por favor instalalo antes de continuar."; exit 1; }
command -v 7z >/dev/null 2>&1 || { echo "[ERROR] '7z' (p7zip) no se encuentra instalado."; exit 1; }
command -v curl >/dev/null 2>&1 || { echo "[ERROR] 'curl' no se encuentra instalado."; exit 1; }

# Obtencion de la ruta del juego base
if [ -z "$1" ]; then
    echo ""
    echo "Ingrese la ruta absoluta del juego base BotW (desempaquetado):"
    echo "Ejemplo: /home/$USER/Juegos/The Legend of Zelda Breath of the Wild [ALZE01]"
    read -r -p "Ruta del juego: " GAME_PATH
else
    GAME_PATH="$1"
fi

if [ ! -d "$GAME_PATH" ]; then
    echo "[ERROR] El directorio del juego '$GAME_PATH' no existe."
    exit 1
fi

echo "[INFO] Ruta del juego: $GAME_PATH"

# 1. Creacion e inicializacion del Wineprefix de 64 bits
echo "[INFO] Creando prefijo Wine en $PREFIX_DIR..."
mkdir -p "$PREFIX_DIR"
WINEARCH=win64 WINEPREFIX="$PREFIX_DIR" wineboot -u

# 2. Descarga e instalacion de dependencias (.NET 8 Desktop y Visual C++)
TEMP_DL="/tmp/botw_mp_setup"
mkdir -p "$TEMP_DL"

echo "[INFO] Verificando e instalando Microsoft Visual C++ Redistributable (x64)..."
VC_FILE="$SCRIPT_DIR/VC_redist.x64.exe"
if [ ! -f "$VC_FILE" ]; then
    VC_FILE="$TEMP_DL/VC_redist.x64.exe"
    if [ ! -f "$VC_FILE" ]; then
        echo "[INFO] Descargando VC_redist.x64.exe..."
        curl -L "$RELEASE_URL/VC_redist.x64.exe" -o "$VC_FILE" || curl -L "https://aka.ms/vs/17/release/vc_redist.x64.exe" -o "$VC_FILE"
    fi
fi
WINEARCH=win64 WINEPREFIX="$PREFIX_DIR" wine "$VC_FILE" /install /quiet /norestart

echo "[INFO] Verificando e instalando Microsoft .NET 8.0 Desktop Runtime (x64)..."
DOTNET_FILE="$SCRIPT_DIR/windowsdesktop-runtime-8.0.31-win-x64.exe"
if [ ! -f "$DOTNET_FILE" ]; then
    DOTNET_FILE="$TEMP_DL/windowsdesktop-runtime-8.0.31-win-x64.exe"
    if [ ! -f "$DOTNET_FILE" ]; then
        echo "[INFO] Descargando .NET 8.0 Desktop Runtime..."
        curl -L "$RELEASE_URL/windowsdesktop-runtime-8.0.31-win-x64.exe" -o "$DOTNET_FILE" || curl -L "https://builds.dotnet.microsoft.com/dotnet/WindowsDesktop/8.0.31/windowsdesktop-runtime-8.0.31-win-x64.exe" -o "$DOTNET_FILE"
    fi
fi
WINEARCH=win64 WINEPREFIX="$PREFIX_DIR" wine "$DOTNET_FILE" /install /quiet /norestart

# 3. Descarga y descompresion de Cemu 1.26.2 y MilkBarLauncher
echo "[INFO] Extrayendo componentes a drive_c..."
mkdir -p "$DRIVE_C/cemu_1.26.2" "$DRIVE_C/MilkBarLauncher" "$DRIVE_C/Games"

CEMU_ARCHIVE="$SCRIPT_DIR/cemu_1.26.2.7z"
if [ ! -f "$CEMU_ARCHIVE" ]; then
    CEMU_ARCHIVE="$SCRIPT_DIR/cemu_1.26.2_lite.7z"
fi
if [ ! -f "$CEMU_ARCHIVE" ]; then
    CEMU_ARCHIVE="$TEMP_DL/cemu_1.26.2.7z"
    if [ ! -f "$CEMU_ARCHIVE" ]; then
        echo "[INFO] Descargando paquete Cemu preconfigurado desde Releases..."
        curl -L "$RELEASE_URL/cemu_1.26.2.7z" -o "$CEMU_ARCHIVE"
    fi
fi
echo "[INFO] Extrayendo Cemu..."
7z x -y "$CEMU_ARCHIVE" -o"$DRIVE_C/"

MBL_ARCHIVE="$SCRIPT_DIR/MilkBarLauncher.zip"
if [ ! -f "$MBL_ARCHIVE" ]; then
    MBL_ARCHIVE="$TEMP_DL/MilkBarLauncher.zip"
    if [ ! -f "$MBL_ARCHIVE" ]; then
        echo "[INFO] Descargando MilkBarLauncher desde Releases..."
        curl -L "$RELEASE_URL/MilkBarLauncher.zip" -o "$MBL_ARCHIVE"
    fi
fi
echo "[INFO] Extrayendo MilkBarLauncher..."
7z x -y "$MBL_ARCHIVE" -o"$DRIVE_C/MilkBarLauncher/"

# 4. Enlace del juego base
GAME_FOLDER_NAME="$(basename "$GAME_PATH")"
ln -sfn "$GAME_PATH" "$DRIVE_C/Games/$GAME_FOLDER_NAME"

# 5. Ajuste de ServerConfig.ini
SERVER_CONFIG="$DRIVE_C/MilkBarLauncher/DedicatedServer/ServerConfig.ini"
if [ -f "$SERVER_CONFIG" ]; then
    sed -i 's/^IP=localhost/IP=127.0.0.1/g' "$SERVER_CONFIG"
fi

# 6. Copia de recursos de Roaming para BOTWM
BOTWM_ROAMING="$DRIVE_C/users/$USER/AppData/Roaming/BOTWM"
mkdir -p "$BOTWM_ROAMING"
SERVER_RES="$DRIVE_C/MilkBarLauncher/DedicatedServer"
if [ -f "$SERVER_RES/BOTWM.DedicatedServer.AppdataFiles.QuestFlagsNames.txt" ]; then
    cp "$SERVER_RES/BOTWM.DedicatedServer.AppdataFiles.QuestFlagsNames.txt" "$BOTWM_ROAMING/QuestFlagsNames.txt"
fi
if [ -f "$SERVER_RES/BOTWM.DedicatedServer.AppdataFiles.ArmorMapping.txt" ]; then
    cp "$SERVER_RES/BOTWM.DedicatedServer.AppdataFiles.ArmorMapping.txt" "$BOTWM_ROAMING/ArmorMapping.txt"
fi

# 7. Configuracion de BCML settings.json y enlace de validacion de modelos
BCML_LOCAL="$DRIVE_C/users/$USER/AppData/Local/bcml"
mkdir -p "$BCML_LOCAL" "$USER_HOME/.config/bcml" "$DRIVE_C/cemu_1.26.2/bcml"

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

# Enlace de compatibilidad para validacion interna de Milk Bar Launcher
ln -sfn "$DRIVE_C/cemu_1.26.2/graphicPacks/BreathOfTheWild_BCML" "$BCML_LOCAL/merged"

# 8. Ajuste de settings.xml de Cemu
CEMU_SETTINGS="$DRIVE_C/cemu_1.26.2/settings.xml"
if [ -f "$CEMU_SETTINGS" ]; then
    sed -i 's|<Entry>C:/Users/IK/Roms/WiiU/Base Games/The Legend of Zelda Breath of the Wild \[ALZE01\]/code/U-King.rpx</Entry>|<Entry>C:\\Games\\'"$GAME_FOLDER_NAME"'\\code\\U-King.rpx</Entry>|g' "$CEMU_SETTINGS"
    sed -i 's|<Entry>C:\\Users\\IK\\Roms\\WiiU\\Base Games</Entry>|<Entry>C:\\Games</Entry>|g' "$CEMU_SETTINGS"
    sed -i 's|<path>C:\\Users\\IK\\Roms\\WiiU\\Base Games\\The Legend of Zelda Breath of the Wild \[ALZE01\]\\code\\U-King.rpx</path>|<path>C:\\Games\\'"$GAME_FOLDER_NAME"'\\code\\U-King.rpx</path>|g' "$CEMU_SETTINGS"
    sed -i 's|<Entry filename="\.\.\\\.\.\\AppData\\Local\\bcml\\merged\\rules\.txt"/>|<Entry filename="graphicPacks\\BreathOfTheWild_BCML\\rules.txt"/>|g' "$CEMU_SETTINGS"
fi

# 9. Creacion de scripts de ejecucion
mkdir -p "$LAUNCHER_DIR"

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
export WINE_FULLSCREEN_FSR=1
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

# 10. Creacion de entradas de escritorio (.desktop)
DESKTOP_DIR="$USER_HOME/.local/share/applications"
mkdir -p "$DESKTOP_DIR"

cat <<EOF > "$DESKTOP_DIR/zelda-botw-milkbar.desktop"
[Desktop Entry]
Name=Zelda BotW Multiplayer (Milk Bar Launcher)
Comment=Lanzador de The Legend of Zelda: Breath of the Wild Multiplayer
Exec=$LAUNCHER_DIR/2_iniciar_milkbar.sh
Icon=wine
Terminal=false
Type=Application
Categories=Game;
EOF

cat <<EOF > "$DESKTOP_DIR/zelda-botw-server.desktop"
[Desktop Entry]
Name=Zelda BotW Dedicated Server (Multiplayer)
Comment=Servidor dedicado para The Legend of Zelda: Breath of the Wild Multiplayer
Exec=foot $LAUNCHER_DIR/1_iniciar_servidor.sh
Icon=wine
Terminal=true
Type=Application
Categories=Game;
EOF

cat <<EOF > "$DESKTOP_DIR/zelda-botw-cemu.desktop"
[Desktop Entry]
Name=Zelda BotW Cemu (1.26.2 MP)
Comment=Emulador Cemu para Breath of the Wild Multiplayer
Exec=$LAUNCHER_DIR/3_iniciar_cemu.sh
Icon=wine
Terminal=false
Type=Application
Categories=Game;
EOF

echo "========================================================================"
echo "[SUCCESS] Instalacion completada exitosamente."
echo "Scripts de ejecucion generados en: $LAUNCHER_DIR"
echo "========================================================================"
