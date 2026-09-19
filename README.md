# The Legend of Zelda: Breath of the Wild Multiplayer on Linux / SteamOS

[English (README.md)](README.md) | [Español (README.es.md)](README.es.md)

An automated deployment installer, dedicated TUI Manager, and comprehensive technical environment to host and play *The Legend of Zelda: Breath of the Wild Multiplayer* (*Milk Bar Launcher* / *Cemu 1.26.2*) on GNU/Linux and SteamOS via Wine.

---

## 1. Quick Installation

Run this single command in your terminal to download, configure, and install everything automatically:

```bash
curl -sSL https://raw.githubusercontent.com/CarlosEvCode/botw-multiplayer-linux/main/install.sh | bash
```

### Manual Installation (via Git)

```bash
git clone https://github.com/CarlosEvCode/botw-multiplayer-linux.git
cd botw-multiplayer-linux
chmod +x install.sh
./install.sh
```

> **Note:** The installer automatically provisions 64-bit Wine, Microsoft .NET 8.0 Desktop Runtime, Visual C++ 2015–2022, Cemu 1.26.2 (with merged BCML patches and 32 Link player models), Milk Bar Launcher (Server & Headless Client), and the `botw-manager` TUI CLI.

---

## 2. TUI Manager (`botw-manager`)

Manage servers, clients, gamemodes, routes, and network IPs directly from your terminal:

```bash
botw-manager
```

### Key Features & Controls
* **Global Navigation:** Press `Tab` / `Shift+Tab` or number keys `[1]`–`[6]` to switch between tabs at any time.
* **1. Dashboard:**
  * `[s]`: Start / Stop Dedicated Server (live log console stream).
  * `[c]`: Connect to Game (launches Cemu and injects multiplayer mod headlessly via `MilkBar.CLI`).
  * `[x]`: Disconnect / Stop active client session.
  * `[e]`: Launch Cemu 1.26.2 (Standalone / Offline).
  * `[t]`: Copy ZeroTier / VPN / LAN IP to clipboard.
  * `[q]`: Exit manager.
* **2. Client & Connect:**
  * Configure multiplayer connection parameters: Target Server IP, Port, Password, Player Name, and Character Model.
  * `[c]` / `[Enter]`: Connect to Game.
  * `[l]`: Set target IP to `127.0.0.1` (Local Host).
  * `[t]`: Set target IP to detected ZeroTier / VPN / LAN IP.
  * `[s]`: Save connection parameters.
  * `[x]`: Disconnect active client.
* **3. Gamemodes:**
  * `←` / `→`: Cycle game modes:
    * **Standard Co-op (Free):** Full story co-op with customizable sync options.
    * **Hunter vs Speedrunner:** Competitive Manhunt mode (automatically enforces individual player progress).
    * **DeathSwap:** High-stakes survival swap mode (auto-locks survival rules).
  * `Space`: Toggle individual synchronization rules (Quests, Shrines, Towers, Koroks, Enemies, Dungeons, Locations).
  * `s`: Save rules (auto-restarts server if currently active).
* **4. Paths & DLC:**
  * `↓` / `↑`: Select field (Base Game, Update v208, DLC v80).
  * `f`: Open native GUI folder browser (*Zenity* / *Kdialog*) to select folders with one click.
  * `Enter`: Save and rebuild Wine/BCML symbolic links.
* **5. Network:**
  * View active ZeroTier (Primary), Tailscale, and LAN IP addresses with connection guides.
* **6. Settings:**
  * **Language Switcher:** Toggle instantly between **English** (default) and **Español**.
  * Project info, version, and repository links.

---

## 3. Multiplayer Connection Guide

### A. Local Host (Playing on the same PC)
1. Launch `botw-manager` and press `[s]` to start the Dedicated Server.
2. Press `[c]` to connect immediately (or go to Tab `2. Client & Connect` and verify `127.0.0.1:5050`).
3. Cemu will launch automatically with multiplayer active, without any background Wine GUI windows.

### B. Local Network (LAN / Same Wi-Fi or Ethernet)
1. The Host starts the Dedicated Server in `botw-manager` (press `[s]`).
2. The Host shares their **Local LAN IP** shown in `botw-manager` (e.g., `192.168.1.50`).
3. Other players open `botw-manager` on their PC, go to Tab `2. Client & Connect`, set Target IP to the Host's LAN IP (`192.168.1.50`), and press `[c]`.

### C. Internet / Remote Friends (ZeroTier VPN - Recommended)
1. Create a shared virtual network room on [ZeroTier](https://www.zerotier.com/) and have all players join the network ID (`zerotier-cli join <network-id>`).
2. The Host starts the Dedicated Server and copies their ZeroTier IP (press `[t]` in `botw-manager`).
3. Remote friends open `botw-manager`, go to Tab `2. Client & Connect`, paste the Host's ZeroTier IP, and press `[c]`.
*(Tailscale is also supported if you prefer mesh networking across your personal devices).*

---

## 4. Headless Client Architecture & Wayland Stability

The original Windows Milk Bar Launcher used a WPF interface with a 50ms polling timer (`CemuFollower`) that attempted to force reparenting of Win32 windows over Cemu. Under XWayland and tiling window managers (such as Hyprland, Sway, or SteamOS Game Mode), this caused severe window flickering and focus thrashing.

This suite solves the problem by integrating **`MilkBar.CLI`**, a headless mod client developed in [.NET 8](https://github.com/CarlosEvCode/MilkBarLauncher):
* **Zero GUI Overhead:** All mod injection, Named Pipe IPC, and Cemu process management run cleanly in the background.
* **Wayland & SteamOS Native:** No Win32 window reparenting conflicts, no XWayland focus stealing, and complete stability.

Optional Hyprland window rule for Cemu:
```lua
o.window({ class = "^(cemu\\.exe)$" }, {
  opaque = true,
  suppress_event = "maximize",
})
```

---

## 5. Technical Architecture & Linux Compatibility Notes

For reference, the following technical solutions were implemented to ensure stability under Wine on Linux:

1. **Headless Mod Client (`MilkBar.CLI.exe`):**
   * Spawns Cemu, reads configuration parameters directly, initializes Named Pipes, and injects `InjectDLL.dll` without any WPF GUI dependencies.
2. **Shared 64-bit Wineprefix (`~/.local/share/wineprefixes/botw-multiplayer`):**
   * Cemu, the Dedicated Server, and `MilkBar.CLI` run within the exact same Wineprefix so `InjectDLL.dll` and Named Pipes can hook into Cemu's game memory without white screen hangs.
3. **Microsoft .NET 8.0 Desktop Runtime x64:**
   * Installed directly into the Wineprefix to resolve exit code 150 crashes.
4. **BCML Model Package Validation:**
   * Pre-configured with all 32 Link player model packages (`Jugador1.sbactorpack` to `Jugador32.sbactorpack`) in `store_dir/merged/content/Actor/Pack/` and linked to `BreathOfTheWild_BCML`.
5. **Dedicated Server Socket Binding:**
   * `ServerConfig.ini` defaults to `IP=127.0.0.1` and `DefaultGamemode=True` to immediately bind socket `5050` and prevent `localhost` DNS resolve exceptions.
6. **AppData Roaming Resource Sanitization:**
   * Extracted `QuestFlagsNames.txt` and `ArmorMapping.txt` into `%APPDATA%/BOTWM/` to prevent `FileNotFoundException` during server start.

---

## 6. Uninstallation

To completely uninstall and remove all Wineprefix files, launchers, configurations, and binaries cleanly from your system, run:

```bash
botw-manager --uninstall
```

Or via the installer script:
```bash
./install.sh --uninstall
```

> **Note:** Uninstallation will only delete multiplayer suite files and Wine configurations. Your original game dumps, updates, and DLC remain intact and untouched.

---

## 7. Credits & Attribution

* **Milk Bar Launcher & Dedicated Server:** Developed by the [MilkBarModding](https://github.com/MilkBarModding/MilkBarLauncher) community.
* **MilkBar.CLI (Headless Client):** Maintained and adapted for Linux by [CarlosEvCode](https://github.com/CarlosEvCode/MilkBarLauncher).
* **Cemu (Wii U Emulator):** Developed by [Team Cemu](https://cemu.info/) under Mozilla Public License 2.0.
* **BCML:** Created by [NiceneNerd](https://github.com/NiceneNerd/BCML).
* **Linux Manager & Deployment:** Maintained by [CarlosEvCode](https://github.com/CarlosEvCode).

### Disclaimer
This repository does not contain or distribute proprietary Nintendo game assets, official ROM executables (`.rpx`), encryption keys, or copyrighted materials. Users must provide their own legally dumped copy of *The Legend of Zelda: Breath of the Wild*, update files, and DLC.
