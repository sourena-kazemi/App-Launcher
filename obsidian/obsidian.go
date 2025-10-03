package obsidian

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type ChecklistItem struct {
	Text      string
	Checked   bool
	LineIndex int
}

func LoadChecklistItems(filePath string) ([]ChecklistItem, error) {
	var items []ChecklistItem

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lineIndex := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "- [x]") {
			parts := strings.SplitN(line, "- [x]", 2)
			if len(parts) == 2 {
				task := strings.TrimSpace(parts[1])
				items = append(items, ChecklistItem{
					Text:      task,
					Checked:   true,
					LineIndex: lineIndex,
				})
			}
		}
		if strings.Contains(line, "- [ ]") {
			parts := strings.SplitN(line, "- [ ]", 2)
			if len(parts) == 2 {
				task := strings.TrimSpace(parts[1])
				items = append(items, ChecklistItem{
					Text:      task,
					Checked:   false,
					LineIndex: lineIndex,
				})
			}
		}

		lineIndex++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func AddLine(filePath string, newLine string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	lines = append(lines, newLine)

	err = os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		return err
	}

	return nil
}

func RemoveLine(filePath string, lineNumber int) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if lineNumber >= 0 && lineNumber < len(lines) {
		lines = append(lines[:lineNumber], lines[lineNumber+1:]...)
	} else {
		return fmt.Errorf("invalid line number")
	}

	err = os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		return err
	}

	return nil
}
