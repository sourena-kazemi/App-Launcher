package components

import (
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
