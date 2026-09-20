package assets

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed appdata/*
var appdataFS embed.FS

// EnsureAppdataFiles ensures that QuestFlagsNames.txt, ArmorMapping.txt, etc. exist in the target BOTWM roaming directory.
func EnsureAppdataFiles(botwmRoamingDir string) error {
	if err := os.MkdirAll(botwmRoamingDir, 0755); err != nil {
		return fmt.Errorf("failed to create BOTWM roaming directory: %w", err)
	}

	entries, err := fs.ReadDir(appdataFS, "appdata")
	if err != nil {
		return fmt.Errorf("failed to read embedded appdata: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		targetPath := filepath.Join(botwmRoamingDir, entry.Name())
		// If file doesn't exist or is empty, write from embedded FS
		if st, err := os.Stat(targetPath); err != nil || st.Size() == 0 {
			data, err := appdataFS.ReadFile("appdata/" + entry.Name())
			if err != nil {
				continue
			}
			_ = os.WriteFile(targetPath, data, 0644)
		}
	}

	return nil
}
