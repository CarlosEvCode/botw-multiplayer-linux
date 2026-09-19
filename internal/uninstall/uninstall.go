package uninstall

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RunUninstaller performs a complete uninstallation of Zelda BotW Multiplayer suite.
func RunUninstaller(autoConfirm bool) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not determine user home directory: %w", err)
	}

	prefixDir := filepath.Join(home, ".local", "share", "wineprefixes", "botw-multiplayer")
	launcherDir := filepath.Join(home, "Zelda_BotW_Multiplayer")
	localBin := filepath.Join(home, ".local", "bin", "botw-manager")
	desktopEntry := filepath.Join(home, ".local", "share", "applications", "zelda-botw-manager.desktop")
	configDir := filepath.Join(home, ".config", "botw-manager")
	tempDir := filepath.Join(os.TempDir(), "botw_mp_setup")

	fmt.Println("========================================================================")
	fmt.Println(" Zelda: Breath of the Wild Multiplayer - Uninstaller")
	fmt.Println("========================================================================")
	fmt.Println("")
	fmt.Println("This action will permanently remove the following components:")
	fmt.Printf(" [1] Wineprefix directory:    %s\n", prefixDir)
	fmt.Printf(" [2] Launchers directory:     %s\n", launcherDir)
	fmt.Printf(" [3] Application binary:      %s\n", localBin)
	fmt.Printf(" [4] Desktop Menu shortcut:   %s\n", desktopEntry)
	fmt.Printf(" [5] Manager configuration:   %s\n", configDir)
	fmt.Printf(" [6] Temporary setup files:   %s\n", tempDir)
	fmt.Println("")
	fmt.Println("Note: Your original Base Game, Update, and DLC dumps will NOT be deleted.")
	fmt.Println("========================================================================")

	if !autoConfirm {
		fmt.Print("Are you sure you want to proceed with uninstallation? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		input = strings.TrimSpace(strings.ToLower(input))
		if input != "y" && input != "yes" && input != "s" && input != "si" && input != "sí" {
			fmt.Println("Uninstallation aborted by user.")
			return nil
		}
	}

	fmt.Println("")
	fmt.Println("[*] Stopping any active Wine processes for botw-multiplayer...")
	// Terminate wine server for this prefix if running
	_ = exec.Command("wineserver", "-k").Run()

	removeTarget("Wineprefix", prefixDir)
	removeTarget("Launchers Directory", launcherDir)
	removeTarget("Desktop Shortcut", desktopEntry)
	removeTarget("Configuration Directory", configDir)
	removeTarget("Temporary Setup Directory", tempDir)
	removeTarget("Application Binary", localBin)

	fmt.Println("")
	fmt.Println("========================================================================")
	fmt.Println(" [SUCCESS] Zelda BotW Multiplayer uninstalled completely.")
	fmt.Println("========================================================================")
	return nil
}

func removeTarget(name, path string) {
	if _, err := os.Stat(path); err == nil || !os.IsNotExist(err) {
		fmt.Printf("[-] Removing %-24s -> %s\n", name+":", path)
		if err := os.RemoveAll(path); err != nil {
			fmt.Printf("    [WARN] Failed to remove %s: %v\n", path, err)
		}
	} else {
		fmt.Printf("[o] Skipped (not found): %-16s -> %s\n", name+":", path)
	}
}
