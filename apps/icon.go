package apps

import (
	"os"
	"path/filepath"
	"strings"
)

var iconDirs = []string{
	"/usr/share/icons/Papirus/48x48/apps",
	"/usr/share/icons/Papirus/48x48/actions",
	"/usr/share/icons/hicolor/scalable/apps",
	"/usr/share/icons/hicolor/128x128/apps",
	"/usr/share/icons/hicolor/48x48/apps",
	"/usr/share/pixmaps",
}

func findIconPath(iconName string) string {
	if strings.HasPrefix(iconName, "/") && fileExists(iconName) {
		return iconName
	}
	for _, dir := range iconDirs {
		for _, ext := range []string{".png", ".svg", ".xpm"} {
			path := filepath.Join(dir, iconName+ext)
			if fileExists(path) {
				return path
			}
		}
	}
	return "" // fallback: no icon found
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
