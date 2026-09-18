package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/CarlosEvCode/botw-multiplayer-linux/internal/i18n"
)

type BCMLSettings struct {
	CemuDir        string `json:"cemu_dir"`
	GameDir        string `json:"game_dir"`
	GameDirNX      string `json:"game_dir_nx"`
	UpdateDir      string `json:"update_dir"`
	DlcDir         string `json:"dlc_dir"`
	DlcDirNX       string `json:"dlc_dir_nx"`
	StoreDir       string `json:"store_dir"`
	ExportDir      string `json:"export_dir"`
	ExportDirNX    string `json:"export_dir_nx"`
	LoadReverse    bool   `json:"load_reverse"`
	SiteMeta       string `json:"site_meta"`
	NoGuess        bool   `json:"no_guess"`
	Lang           string `json:"lang"`
	NoCemu         bool   `json:"no_cemu"`
	WiiU           bool   `json:"wiiu"`
	NoHardlinks    bool   `json:"no_hardlinks"`
	Force7z        bool   `json:"force_7z"`
	SuppressUpdate bool   `json:"suppress_update"`
	Loaded         bool   `json:"loaded"`
	NSFW           bool   `json:"nsfw"`
	Changelog      bool   `json:"changelog"`
	StripGfx       bool   `json:"strip_gfx"`
	AutoGB         bool   `json:"auto_gb"`
	ShowGB         bool   `json:"show_gb"`
	DarkTheme      bool   `json:"dark_theme"`
	LastVersion    string `json:"last_version"`
}

type ServerConfigData struct {
	IP           string
	Port         string
	Password     string
	Description  string
	EnemySync    bool
	QuestSync    bool
	KorokSync    bool
	TowerSync    bool
	ShrineSync   bool
	LocationSync bool
	DungeonSync  bool
	SpecialMode  int // 0 = Standard Co-op, 1 = Hunter vs Speedrunner, 2 = DeathSwap
}

type ManagerConfig struct {
	PrefixDir   string
	DriveC      string
	CemuDir     string
	MilkBarDir  string
	BaseGame    string
	UpdatePath  string
	DLCPath     string
	User        string
	BCMLSetting BCMLSettings
	ServerCfg   ServerConfigData
}

func GetDefaultPrefix() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local/share/wineprefixes/botw-multiplayer")
}

func NormalizeGameDir(p string) string {
	cleaned := strings.TrimSpace(p)
	if cleaned == "" {
		return ""
	}
	cleaned = strings.TrimRight(cleaned, "/\\")
	baseName := filepath.Base(cleaned)
	if baseName == "content" || baseName == "code" || baseName == "meta" {
		cleaned = filepath.Dir(cleaned)
	}
	return cleaned
}

func PickDirectory(title, initialDir string) (string, error) {
	// Try zenity first
	if _, err := exec.LookPath("zenity"); err == nil {
		args := []string{"--file-selection", "--directory", "--title=" + title}
		if initialDir != "" && dirExists(initialDir) {
			args = append(args, "--filename="+filepath.Clean(initialDir)+"/")
		}
		cmd := exec.Command("zenity", args...)
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(out)), nil
	}

	// Try kdialog
	if _, err := exec.LookPath("kdialog"); err == nil {
		args := []string{"--getexistingdirectory", initialDir, "--title", title}
		cmd := exec.Command("kdialog", args...)
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(out)), nil
	}

	return "", fmt.Errorf("no GUI file picker found (zenity or kdialog required)")
}

func LoadConfig() (*ManagerConfig, error) {
	home, _ := os.UserHomeDir()
	user := os.Getenv("USER")
	if user == "" {
		user = filepath.Base(home)
	}

	prefix := GetDefaultPrefix()
	driveC := filepath.Join(prefix, "drive_c")
	cemuDir := filepath.Join(driveC, "cemu_1.26.2")
	milkBarDir := filepath.Join(driveC, "MilkBarLauncher")

	cfg := &ManagerConfig{
		PrefixDir:  prefix,
		DriveC:     driveC,
		CemuDir:    cemuDir,
		MilkBarDir: milkBarDir,
		User:       user,
		ServerCfg: ServerConfigData{
			IP:          "127.0.0.1",
			Port:        "5050",
			Description: "Explore Hyrule with Friends!",
			SpecialMode: 0,
		},
	}

	// Read BCML settings.json if exists
	bcmlPath := filepath.Join(driveC, "users", user, "AppData/Local/bcml/settings.json")
	if data, err := os.ReadFile(bcmlPath); err == nil {
		var bcml BCMLSettings
		if err := json.Unmarshal(data, &bcml); err == nil {
			cfg.BCMLSetting = bcml
			rawBase := NormalizeGameDir(toLinuxPath(driveC, bcml.GameDir))
			if realPath, err := filepath.EvalSymlinks(rawBase); err == nil {
				cfg.BaseGame = realPath
			} else {
				cfg.BaseGame = rawBase
			}

			rawUp := NormalizeGameDir(toLinuxPath(driveC, bcml.UpdateDir))
			if realUp, err := filepath.EvalSymlinks(rawUp); err == nil {
				cfg.UpdatePath = realUp
			} else {
				cfg.UpdatePath = rawUp
			}

			rawDlc := NormalizeGameDir(toLinuxPath(driveC, bcml.DlcDir))
			if realDlc, err := filepath.EvalSymlinks(rawDlc); err == nil {
				cfg.DLCPath = realDlc
			} else {
				cfg.DLCPath = rawDlc
			}
		}
	}

	// Read ServerConfig.ini
	cfg.LoadServerConfig()

	return cfg, nil
}

func (c *ManagerConfig) LoadServerConfig() {
	iniPath := filepath.Join(c.MilkBarDir, "DedicatedServer/ServerConfig.ini")
	f, err := os.Open(iniPath)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") || line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])

		switch k {
		case "IP":
			c.ServerCfg.IP = v
		case "Port":
			c.ServerCfg.Port = v
		case "Password":
			c.ServerCfg.Password = v
		case "Description":
			c.ServerCfg.Description = v
		case "EnemySync":
			c.ServerCfg.EnemySync = parseBool(v)
		case "QuestSync":
			c.ServerCfg.QuestSync = parseBool(v)
		case "KorokSync":
			c.ServerCfg.KorokSync = parseBool(v)
		case "TowerSync":
			c.ServerCfg.TowerSync = parseBool(v)
		case "ShrineSync":
			c.ServerCfg.ShrineSync = parseBool(v)
		case "LocationSync":
			c.ServerCfg.LocationSync = parseBool(v)
		case "DungeonSync":
			c.ServerCfg.DungeonSync = parseBool(v)
		case "Special":
			if val, err := strconv.Atoi(v); err == nil {
				c.ServerCfg.SpecialMode = val
			}
		}
	}
}

func (c *ManagerConfig) SaveServerConfig(sc ServerConfigData) error {
	c.ServerCfg = sc
	if c.ServerCfg.IP == "" {
		c.ServerCfg.IP = "127.0.0.1"
	}
	if c.ServerCfg.Port == "" {
		c.ServerCfg.Port = "5050"
	}

	iniPath := filepath.Join(c.MilkBarDir, "DedicatedServer/ServerConfig.ini")
	_ = os.MkdirAll(filepath.Dir(iniPath), 0755)

	content := fmt.Sprintf(`[Connection]
IP=%s
Port=%s
Password=%s

[ServerInformation]
Description=%s

[Gamemode]
DefaultGamemode=True

[DefaultGamemode]
Name=Custom
EnemySync=%t
QuestSync=%t
KorokSync=%t
TowerSync=%t
ShrineSync=%t
LocationSync=%t
DungeonSync=%t
# 0 for no special gamemode, 1 for HunterVsSpeedrunner, 2 for DeathSwap
Special=%d
`,
		c.ServerCfg.IP,
		c.ServerCfg.Port,
		c.ServerCfg.Password,
		c.ServerCfg.Description,
		c.ServerCfg.EnemySync,
		c.ServerCfg.QuestSync,
		c.ServerCfg.KorokSync,
		c.ServerCfg.TowerSync,
		c.ServerCfg.ShrineSync,
		c.ServerCfg.LocationSync,
		c.ServerCfg.DungeonSync,
		c.ServerCfg.SpecialMode,
	)

	return os.WriteFile(iniPath, []byte(content), 0644)
}

func (c *ManagerConfig) SavePaths(baseGame, updatePath, dlcPath string) error {
	c.BaseGame = NormalizeGameDir(baseGame)
	c.UpdatePath = NormalizeGameDir(updatePath)
	c.DLCPath = NormalizeGameDir(dlcPath)

	gameFolder := filepath.Base(c.BaseGame)
	if gameFolder == "" || gameFolder == "." {
		gameFolder = "The Legend of Zelda Breath of the Wild"
	}

	// Symlink base game in C:\Games\
	gamesDir := filepath.Join(c.DriveC, "Games")
	_ = os.MkdirAll(gamesDir, 0755)

	// Clean out stale symlinks in C:\Games\ to prevent ghost games in Cemu
	if entries, err := os.ReadDir(gamesDir); err == nil {
		for _, entry := range entries {
			fullEntry := filepath.Join(gamesDir, entry.Name())
			_ = os.Remove(fullEntry)
		}
	}

	validBase, _ := c.ValidateBaseGame()
	if c.BaseGame != "" && validBase {
		targetLink := filepath.Join(gamesDir, gameFolder)
		_ = os.Symlink(c.BaseGame, targetLink)
	}

	// Handle Update symlink if provided
	if c.UpdatePath != "" {
		updateDst := filepath.Join(c.CemuDir, "mlc01/usr/title/0005000e/101c9400")
		if filepath.Clean(c.UpdatePath) != filepath.Clean(updateDst) {
			_ = os.MkdirAll(filepath.Dir(updateDst), 0755)
			_ = os.Remove(updateDst)
			_ = os.Symlink(c.UpdatePath, updateDst)
		}
	}

	// Handle DLC symlink if provided
	if c.DLCPath != "" {
		dlcDst := filepath.Join(c.CemuDir, "mlc01/usr/title/0005000c/101c9400")
		if filepath.Clean(c.DLCPath) != filepath.Clean(dlcDst) {
			_ = os.MkdirAll(filepath.Dir(dlcDst), 0755)
			_ = os.Remove(dlcDst)
			_ = os.Symlink(c.DLCPath, dlcDst)
		}
	}

	// Prepare BCML Settings
	c.BCMLSetting = BCMLSettings{
		CemuDir:     "C:/cemu_1.26.2",
		GameDir:     fmt.Sprintf("C:/Games/%s/content", gameFolder),
		UpdateDir:   "C:/cemu_1.26.2/mlc01/usr/title/0005000e/101c9400/content",
		DlcDir:      "C:/cemu_1.26.2/mlc01/usr/title/0005000c/101c9400/content/0010",
		StoreDir:    fmt.Sprintf("C:/users/%s/AppData/Local/bcml", c.User),
		ExportDir:   "C:/cemu_1.26.2/graphicPacks/BreathOfTheWild_BCML",
		Lang:        "USen",
		WiiU:        true,
		NoHardlinks: true,
		Loaded:      true,
		Changelog:   true,
		AutoGB:      true,
		LastVersion: "3.10.8",
	}

	data, err := json.MarshalIndent(c.BCMLSetting, "", "  ")
	if err != nil {
		return err
	}

	// Write to all necessary BCML paths
	home, _ := os.UserHomeDir()
	p1 := filepath.Join(c.DriveC, "users", c.User, "AppData/Local/bcml/settings.json")
	p2 := filepath.Join(c.CemuDir, "bcml/settings.json")
	p3 := filepath.Join(home, ".config/bcml/settings.json")

	_ = os.MkdirAll(filepath.Dir(p1), 0755)
	_ = os.MkdirAll(filepath.Dir(p2), 0755)
	_ = os.MkdirAll(filepath.Dir(p3), 0755)

	_ = os.WriteFile(p1, data, 0644)
	_ = os.WriteFile(p2, data, 0644)
	_ = os.WriteFile(p3, data, 0644)

	// Ensure merged link
	mergedLink := filepath.Join(c.DriveC, "users", c.User, "AppData/Local/bcml/merged")
	_ = os.Remove(mergedLink)
	_ = os.Symlink(filepath.Join(c.CemuDir, "graphicPacks/BreathOfTheWild_BCML"), mergedLink)

	return nil
}

func (c *ManagerConfig) ValidateBaseGame() (bool, string) {
	if c.BaseGame == "" {
		return false, i18n.T("paths.val_not_set")
	}
	norm := NormalizeGameDir(c.BaseGame)
	rpx1 := filepath.Join(norm, "code/U-King.rpx")
	rpx2 := filepath.Join(c.BaseGame, "code/U-King.rpx")
	content1 := filepath.Join(norm, "content")
	content2 := filepath.Join(c.BaseGame, "content")
	code1 := filepath.Join(norm, "code")
	code2 := filepath.Join(c.BaseGame, "code")

	if fileExists(rpx1) || fileExists(rpx2) {
		return true, i18n.T("paths.val_rpx")
	}
	if (dirExists(content1) && dirExists(code1)) || (dirExists(content2) && dirExists(code2)) {
		return true, i18n.T("paths.val_content")
	}
	return false, i18n.T("paths.val_err_base")
}

func (c *ManagerConfig) ValidateUpdate() (bool, string) {
	mlcUpdate := filepath.Join(c.CemuDir, "mlc01/usr/title/0005000e/101c9400")
	if dirExists(mlcUpdate) {
		return true, i18n.T("paths.val_up_cemu")
	}
	if c.UpdatePath != "" && dirExists(c.UpdatePath) {
		return true, i18n.T("paths.val_up_ext")
	}
	return false, i18n.T("paths.val_up_err")
}

func (c *ManagerConfig) ValidateDLC() (bool, string) {
	mlcDlc := filepath.Join(c.CemuDir, "mlc01/usr/title/0005000c/101c9400")
	if dirExists(mlcDlc) {
		return true, i18n.T("paths.val_dlc_cemu")
	}
	if c.DLCPath != "" && dirExists(c.DLCPath) {
		return true, i18n.T("paths.val_dlc_ext")
	}
	return false, i18n.T("paths.val_dlc_err")
}

func toLinuxPath(driveC, winPath string) string {
	if winPath == "" {
		return ""
	}
	cleaned := strings.ReplaceAll(winPath, "\\", "/")
	if strings.HasPrefix(cleaned, "C:/") {
		return filepath.Join(driveC, cleaned[3:])
	}
	if strings.HasPrefix(cleaned, "Z:/") {
		return "/" + cleaned[3:]
	}
	return winPath
}

func fileExists(p string) bool {
	st, err := os.Stat(p)
	return err == nil && !st.IsDir()
}

func dirExists(p string) bool {
	st, err := os.Stat(p)
	return err == nil && st.IsDir()
}

func parseBool(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "true" || s == "1" || s == "yes"
}
