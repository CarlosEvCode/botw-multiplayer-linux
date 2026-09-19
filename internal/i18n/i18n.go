package i18n

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Language string

const (
	LangEN Language = "en"
	LangES Language = "es"
)

var CurrentLang Language = LangEN

type Config struct {
	Language Language `json:"language"`
}

func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config/botw-manager/settings.json")
}

func LoadLanguage() Language {
	path := GetConfigPath()
	if data, err := os.ReadFile(path); err == nil {
		var cfg Config
		if err := json.Unmarshal(data, &cfg); err == nil && (cfg.Language == LangEN || cfg.Language == LangES) {
			CurrentLang = cfg.Language
			return CurrentLang
		}
	}
	CurrentLang = LangEN
	return CurrentLang
}

func SaveLanguage(lang Language) error {
	CurrentLang = lang
	path := GetConfigPath()
	_ = os.MkdirAll(filepath.Dir(path), 0755)
	data, err := json.MarshalIndent(Config{Language: lang}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func SetLanguage(lang Language) {
	CurrentLang = lang
	_ = SaveLanguage(lang)
}

func T(key string) string {
	if dict, ok := translations[CurrentLang]; ok {
		if val, exists := dict[key]; exists {
			return val
		}
	}
	// Fallback to English
	if dict, ok := translations[LangEN]; ok {
		if val, exists := dict[key]; exists {
			return val
		}
	}
	return key
}

var translations = map[Language]map[string]string{
	LangEN: {
		"app.title": "BOTW MULTIPLAYER MANAGER",

		// Tabs
		"tab.dashboard": "1. Dashboard",
		"tab.client":    "2. Client & Connect",
		"tab.gamemodes": "3. Gamemodes",
		"tab.paths":     "4. Paths & DLC",
		"tab.network":   "5. Network",
		"tab.settings":  "6. Settings",

		// Dashboard
		"dash.services_title": "SERVICES STATUS",
		"dash.srv_dedicated":  "Dedicated Server",
		"dash.srv_client":     "Multiplayer Game Client",
		"dash.srv_cemu":       "Cemu 1.26.2 (Offline)",
		"dash.srv_mode":       "Server Mode:",
		"dash.net_title":      "CONNECTIVITY",
		"dash.net_zerotier":   "ZeroTier IP:",
		"dash.net_tailscale":  "Tailscale IP:",
		"dash.net_local":      "Local LAN IP:",
		"dash.net_port":       "Port:",
		"dash.net_pass":       "Server Password:",
		"dash.net_nopass":     "No password",
		"dash.net_notfound":   "Not detected",
		"dash.console_title":  "CONSOLE & LOGS",

		// Client Tab
		"client.title":         "MULTIPLAYER CLIENT & DIRECT CONNECTION",
		"client.header_params": "CONNECTION PARAMETERS:",
		"client.target_ip":     "1. Target Server IP",
		"client.target_port":   "2. Server Port",
		"client.password":      "3. Server Password (optional)",
		"client.player_name":   "4. Player Display Name",
		"client.model":         "5. Character Model",
		"client.quick_actions": "ACTIONS & SHORTCUTS:",
		"client.btn_connect":   "[c / Enter: Connect to Game]",
		"client.btn_disconnect":"[x: Disconnect Client]",
		"client.btn_local":     "[l: Set to 127.0.0.1 (Local Host)]",
		"client.btn_vpn":       "[t: Set to ZeroTier / VPN IP]",
		"client.btn_save":      "[s: Save Parameters]",
		"client.status":        "Client Status:",
		"client.status_ready":  "Ready to launch Cemu and connect.",
		"client.status_running":"Multiplayer session active in Cemu.",
		"client.footer_hint":   "[Tab: Tab]  [↓/↑: Select field]  [c/Enter: Connect]  [s: Save]  [l: Local IP]  [t: ZeroTier IP]  [Esc: Back]",

		// Gamemodes
		"gm.title":            "GAMEMODES & DEDICATED SERVER CONFIGURATION",
		"gm.srv_running":      ">> SERVER RUNNING: Changes will apply by restarting the server on save.",
		"gm.srv_stopped":      ">> SERVER STOPPED: Configure your rules and press [s] to start.",
		"gm.special_mode":     "Special Gamemode",
		"gm.mode_coop":        "Standard Co-op (Free)",
		"gm.mode_hunter":      "Hunter vs Speedrunner",
		"gm.mode_deathswap":   "DeathSwap",
		"gm.header_coop":      "SYNCHRONIZATION OPTIONS (FREE / FULLY CUSTOMIZABLE):",
		"gm.header_hunter":    "SYNCHRONIZATION OPTIONS (INDIVIDUAL PROGRESS - HUNTER VS SPEEDRUNNER):",
		"gm.header_deathswap": "SYNCHRONIZATION OPTIONS (DISABLED FOR SURVIVAL - DEATHSWAP):",
		"gm.sync_quest":       "Sync Quests (QuestSync)",
		"gm.sync_shrine":      "Sync Shrines (ShrineSync)",
		"gm.sync_tower":       "Sync Towers (TowerSync)",
		"gm.sync_korok":       "Sync Koroks (KorokSync)",
		"gm.sync_enemy":       "Sync Enemies (EnemySync)",
		"gm.sync_dungeon":     "Sync Dungeons (DungeonSync)",
		"gm.sync_location":    "Sync Locations (LocationSync)",
		"gm.locked":           "[-] Locked",
		"gm.server_params":    "SERVER PARAMETERS:",
		"gm.password":         "Password",
		"gm.description":      "Description",
		"gm.footer_hint":      "[Tab: Tab]  [←/→: Change mode]  [Space: Toggle option]  [s: Save & Apply]  [Esc: Back]",

		// Paths
		"paths.title":        "GAME, UPDATE & DLC PATHS MANAGER",
		"paths.base":         "1. Base Game Directory (Zelda BotW):",
		"paths.update":       "2. Update Directory (Update v208):",
		"paths.dlc":          "3. DLC Directory (v80):",
		"paths.status":       "Status:",
		"paths.browse":       "[f: Browse]",
		"paths.footer_hint":  "[Tab: Tab]  [↓/↑: Change field]  [f: Open GUI Explorer]  [Enter: Save & Apply]  [Esc: Back]",
		"paths.val_not_set":  "Path not configured",
		"paths.val_rpx":      "Valid (U-King.rpx executable found)",
		"paths.val_content":  "Valid (Content folder found)",
		"paths.val_dir":      "Valid (Game directory found)",
		"paths.val_err_base": "U-King.rpx or content folder not found",
		"paths.val_up_cemu":  "Installed in Cemu (v208)",
		"paths.val_up_ext":   "Externally linked",
		"paths.val_up_err":   "Not found (Requires Update v208)",
		"paths.val_dlc_cemu": "Installed in Cemu (v80)",
		"paths.val_dlc_ext":  "Externally linked",
		"paths.val_dlc_err":  "Not found (Requires DLC v80)",

		// Network
		"net.title":        "NETWORK & MULTIPLAYER DETAILS",
		"net.interfaces":   "DETECTED NETWORK INTERFACES:",
		"net.zerotier":     "ZeroTier VPN (Primary):",
		"net.tailscale":    "Tailscale VPN (Secondary):",
		"net.local":        "Local Network (LAN):",
		"net.inactive":     "Inactive / Not Connected",
		"net.instructions": "MULTIPLAYER CONNECTION INSTRUCTIONS:",
		"net.host_info":    "- Host (Server): Connects locally to %s on port %s",
		"net.client_info":  "- Remote Clients: Enter your ZeroTier / VPN IP (%s) in Client Tab",
		"net.port_info":    "- Default Port: %s (TCP/UDP)",
		"net.zt_tip":       "- Recommended: Use ZeroTier shared network ID for multi-user peer rooms.",
		"net.footer_hint":  "[t: Copy ZeroTier/VPN IP]  [Tab: Next tab]",

		// Settings & About
		"set.title":          "SETTINGS",
		"set.pref_title":     "APPLICATION PREFERENCES:",
		"set.lang_label":     "Language / Idioma:",
		"set.about_title":    "ABOUT THIS PROJECT:",
		"set.app_name":      "Zelda BotW Multiplayer Manager for Linux / SteamOS",
		"set.version":        "Version: v1.0.0",
		"set.author":         "Author: CarlosEvCode",
		"set.github":         "GitHub: https://github.com/CarlosEvCode/botw-multiplayer-linux",
		"set.stack":          "Built with: Go + Bubble Tea + Wine + Cemu 1.26.2 + Milk Bar Launcher + BCML",
		"set.desc":           "Full automated manager for hosting, connecting and managing Zelda BotW Multiplayer on Linux.",
		"set.footer_hint":    "[Tab: Tab]  [←/→ / Space: Switch Language]  [Esc: Back]",

		// Footer
		"foot.server":    "Server",
		"foot.client":    "Connect",
		"foot.cemu":      "Cemu",
		"foot.gamemodes": "Gamemodes",
		"foot.paths":     "Paths",
		"foot.copy_ip":   "Copy IP",
		"foot.tab":       "Tab",
		"foot.quit":      "Quit",

		// Notifications
		"notif.srv_stopping":     "Stopping dedicated server...",
		"notif.srv_started":      "Server started on port %s",
		"notif.srv_error":        "Error starting server: %v",
		"notif.client_starting":  "Launching Cemu & connecting to %s:%s...",
		"notif.client_stopping":  "Disconnecting multiplayer client...",
		"notif.client_saved":     "Client connection parameters saved!",
		"notif.client_set_local": "Target IP set to 127.0.0.1 (Local Host)",
		"notif.client_set_vpn":   "Target IP set to ZeroTier/VPN IP: %s",
		"notif.cemu_starting":    "Launching Cemu 1.26.2...",
		"notif.ip_copied":        "Tailscale IP copied: %s",
		"notif.zt_ip_copied":     "ZeroTier IP copied: %s",
		"notif.local_ip_copied":  "Local IP copied: %s",
		"notif.picker_opening":   "Opening folder explorer...",
		"notif.paths_saved":      "Paths updated and linked successfully!",
		"notif.cfg_saved_reload": "Configuration saved! Restarting server with new rules...",
		"notif.cfg_saved":        "Configuration saved! Ready for server launch.",
		"notif.mode_coop":        "Mode: Standard Co-op - Synchronizations unlocked",
		"notif.mode_hunter":      "Mode: Hunter vs Speedrunner - Individual progress enabled",
		"notif.mode_deathswap":   "Mode: DeathSwap - Survival mode enabled",
		"notif.sync_locked":      "Synchronization locked in competitive mode",
		"notif.lang_changed":     "Language switched to English!",
	},
	LangES: {
		"app.title": "BOTW MULTIPLAYER MANAGER",

		// Tabs
		"tab.dashboard": "1. Dashboard",
		"tab.client":    "2. Cliente & Conexión",
		"tab.gamemodes": "3. Modos de Juego",
		"tab.paths":     "4. Rutas & DLC",
		"tab.network":   "5. Red",
		"tab.settings":  "6. Ajustes",

		// Dashboard
		"dash.services_title": "ESTADO DE SERVICIOS",
		"dash.srv_dedicated":  "Servidor Dedicado",
		"dash.srv_client":     "Cliente Multijugador",
		"dash.srv_cemu":       "Cemu 1.26.2 (Sin conexión)",
		"dash.srv_mode":       "Modo Servidor:",
		"dash.net_title":      "CONECTIVIDAD",
		"dash.net_zerotier":   "IP ZeroTier:",
		"dash.net_tailscale":  "IP Tailscale:",
		"dash.net_local":      "IP Local LAN:",
		"dash.net_port":       "Puerto:",
		"dash.net_pass":       "Clave Servidor:",
		"dash.net_nopass":     "Sin clave",
		"dash.net_notfound":   "No detectado",
		"dash.console_title":  "CONSOLA & REGISTROS",

		// Client Tab
		"client.title":         "CONFIGURACION DEL CLIENTE Y CONEXION DIRECTA",
		"client.header_params": "PARAMETROS DE CONEXION:",
		"client.target_ip":     "1. IP del Servidor Destino",
		"client.target_port":   "2. Puerto del Servidor",
		"client.password":      "3. Contrasena del Servidor (opcional)",
		"client.player_name":   "4. Nombre del Jugador",
		"client.model":         "5. Modelo del Personaje",
		"client.quick_actions": "ACCIONES & ATAJOS:",
		"client.btn_connect":   "[c / Enter: Conectar al Juego]",
		"client.btn_disconnect":"[x: Desconectar Cliente]",
		"client.btn_local":     "[l: Fijar a 127.0.0.1 (Anfitrion Local)]",
		"client.btn_vpn":       "[t: Fijar a IP de ZeroTier/VPN]",
		"client.btn_save":      "[s: Guardar Parametros]",
		"client.status":        "Estado del Cliente:",
		"client.status_ready":  "Listo para abrir Cemu y conectar.",
		"client.status_running":"Sesion multijugador activa en Cemu.",
		"client.footer_hint":   "[Tab: Pestana]  [↓/↑: Seleccionar campo]  [c/Enter: Conectar]  [s: Guardar]  [l: IP Local]  [t: IP ZeroTier]  [Esc: Volver]",

		// Gamemodes
		"gm.title":            "CONFIGURACION DE GAMEMODES Y SERVIDOR DEDICADO",
		"gm.srv_running":      ">> SERVIDOR ACTIVO: Los cambios se aplicaran reiniciando el servidor al guardar.",
		"gm.srv_stopped":      ">> SERVIDOR DETENIDO: Configure sus reglas y presione [s] para iniciar.",
		"gm.special_mode":     "Modo de Juego Especial",
		"gm.mode_coop":        "Cooperativo Estandar (Libre)",
		"gm.mode_hunter":      "Hunter vs Speedrunner",
		"gm.mode_deathswap":   "DeathSwap",
		"gm.header_coop":      "OPCIONES DE SINCRONIZACION (MODO LIBRE / PERSONALIZABLE):",
		"gm.header_hunter":    "OPCIONES DE SINCRONIZACION (PROGRESO INDIVIDUAL - HUNTER VS SPEEDRUNNER):",
		"gm.header_deathswap": "OPCIONES DE SINCRONIZACION (DESACTIVADAS POR SUPERVIVENCIA - DEATHSWAP):",
		"gm.sync_quest":       "Sincronizar Misiones (QuestSync)",
		"gm.sync_shrine":      "Sincronizar Santuarios (ShrineSync)",
		"gm.sync_tower":       "Sincronizar Torres (TowerSync)",
		"gm.sync_korok":       "Sincronizar Kologs (KorokSync)",
		"gm.sync_enemy":       "Sincronizar Enemigos (EnemySync)",
		"gm.sync_dungeon":     "Sincronizar Mazmorras (DungeonSync)",
		"gm.sync_location":    "Sincronizar Ubicaciones (LocationSync)",
		"gm.locked":           "[-] Bloqueado",
		"gm.server_params":    "PARAMETROS DEL SERVIDOR:",
		"gm.password":         "Contrasena",
		"gm.description":      "Descripcion",
		"gm.footer_hint":      "[Tab: Pestana]  [←/→: Cambiar modo]  [Espacio: Alternar opcion]  [s: Guardar y Aplicar]  [Esc: Volver]",

		// Paths
		"paths.title":        "GESTOR DE RUTAS DEL JUEGO, UPDATE Y DLC",
		"paths.base":         "1. Carpeta del Juego Base (Zelda BotW):",
		"paths.update":       "2. Carpeta de Actualizacion (Update v208):",
		"paths.dlc":          "3. Carpeta de DLC (v80):",
		"paths.status":       "Estado:",
		"paths.browse":       "[f: Explorar]",
		"paths.footer_hint":  "[Tab: Pestana]  [↓/↑: Cambiar campo]  [f: Abrir Explorador]  [Enter: Guardar y Aplicar]  [Esc: Volver]",
		"paths.val_not_set":  "Ruta no configurada",
		"paths.val_rpx":      "Valido (Ejecutable U-King.rpx detectado)",
		"paths.val_content":  "Valido (Carpeta content detectada)",
		"paths.val_dir":      "Valido (Carpeta del juego detectada)",
		"paths.val_err_base": "No se encontro code/U-King.rpx ni carpeta content",
		"paths.val_up_cemu":  "Instalado en Cemu (v208)",
		"paths.val_up_ext":   "Vinculado externamente",
		"paths.val_up_err":   "No encontrado (Requiere Update v208)",
		"paths.val_dlc_cemu": "Instalado en Cemu (v80)",
		"paths.val_dlc_ext":  "Vinculado externamente",
		"paths.val_dlc_err":  "No encontrado (Requiere DLC v80)",

		// Network
		"net.title":        "DETALLES DE RED & MULTIJUGADOR",
		"net.interfaces":   "INTERFACES DE RED DETECTADAS:",
		"net.zerotier":     "ZeroTier VPN (Principal):",
		"net.tailscale":    "Tailscale VPN (Secundaria):",
		"net.local":        "Red Local (LAN):",
		"net.inactive":     "Inactivo / No Conectado",
		"net.instructions": "INSTRUCCIONES DE CONEXION MULTIJUGADOR:",
		"net.host_info":    "- Anfitrion (Servidor): Conecta localmente a %s en el puerto %s",
		"net.client_info":  "- Clientes remotos: Ingresan tu IP de ZeroTier / VPN (%s) en la pestana Cliente",
		"net.port_info":    "- Puerto por defecto: %s (TCP/UDP)",
		"net.zt_tip":       "- Recomendado: Usar la red compartida de ZeroTier para salas virtuales entre amigos.",
		"net.footer_hint":  "[t: Copiar IP de ZeroTier/VPN]  [Tab: Siguiente pestana]",

		// Settings & About
		"set.title":          "AJUSTES",
		"set.pref_title":     "PREFERENCIAS DE LA APLICACION:",
		"set.lang_label":     "Idioma / Language:",
		"set.about_title":    "ACERCA DE ESTE PROYECTO:",
		"set.app_name":      "Zelda BotW Multiplayer Manager para Linux / SteamOS",
		"set.version":        "Version: v1.0.0",
		"set.author":         "Autor: CarlosEvCode",
		"set.github":         "GitHub: https://github.com/CarlosEvCode/botw-multiplayer-linux",
		"set.stack":          "Creado con: Go + Bubble Tea + Wine + Cemu 1.26.2 + Milk Bar Launcher + BCML",
		"set.desc":           "Gestor automatizado completo para hostear, conectar y jugar Zelda BotW Multiplayer en Linux.",
		"set.footer_hint":    "[Tab: Pestana]  [←/→ / Espacio: Cambiar Idioma]  [Esc: Volver]",

		// Footer
		"foot.server":    "Servidor",
		"foot.client":    "Conectar",
		"foot.cemu":      "Cemu",
		"foot.gamemodes": "Gamemodes",
		"foot.paths":     "Rutas",
		"foot.copy_ip":   "Copiar IP",
		"foot.tab":       "Pestana",
		"foot.quit":      "Salir",

		// Notifications
		"notif.srv_stopping":     "Deteniendo servidor dedicado...",
		"notif.srv_started":      "Servidor iniciado en puerto %s",
		"notif.srv_error":        "Error iniciando servidor: %v",
		"notif.client_starting":  "Lanzando Cemu y conectando a %s:%s...",
		"notif.client_stopping":  "Desconectando cliente multijugador...",
		"notif.client_saved":     "Parametros de conexion guardados!",
		"notif.client_set_local": "IP fijada a 127.0.0.1 (Anfitrion Local)",
		"notif.client_set_vpn":   "IP fijada a la IP de ZeroTier/VPN detectada: %s",
		"notif.cemu_starting":    "Lanzando Cemu 1.26.2...",
		"notif.ip_copied":        "IP de Tailscale copiada: %s",
		"notif.zt_ip_copied":     "IP de ZeroTier copiada: %s",
		"notif.local_ip_copied":  "IP Local copiada: %s",
		"notif.picker_opening":   "Abriendo explorador de carpetas...",
		"notif.paths_saved":      "Rutas actualizadas y enlazadas correctamente!",
		"notif.cfg_saved_reload": "Configuracion guardada! Reiniciando servidor con las nuevas reglas...",
		"notif.cfg_saved":        "Configuracion guardada! Lista para el inicio del servidor.",
		"notif.mode_coop":        "Modo: Cooperativo Libre - Sincronizaciones desbloqueadas",
		"notif.mode_hunter":      "Modo: Hunter vs Speedrunner - Progreso individual activado",
		"notif.mode_deathswap":   "Modo: DeathSwap - Supervivencia individual activada",
		"notif.sync_locked":      "Sincronizacion bloqueada en modo competitivo",
		"notif.lang_changed":     "Idioma cambiado a Español!",
	},
}
