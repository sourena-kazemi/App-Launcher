package components

import (
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type extendedEntry struct {
	widget.Entry
	window fyne.Window
	list   *widget.List
}

func (e *extendedEntry) TypedKey(keyEvent *fyne.KeyEvent) {
	if keyEvent.Name == fyne.KeyEscape {
		e.window.Close()
		return
	}
	if keyEvent.Name == fyne.KeyDown {
		e.window.Canvas().Focus(e.list)
	}
	e.Entry.TypedKey(keyEvent)
}

// func (e *extendedEntry) TypedShortcut(shortcut fyne.Shortcut) {
// 	if _, ok := shortcut.(*desktop.CustomShortcut); !ok {
// 		fmt.Println("Shortcut typed:", shortcut.ShortcutName())
// 		e.Entry.TypedShortcut(shortcut)
// 		return
// 	}
// }

func NewExtendedEntry(window fyne.Window, list *widget.List) *extendedEntry {
	entry := &extendedEntry{window: window, list: list}
	entry.ExtendBaseWidget(entry)
	return entry
}

type soundCommandEntry struct {
	widget.Entry
	window      fyne.Window
	upCommand   string
	downCommand string
}

func (e *soundCommandEntry) TypedKey(keyEvent *fyne.KeyEvent) {
	if keyEvent.Name == fyne.KeyEscape {
		e.window.Close()
		return
	}
	if keyEvent.Name == fyne.KeyDown {
		exec.Command("sh", "-c", e.downCommand).Run()
		result, err := exec.Command("bash", "-c", "awk -F'[][]' '/Left:/ { print $2 }' <(amixer sget Master)").Output()
		if err == nil {
			e.Entry.SetPlaceHolder(strings.TrimSuffix(strings.TrimSpace(string(result)), "%"))
		}
	}
	if keyEvent.Name == fyne.KeyUp {
		exec.Command("sh", "-c", e.upCommand).Run()
		result, err := exec.Command("bash", "-c", "awk -F'[][]' '/Left:/ { print $2 }' <(amixer sget Master)").Output()
		if err == nil {
			e.Entry.SetPlaceHolder(strings.TrimSuffix(strings.TrimSpace(string(result)), "%"))
		}
	}
	e.Entry.TypedKey(keyEvent)
}

func NewSoundCommandEntry(window fyne.Window, upCommand string, downCommand string) *soundCommandEntry {
	entry := &soundCommandEntry{window: window, upCommand: upCommand, downCommand: downCommand}
	entry.ExtendBaseWidget(entry)
	return entry
}
