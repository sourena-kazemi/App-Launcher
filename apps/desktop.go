package apps

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type AppEntry struct {
	Name string
	Exec string
	Icon string
}

func contains(list []string, item string) bool {
	for _, v := range list {
		if strings.EqualFold(v, item) {
			return true
		}
	}
	return false
}

func parseDesktopFile(path string) (*AppEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// currentDE := os.Getenv("XDG_CURRENT_DESKTOP")
	currentDE := "GNOME"

	entry := &AppEntry{}
	scanner := bufio.NewScanner(file)
	inDesktopSection := false

	var typeIsApp bool
	var noDisplay, hidden bool
	var onlyShowIn, notShowIn []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "[Desktop Entry]" {
			inDesktopSection = true
			continue
		}
		if !inDesktopSection || strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "Name="):
			if entry.Name != "" && entry.Exec != "" && entry.Icon != "" {
				continue
			}
			entry.Name = strings.TrimPrefix(line, "Name=")
		case strings.HasPrefix(line, "Exec="):
			if entry.Name != "" && entry.Exec != "" && entry.Icon != "" {
				continue
			}
			entry.Exec = strings.TrimPrefix(line, "Exec=")
		case strings.HasPrefix(line, "Icon="):
			if entry.Name != "" && entry.Exec != "" && entry.Icon != "" {
				continue
			}
			entry.Icon = strings.TrimPrefix(line, "Icon=")
		case strings.HasPrefix(line, "Type="):
			typeIsApp = strings.TrimPrefix(line, "Type=") == "Application"
		case strings.HasPrefix(line, "NoDisplay="):
			noDisplay = strings.TrimPrefix(line, "NoDisplay=") == "true"
		case strings.HasPrefix(line, "Hidden="):
			hidden = strings.TrimPrefix(line, "Hidden=") == "true"
		case strings.HasPrefix(line, "OnlyShowIn="):
			raw := strings.TrimPrefix(line, "OnlyShowIn=")
			onlyShowIn = strings.Split(strings.TrimSuffix(raw, ";"), ";")
		case strings.HasPrefix(line, "NotShowIn="):
			raw := strings.TrimPrefix(line, "NotShowIn=")
			notShowIn = strings.Split(strings.TrimSuffix(raw, ";"), ";")
		}
	}
	if len(onlyShowIn) > 0 && !contains(onlyShowIn, currentDE) {
		return nil, fmt.Errorf("not shown in current DE")
	}
	if contains(notShowIn, currentDE) {
		return nil, fmt.Errorf("excluded from current DE")
	}
	if !typeIsApp || noDisplay || hidden || entry.Name == "" || entry.Exec == "" {
		return nil, fmt.Errorf("invalid or hidden entry")
	}

	entry.Icon = findIconPath(entry.Icon)
	return entry, nil
}

func FindDesktopEntries() (map[string]AppEntry, []string) {
	paths := []string{
		"/usr/share/applications",
		filepath.Join(os.Getenv("HOME"), ".local/share/applications"),
	}

	entries := make(map[string]AppEntry)
	names := []string{}

	for _, dir := range paths {
		files, _ := filepath.Glob(filepath.Join(dir, "*.desktop"))
		for _, file := range files {
			entry, err := parseDesktopFile(file)
			if err == nil {
				names = append(names, entry.Name)
				entries[entry.Name] = AppEntry{
					Name: entry.Name,
					Exec: entry.Exec,
					Icon: entry.Icon,
				}
			}
		}
	}
	return entries, names
}
