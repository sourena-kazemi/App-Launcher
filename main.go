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
	"github.com/sourena-kazemi/App-Launcher/components"
)

type menu struct {
	CalculatorResult   string
	AppEntries         map[string]apps.AppEntry
	AppNames           []string
	SelectedAppEntries []apps.AppEntry
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

		// selectedEntries := []apps.AppEntry{}
		entries, names := apps.FindDesktopEntries()
		menu := menu{}
		menu.AppEntries = entries
		menu.AppNames = names

		list := widget.NewList(
			func() int {
				if menu.CalculatorResult != "" {
					return len(menu.SelectedAppEntries) + 1
				}
				return len(menu.SelectedAppEntries)
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
						o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[2].(*widget.Label).SetText(menu.CalculatorResult)
						o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[2].(*widget.Label).Refresh()
						o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*canvas.Image).Resource = nil
						o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*canvas.Image).Refresh()
					} else {
						o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[2].(*widget.Label).SetText(menu.AppEntries[menu.SelectedAppEntries[i-1].Name].Name)
						o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[2].(*widget.Label).Refresh()
						o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*canvas.Image).File = menu.AppEntries[menu.SelectedAppEntries[i-1].Name].Icon
						o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*canvas.Image).FillMode = canvas.ImageFillContain
						o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*canvas.Image).SetMinSize(fyne.NewSize(32, 32))
						o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*canvas.Image).Refresh()
					}
				} else {
					o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[2].(*widget.Label).SetText(menu.AppEntries[menu.SelectedAppEntries[i].Name].Name)
					o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[2].(*widget.Label).Refresh()
					o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*canvas.Image).File = menu.AppEntries[menu.SelectedAppEntries[i].Name].Icon
					o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*canvas.Image).FillMode = canvas.ImageFillContain
					o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*canvas.Image).SetMinSize(fyne.NewSize(32, 32))
					o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*canvas.Image).Refresh()
				}
			},
		)

		list.OnSelected = func(i int) {
			var cmd string
			if menu.CalculatorResult == "" {
				cmd = cleanExec(menu.AppEntries[menu.SelectedAppEntries[i].Name].Exec)
			} else {
				if i == 0 {
					cmd = fmt.Sprintf("echo %s | xclip -selection clipboard", menu.CalculatorResult)
				} else {
					cmd = cleanExec(menu.AppEntries[menu.SelectedAppEntries[i-1].Name].Exec)
				}
			}
			if cmd != "" {
				exec.Command("sh", "-c", cmd).Start()
			}
			splash.Close()
		}

		input := components.NewExtendedEntry(splash, list)
		input.SetPlaceHolder("Type to search")
		inputWrapper := container.New(layout.NewVBoxLayout(), input)

		input.OnChanged = func(s string) {
			result, err := calculator.Calculate(s)
			if err == nil {
				menu.CalculatorResult = fmt.Sprint(result)
			} else {
				menu.CalculatorResult = ""
			}

			selectedNames := fuzzy.FindNormalizedFold(s, menu.AppNames)
			menu.SelectedAppEntries = []apps.AppEntry{}
			for i := 0; i < len(selectedNames); i++ {
				menu.SelectedAppEntries = append(menu.SelectedAppEntries, menu.AppEntries[selectedNames[i]])
			}
			list.Refresh()

			itemsCount := len(menu.SelectedAppEntries)
			if menu.CalculatorResult != "" {
				itemsCount += 1
			}
			if itemsCount != 0 {
				for i := 0; i < itemsCount; i++ {
					list.SetItemHeight(i, 40)
				}
				// content := container.NewBorder(inputWrapper, nil, nil, nil, baseList)
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

		input.OnSubmitted = func(s string) {
			if menu.CalculatorResult != "" || len(menu.SelectedAppEntries) > 0 {
				list.Select(0)
			}
		}

		input.Resize(fyne.NewSize(800, 40))
		splash.SetContent(inputWrapper)
		splash.Canvas().Focus(input)
	}

	a.Run()
}
