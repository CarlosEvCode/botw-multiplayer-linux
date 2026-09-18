# The Legend of Zelda: Breath of the Wild Multiplayer on Linux

[English](README.md) | [Español](README.es.md)

This repository provides an automated installation script and comprehensive technical documentation to run the multiplayer mod (*Milk Bar Launcher* / *Breath of the Wild Multiplayer*) on GNU/Linux distributions and SteamOS via Wine.

---

## 1. System Architecture

The multiplayer mod architecture consists of three core components:

1. **Cemu Emulator (v1.26.2 x86_64 for Windows)**:
   - Emulates the Wii U release of *The Legend of Zelda: Breath of the Wild*.
   - Requires specific merged Graphic Packs enabled (`BreathOfTheWild_BCML`, `bcmlPatches/MilkBarLauncher`, and `ExtendedMemory`).

2. **Milk Bar Launcher (.NET / WinUI Client)**:
   - Handles player authentication, server discovery, and connection state.
   - Injects `InjectDLL.dll` into the Cemu runtime process and communicates via Named Pipes to read and synchronize game memory.

3. **MBL Dedicated Server (.NET Dedicated Server)**:
   - Authoritative game server synchronizing player positions, animations, inventories, quests, and world states over UDP/TCP port `5050`.

---

## 2. Linux & Wine Compatibility Analysis and Fixes

Through environment debugging and process tracing on Linux, five critical failure points were identified and resolved:

### A. Process Isolation in Wine
* **Issue:** When Cemu and Milk Bar Launcher run in separate Wineprefixes or standalone containerized environments, Milk Bar Launcher cannot access Cemu memory spaces or Named Pipes, causing indefinite hangs (*white screen "Cemu 1.26.2f - loading..."*).
* **Fix:** Both Cemu and Milk Bar Launcher must strictly reside and execute within the same 64-bit Wineprefix (`~/.local/share/wineprefixes/botw-multiplayer`).

### B. .NET Runtime Requirements
* **Issue:** Modern releases of Milk Bar Launcher require `.NET Desktop Runtime 8.0 (x64)`. The exclusive presence of .NET 6.0 or older runtimes results in immediate termination with exit code 150.
* **Fix:** The setup script automatically provisions `windowsdesktop-runtime-8.0.x-win-x64.exe` and `VC_redist.x64.exe` into the prefix.

### C. BCML Model Validation ("Mod is not setup on BCML")
* **Issue:** Upon clicking *Connect*, Milk Bar Launcher verifies the existence of all 32 player actor packages (`Jugador1.sbactorpack` to `Jugador32.sbactorpack`) in `store_dir/merged/content/Actor/Pack/` as defined in `settings.json`. If this path is missing or unlinked, the client halts connection.
* **Fix:** `settings.json` is generated using internal virtual `C:` drive mappings, and a symbolic link from `BreathOfTheWild_BCML` to `AppData/Local/bcml/merged` is established.

### D. Dedicated Server "localhost" Resolution Crash
* **Issue:** Under Wine, .NET socket binding may throw a `NullReferenceException` when resolving string hostname `localhost`.
* **Fix:** The `IP` configuration parameter in `ServerConfig.ini` is explicitly bound to `127.0.0.1`.

### E. Missing Roaming AppData Resources
* **Issue:** The dedicated server looks for `QuestFlagsNames.txt` and `ArmorMapping.txt` directly inside `%APPDATA%/BOTWM/`. When extracted with assembly prefixes (`BOTWM.DedicatedServer.AppdataFiles.*`), the server aborts with `FileNotFoundException`.
* **Fix:** Required configuration files are automatically sanitized, copied, and placed in `AppData/Roaming/BOTWM/`.

---

## 3. Repository and Release Artifacts

| Resource | Location | Description |
| :--- | :--- | :--- |
| `setup_botw_multiplayer.sh` | Repository | Automated deployment and configuration bash script. |
| `cemu_1.26.2.7z` | GitHub Releases | Preconfigured Cemu 1.26.2 build (~146 MB) containing Graphic Packs, 32 player models, and merged BCML patches. |
| `MilkBarLauncher.zip` | GitHub Releases | Milk Bar Launcher client, injection binaries, and dedicated server. |
| `windowsdesktop-runtime-8.0.31-win-x64.exe` | GitHub Releases | Official Microsoft .NET Desktop Runtime 8.0 x64 installer. |
| `VC_redist.x64.exe` | GitHub Releases | Official Microsoft Visual C++ Redistributable 2015-2022 x64 installer. |

---

## 4. Installation Guide

### System Prerequisites
Install the required base utilities using your distribution package manager:

```bash
# Arch Linux / CachyOS / Manjaro:
sudo pacman -S wine winetricks p7zip curl

# Ubuntu / Debian:
sudo apt install wine winetricks p7zip-full curl

# Fedora:
sudo dnf install wine winetricks p7zip p7zip-plugins curl
```

### Automated Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/CarlosEvCode/botw-multiplayer-linux.git
   cd botw-multiplayer-linux
   chmod +x setup_botw_multiplayer.sh
   ```

2. Run the installer specifying the absolute path to your unpacked BotW base game folder:
   ```bash
   ./setup_botw_multiplayer.sh "/absolute/path/to/The Legend of Zelda Breath of the Wild"
   ```

*Note: If binary archives are not present locally, the installer will automatically download them from the repository's GitHub Releases.*

The script performs the following tasks automatically:
1. Initializes a 64-bit Wineprefix at `~/.local/share/wineprefixes/botw-multiplayer`.
2. Silently installs Microsoft Visual C++ and .NET 8.0 Desktop Runtime.
3. Extracts and structures Cemu 1.26.2 and Milk Bar Launcher inside `drive_c`.
4. Symlinks the base game, adjusts configuration files, and links BCML merged mod paths.
5. Generates launcher scripts in `~/Zelda_BotW_Multiplayer/` and system desktop shortcuts (`~/.local/share/applications/`).

---

## 5. Running and Connecting

### Option A: Interactive TUI Manager (Recommended)
You can manage the server, client, routes, and network in a single unified terminal UI:
```bash
botw-manager
# or: ~/Zelda_BotW_Multiplayer/0_manager_tui.sh
```
* **Keybindings:**
  * `[s]` Start/Stop Dedicated Server
  * `[m]` Launch Milk Bar Launcher
  * `[c]` Launch Cemu 1.26.2
  * `[t]` Copy Tailscale / Local IP to clipboard
  * `[r]` Game paths, Update & DLC configuration tab
  * `[i]` Send console command to the dedicated server
  * `[Tab]` Switch views / tabs
  * `[q]` Quit

### Option B: Standalone Scripts
Individual control scripts are also available in `~/Zelda_BotW_Multiplayer/`:

1. **Start Dedicated Server (Host only):**
   ```bash
   ~/Zelda_BotW_Multiplayer/1_iniciar_servidor.sh
   ```
   - Enter `0` for standard co-op or `1` for special game modes (Prop Hunt, etc.).

2. **Start Client (All players):**
   ```bash
   ~/Zelda_BotW_Multiplayer/2_iniciar_milkbar.sh
   ```
   - **Local Host:** Connect to `127.0.0.1`, port `5050`.
   - **Remote Clients (via VPN such as Tailscale / ZeroTier):** Connect to Host VPN IP, port `5050`.

3. **Cemu Configuration (Controllers / Graphics):**
   ```bash
   ~/Zelda_BotW_Multiplayer/3_iniciar_cemu.sh
   ```

---

## 6. Credits and Attribution

- **Milk Bar Launcher & Dedicated Server**: Developed and maintained by the [MilkBarModding](https://github.com/MilkBarModding/MilkBarLauncher) community.
- **Cemu (Wii U Emulator)**: Developed by [Team Cemu](https://cemu.info/) under the Mozilla Public License 2.0.
- **BCML (BotW Cross-Platform Mod Loader)**: Created by [NiceneNerd](https://github.com/NiceneNerd/BCML).

### Disclaimer
This repository and its associated resources do not contain, host, or distribute proprietary Nintendo game assets, official executable binaries (`.rpx`), encryption keys, or copyrighted materials. Users must supply their own legally acquired base game, update, and DLC files.
