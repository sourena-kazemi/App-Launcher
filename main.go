package main

import (
	"fmt"
	"image/color"
	"math"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/mnogu/go-calculator"
	"github.com/sourena-kazemi/App-Launcher/apps"
)

type menu struct {
	CalculatorResult string
	AppEntries       []apps.AppEntry
}

func calculateHeight(itemCount int) float32 {
	inputHeight := float32(40)
	itemHeight := float32(40)

	if itemCount == 0 {
		return inputHeight
	}
	return inputHeight + (itemHeight * float32(math.Min(float64(itemCount), 10))) + (float32(math.Min(float64(itemCount), 10)+1) * 4)
}

func cleanExec(raw string) string {
	// Remove placeholders like %U, %f, etc.
	return strings.FieldsFunc(raw, func(r rune) bool {
		return r == '%' || r == '\n'
	})[0]
}

func main() {
	a := app.New()
	a.Settings().SetTheme(&myTheme{})

	if drv, ok := a.Driver().(desktop.Driver); ok {
		splash := drv.CreateSplashWindow()
		splash.SetFixedSize(true)
		splash.Resize(fyne.NewSize(800, 400))
		splash.Show()
		splash.Resize(fyne.NewSize(800, 40))

		entries, names := apps.FindDesktopEntries()
		selectedEntries := []apps.AppEntry{}
		menu := menu{}

		input := widget.NewEntry()
		input.SetPlaceHolder("Type to search")
		inputWrapper := container.New(layout.NewVBoxLayout(), input)

		list := widget.NewList(
			func() int {
				if menu.CalculatorResult != "" {
					return len(menu.AppEntries) + 1
				}
				return len(menu.AppEntries)
			},
			func() fyne.CanvasObject {
				leftPad := canvas.NewRectangle(color.Transparent)
				leftPad.SetMinSize(fyne.NewSize(12, 1))
				label := widget.NewLabel("")
				img := canvas.NewImageFromFile("")
				wrapper := container.New(layout.NewVBoxLayout(), layout.NewSpacer(), container.New(layout.NewHBoxLayout(), leftPad, img, label), layout.NewSpacer())
				return wrapper
			},
			func(i widget.ListItemID, o fyne.CanvasObject) {
				if menu.CalculatorResult != "" {
					if i == 0 {
						fmt.Println("hey, it's working?")
						o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[2].(*widget.Label).SetText(menu.CalculatorResult)
						img := canvas.NewImageFromResource(nil)
						o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1] = img
						fmt.Println("hey, it's still working?")
					} else {
						fmt.Println(entries[menu.AppEntries[i-1].Name].Name)
						o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[2].(*widget.Label).SetText(entries[menu.AppEntries[i-1].Name].Name)
						img := canvas.NewImageFromFile(entries[menu.AppEntries[i-1].Name].Icon)
						img.FillMode = canvas.ImageFillContain
						img.SetMinSize(fyne.NewSize(32, 32))
						o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1] = img
						fmt.Println(entries[menu.AppEntries[i-1].Name].Name)
					}
				} else {
					o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[2].(*widget.Label).SetText(entries[menu.AppEntries[i].Name].Name)
					img := canvas.NewImageFromFile(entries[menu.AppEntries[i].Name].Icon)
					img.FillMode = canvas.ImageFillContain
					img.SetMinSize(fyne.NewSize(32, 32))
					o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1] = img
				}
			},
		)
		list.OnSelected = func(i int) {
			cmd := cleanExec(entries[menu.AppEntries[i].Name].Exec)
			if cmd != "" {
				exec.Command("sh", "-c", cmd).Start()
			}
		}

		input.OnChanged = func(s string) {
			result, err := calculator.Calculate(s)
			if err == nil {
				menu.CalculatorResult = fmt.Sprint(result)
			} else {
				menu.CalculatorResult = ""
			}

			selectedNames := fuzzy.FindNormalizedFold(s, names)
			selectedEntries = []apps.AppEntry{}
			for i := 0; i < len(selectedNames); i++ {
				selectedEntries = append(selectedEntries, entries[selectedNames[i]])
			}
			menu.AppEntries = selectedEntries
			list.Refresh()
			itemsCount := len(menu.AppEntries)
			if menu.CalculatorResult != "" {
				itemsCount += 1
			}
			if itemsCount != 0 {
				for i := 0; i < itemsCount; i++ {
					list.SetItemHeight(i, 40)
				}
				content := container.NewBorder(inputWrapper, nil, nil, nil, list)
				splash.SetContent(content)

				windowHeight := calculateHeight(itemsCount)
				splash.Resize(fyne.NewSize(800, windowHeight))
			} else {
				inputWrapper = container.New(layout.NewVBoxLayout(), layout.NewSpacer(), input, layout.NewSpacer())
				input.Resize(fyne.NewSize(800, 40))
				splash.SetContent(inputWrapper)
				splash.Resize(fyne.NewSize(800, 40))
			}
		}

		input.Resize(fyne.NewSize(800, 40))
		splash.SetContent(inputWrapper)
		splash.Canvas().Focus(input)
	}

	a.Run()
}
