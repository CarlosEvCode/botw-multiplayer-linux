package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
}

func GetDefaultPrefix() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local/share/wineprefixes/botw-multiplayer")
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
	}

	// Read BCML settings.json if exists
	bcmlPath := filepath.Join(driveC, "users", user, "AppData/Local/bcml/settings.json")
	if data, err := os.ReadFile(bcmlPath); err == nil {
		var bcml BCMLSettings
		if err := json.Unmarshal(data, &bcml); err == nil {
			cfg.BCMLSetting = bcml
			cfg.BaseGame = toLinuxPath(driveC, bcml.GameDir)
			cfg.UpdatePath = toLinuxPath(driveC, bcml.UpdateDir)
			cfg.DLCPath = toLinuxPath(driveC, bcml.DlcDir)
		}
	}

	return cfg, nil
}

func (c *ManagerConfig) SavePaths(baseGame, updatePath, dlcPath string) error {
	c.BaseGame = strings.TrimSpace(baseGame)
	c.UpdatePath = strings.TrimSpace(updatePath)
	c.DLCPath = strings.TrimSpace(dlcPath)

	gameFolder := filepath.Base(c.BaseGame)
	if gameFolder == "" || gameFolder == "." {
		gameFolder = "The Legend of Zelda Breath of the Wild"
	}

	// Symlink base game in C:\Games\
	gamesDir := filepath.Join(c.DriveC, "Games")
	_ = os.MkdirAll(gamesDir, 0755)
	if c.BaseGame != "" {
		targetLink := filepath.Join(gamesDir, gameFolder)
		_ = os.Remove(targetLink)
		_ = os.Symlink(c.BaseGame, targetLink)
	}

	// Handle Update symlink if provided
	if c.UpdatePath != "" {
		updateDst := filepath.Join(c.CemuDir, "mlc01/usr/title/0005000e/101c9400")
		_ = os.MkdirAll(filepath.Dir(updateDst), 0755)
		_ = os.Remove(updateDst)
		_ = os.Symlink(c.UpdatePath, updateDst)
	}

	// Handle DLC symlink if provided
	if c.DLCPath != "" {
		dlcDst := filepath.Join(c.CemuDir, "mlc01/usr/title/0005000c/101c9400")
		_ = os.MkdirAll(filepath.Dir(dlcDst), 0755)
		_ = os.Remove(dlcDst)
		_ = os.Symlink(c.DLCPath, dlcDst)
	}

	// Prepare BCML Settings
	c.BCMLSetting = BCMLSettings{
		CemuDir:        "C:/cemu_1.26.2",
		GameDir:        fmt.Sprintf("C:/Games/%s/content", gameFolder),
		UpdateDir:      "C:/cemu_1.26.2/mlc01/usr/title/0005000e/101c9400/content",
		DlcDir:         "C:/cemu_1.26.2/mlc01/usr/title/0005000c/101c9400/content/0010",
		StoreDir:       fmt.Sprintf("C:/users/%s/AppData/Local/bcml", c.User),
		ExportDir:      "C:/cemu_1.26.2/graphicPacks/BreathOfTheWild_BCML",
		Lang:           "USen",
		WiiU:           true,
		NoHardlinks:    true,
		Loaded:         true,
		Changelog:      true,
		AutoGB:         true,
		LastVersion:    "3.10.8",
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
		return false, "Ruta no configurada"
	}
	rpx := filepath.Join(c.BaseGame, "code/U-King.rpx")
	content := filepath.Join(c.BaseGame, "content")
	if fileExists(rpx) || dirExists(content) {
		return true, "Valido (Base Game)"
	}
	return false, "No se encontro code/U-King.rpx ni carpeta content"
}

func (c *ManagerConfig) ValidateUpdate() (bool, string) {
	mlcUpdate := filepath.Join(c.CemuDir, "mlc01/usr/title/0005000e/101c9400")
	if dirExists(mlcUpdate) {
		return true, "Instalado en Cemu (v208)"
	}
	if c.UpdatePath != "" && dirExists(c.UpdatePath) {
		return true, "Vinculado externamente"
	}
	return false, "No encontrado (Requiere Update v208)"
}

func (c *ManagerConfig) ValidateDLC() (bool, string) {
	mlcDlc := filepath.Join(c.CemuDir, "mlc01/usr/title/0005000c/101c9400")
	if dirExists(mlcDlc) {
		return true, "Instalado en Cemu (v80)"
	}
	if c.DLCPath != "" && dirExists(c.DLCPath) {
		return true, "Vinculado externamente"
	}
	return false, "No encontrado (Requiere DLC v80)"
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
