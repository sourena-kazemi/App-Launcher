package todo

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/sourena-kazemi/App-Launcher/components"
	"github.com/sourena-kazemi/App-Launcher/obsidian"
	"github.com/sourena-kazemi/App-Launcher/util"
)

func Todo(app fyne.App) {
	if drv, ok := app.Driver().(desktop.Driver); ok {
		splash := drv.CreateSplashWindow()
		splash.SetFixedSize(true)
		splash.Resize(fyne.NewSize(800, 400))
		splash.Show()

		filePath := "/home/sourena/Documents/Obsidian/Personal/Todo.md"
		backupPath := "/home/sourena/Documents/Obsidian/Personal/Todo History.md"

		items, err := obsidian.LoadChecklistItems(filePath)
		if err != nil {
			errorLabel := widget.NewLabel(fmt.Sprintf("Error loading file: %v", err))
			splash.SetContent(container.NewCenter(errorLabel))
			return
		}

		list := widget.NewList(
			func() int {
				return len(items)
			},
			func() fyne.CanvasObject {
				leftPad := canvas.NewRectangle(color.Transparent)
				leftPad.SetMinSize(fyne.NewSize(12, 1))
				label := widget.NewLabel("")
				wrapper := container.New(layout.NewVBoxLayout(), layout.NewSpacer(), label, layout.NewSpacer())
				return wrapper
			},
			func(i widget.ListItemID, o fyne.CanvasObject) {
				o.(*fyne.Container).Objects[1].(*widget.Label).SetText(items[i].Text)
				o.(*fyne.Container).Objects[1].(*widget.Label).Refresh()
			},
		)

		list.OnSelected = func(id widget.ListItemID) {
			fmt.Print(id)
			err := obsidian.RemoveLine(filePath, items[id].LineIndex)
			if err != nil {
				log.Fatal(err)
			}
			currentTime := time.Now()
			formattedTime := currentTime.Format("2006-01-02 15:04:05")

			err = obsidian.AddLine(backupPath, fmt.Sprintf("- [x] %s %s", items[id].Text, formattedTime))
			if err != nil {
				log.Fatal(err)
			}
			items, err = obsidian.LoadChecklistItems(filePath)
			if err != nil {
				log.Fatal(err)
			}
			list.Refresh()

		}

		input := components.NewExtendedEntry(splash, list)
		input.SetPlaceHolder("Add Todo")
		inputWrapper := container.New(layout.NewVBoxLayout(), input)

		input.OnSubmitted = func(s string) {
			err := obsidian.AddLine(filePath, fmt.Sprintf("- [ ] %s", s))
			if err != nil {
				log.Fatal(err)
			}
			items, err = obsidian.LoadChecklistItems(filePath)
			if err != nil {
				log.Fatal(err)
			}

			list.Refresh()

			for i := 0; i < len(items); i++ {
				list.SetItemHeight(i, 40)
			}
			content := container.NewBorder(inputWrapper, nil, nil, nil, list)
			splash.SetContent(content)
			windowHeight := util.CalculateHeight(len(items))
			splash.Resize(fyne.NewSize(800, windowHeight))
		}

		inputWrapper = container.New(layout.NewVBoxLayout(), layout.NewSpacer(), input, layout.NewSpacer())

		if len(items) == 0 {
			input.Resize(fyne.NewSize(800, 40))
			splash.SetContent(inputWrapper)
			splash.Resize(fyne.NewSize(800, 40))
		} else {
			for i := 0; i < len(items); i++ {
				list.SetItemHeight(i, 40)
			}
			content := container.NewBorder(inputWrapper, nil, nil, nil, list)
			splash.SetContent(content)
			windowHeight := util.CalculateHeight(len(items))
			splash.Resize(fyne.NewSize(800, windowHeight))
		}

		splash.Canvas().Focus(input)
	}
}
