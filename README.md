# The Legend of Zelda: Breath of the Wild Multiplayer on Linux / SteamOS

[English (README.md)](README.md) | [Español (README.es.md)](README.es.md)

An automated deployment installer, dedicated TUI Manager, and comprehensive technical environment to host and play *The Legend of Zelda: Breath of the Wild Multiplayer* (*Milk Bar Launcher* / *Cemu 1.26.2*) on GNU/Linux and SteamOS via Wine.

---

## 🚀 Quick Install (1-Line Command)

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

> **Note:** The installer automatically provisions 64-bit Wine, Microsoft .NET 8.0 Desktop Runtime, Visual C++ 2015–2022, Cemu 1.26.2 (with merged BCML patches & 32 Link player models), Milk Bar Launcher, and the `botw-manager` TUI CLI.

---

## 🎮 TUI Manager (`botw-manager`)

Manage servers, clients, gamemodes, routes, and network IPs directly from your terminal:

```bash
botw-manager
```

### Key Features & Controls
* **Global Navigation:** Press **`Tab`** / **`Shift+Tab`** to switch between tabs at any time.
* **1. Dashboard:**
  * **`[s]`**: Start / Stop Dedicated Server (live log console stream).
  * **`[m]`**: Launch Milk Bar Launcher.
  * **`[c]`**: Launch Cemu 1.26.2.
  * **`[t]`**: Copy Tailscale / LAN IP to clipboard.
  * **`[q]`**: Exit manager.
* **2. Gamemodes:**
  * **`←` / `→`**: Cycle game modes:
    * **Standard Co-op (Free):** Full story co-op with customizable sync options.
    * **Hunter vs Speedrunner:** Competitive Manhunt mode (automatically enforces individual player progress).
    * **DeathSwap:** High-stakes survival swap mode (auto-locks survival rules).
  * **`Space`**: Toggle individual synchronization rules (Quests, Shrines, Towers, Koroks, Enemies, Dungeons, Locations).
  * **`s`**: Save rules (auto-restarts server if currently active).
* **3. Paths & DLC:**
  * **`↓` / `↑`**: Select field (Base Game, Update v208, DLC v80).
  * **`f`**: Open native GUI folder browser (*Zenity* / *Kdialog*) to select folders with one click.
  * **`Enter`**: Save and rebuild Wine/BCML symbolic links.
* **4. Network:**
  * View active Tailscale VPN, LAN, and ZeroTier IP addresses with connection guides.
* **5. Settings:**
  * **Language Switcher:** Toggle instantly between **English** (default) and **Español**.
  * Project info, version, and repository links.

---

## 🌐 Multiplayer Connection Guide

### Host (Player hosting the Server)
1. Launch `botw-manager` and press **`s`** to start the Dedicated Server.
2. Press **`m`** to launch Milk Bar Launcher.
3. In Milk Bar Launcher, set:
   * **IP:** `127.0.0.1`
   * **Port:** `5050`
4. Click **Connect** (Cemu will launch automatically and sync your game).

### Remote Friends (Players joining over the Internet)
1. Install [Tailscale](https://tailscale.com/) (or ZeroTier) on both Host and Client machines.
2. Host copies their Tailscale IP (press **`t`** in `botw-manager`).
3. Remote players open Milk Bar Launcher and enter:
   * **IP:** `Host's Tailscale IP` (e.g., `100.x.y.z`)
   * **Port:** `5050`
4. Click **Connect**.

---

## 🪟 Hyprland & Wayland Window Stability

If you use **Hyprland**, add the following window rules to `~/.config/hypr/hyprland.lua` to prevent XWayland focus flickering and send Milk Bar Launcher silently to background Workspace 3:

```lua
o.window({ class = ".*(milk bar launcher|MilkBar).*" }, {
  float = true,
  center = true,
  size = "1188 670",
  suppress_event = "activate maximize fullscreen",
  workspace = "3 silent",
})
o.window({ class = "^(cemu\\.exe)$" }, {
  opaque = true,
  suppress_event = "maximize",
})
```

---

## 🛠️ Technical Architecture & Linux Compatibility Notes

For reference, the following technical solutions were implemented to ensure stability under Wine on Linux:

1. **Shared 64-bit Wineprefix (`~/.local/share/wineprefixes/botw-multiplayer`):**
   * Cemu and Milk Bar Launcher must run within the exact same Wineprefix so `InjectDLL.dll` and Named Pipes can hook into Cemu's game memory without white screen hangs.
2. **Microsoft .NET 8.0 Desktop Runtime x64:**
   * Installed directly into the Wineprefix to resolve exit code 150 crashes.
3. **BCML Model Package Validation:**
   * Pre-configured with all 32 Link player model packages (`Jugador1.sbactorpack` to `Jugador32.sbactorpack`) in `store_dir/merged/content/Actor/Pack/` and linked to `BreathOfTheWild_BCML`.
4. **Dedicated Server Socket Binding:**
   * `ServerConfig.ini` defaults to `IP=127.0.0.1` and `DefaultGamemode=True` to immediately bind socket `5050` and prevent `localhost` DNS resolve exceptions.
5. **AppData Roaming Resource Sanitization:**
   * Extracted `QuestFlagsNames.txt` and `ArmorMapping.txt` into `%APPDATA%/BOTWM/` to prevent `FileNotFoundException` during server start.

---

## 📜 Credits & Attribution

* **Milk Bar Launcher & Dedicated Server:** Developed by the [MilkBarModding](https://github.com/MilkBarModding/MilkBarLauncher) community.
* **Cemu (Wii U Emulator):** Developed by [Team Cemu](https://cemu.info/) under Mozilla Public License 2.0.
* **BCML:** Created by [NiceneNerd](https://github.com/NiceneNerd/BCML).
* **Linux Manager & Deployment:** Maintained by [CarlosEvCode](https://github.com/CarlosEvCode).

### Disclaimer
This repository does not contain or distribute proprietary Nintendo game assets, official ROM executables (`.rpx`), encryption keys, or copyrighted materials. Users must provide their own legally dumped copy of *The Legend of Zelda: Breath of the Wild*, update files, and DLC.
