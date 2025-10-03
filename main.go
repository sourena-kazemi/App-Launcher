package main

import (
	"fmt"
	"image/color"
	"os/exec"

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
	"github.com/sourena-kazemi/App-Launcher/commands"
	"github.com/sourena-kazemi/App-Launcher/components"
	"github.com/sourena-kazemi/App-Launcher/theme"
	"github.com/sourena-kazemi/App-Launcher/util"
)

type menu struct {
	CalculatorResult string

	AppEntries map[string]apps.AppEntry
	AppNames   []string

	Commands      map[string]func(fyne.App)
	CommandsNames []string
	CommandIcon   string

	SelectedItems []string
}

func main() {
	a := app.New()
	a.Settings().SetTheme(&theme.Theme{})

	if drv, ok := a.Driver().(desktop.Driver); ok {
		splash := drv.CreateSplashWindow()
		splash.SetFixedSize(true)
		splash.Resize(fyne.NewSize(800, 400))
		splash.Show()
		splash.Resize(fyne.NewSize(800, 40))

		entries, names := apps.FindDesktopEntries()
		menu := menu{}
		menu.AppEntries = entries
		menu.AppNames = names
		menu.Commands = commands.Commands
		menu.CommandsNames = commands.GetCommandsNames()
		menu.CommandIcon = commands.CommandIcon

		list := widget.NewList(
			func() int {
				if menu.CalculatorResult != "" {
					return len(menu.SelectedItems) + 1
				}
				return len(menu.SelectedItems)
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
						o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*canvas.Image).File = ""
						o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*canvas.Image).Refresh()
					} else {
						if _, ok := menu.Commands[menu.SelectedItems[i-1]]; ok {
							o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*canvas.Image).File = menu.CommandIcon
						} else {
							o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*canvas.Image).File = menu.AppEntries[menu.SelectedItems[i-1]].Icon

						}

						o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[2].(*widget.Label).SetText(menu.SelectedItems[i-1])

						o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*canvas.Image).FillMode = canvas.ImageFillContain
						o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*canvas.Image).SetMinSize(fyne.NewSize(32, 32))

						o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[2].(*widget.Label).Refresh()
						o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*canvas.Image).Refresh()
					}
				} else {
					if _, ok := menu.Commands[menu.SelectedItems[i]]; ok {
						o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*canvas.Image).File = menu.CommandIcon
					} else {
						o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*canvas.Image).File = menu.AppEntries[menu.SelectedItems[i]].Icon
					}
					o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[2].(*widget.Label).SetText(menu.SelectedItems[i])

					o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*canvas.Image).FillMode = canvas.ImageFillContain
					o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*canvas.Image).SetMinSize(fyne.NewSize(32, 32))

					o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[2].(*widget.Label).Refresh()
					o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*canvas.Image).Refresh()
				}
			},
		)

		list.OnSelected = func(i widget.ListItemID) {
			var cmd string
			if menu.CalculatorResult == "" {
				if command, ok := menu.Commands[menu.SelectedItems[i]]; ok {
					command(a)
				} else {
					cmd = util.CleanExec(menu.AppEntries[menu.SelectedItems[i]].Exec)
				}
			} else {
				if i == 0 {
					cmd = fmt.Sprintf("echo %s | xclip -selection clipboard", menu.CalculatorResult)
				} else {
					if command, ok := menu.Commands[menu.SelectedItems[i-1]]; ok {
						command(a)
					} else {
						cmd = util.CleanExec(menu.AppEntries[menu.SelectedItems[i-1]].Exec)
					}
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

			selectedNames := fuzzy.FindNormalizedFold(s, append(menu.AppNames, menu.CommandsNames...))
			menu.SelectedItems = []string{}
			for i := 0; i < len(selectedNames); i++ {
				menu.SelectedItems = append(menu.SelectedItems, selectedNames[i])
			}
			list.Refresh()

			itemsCount := len(menu.SelectedItems)
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

				windowHeight := util.CalculateHeight(itemsCount)
				splash.Resize(fyne.NewSize(800, windowHeight))
			} else {
				inputWrapper = container.New(layout.NewVBoxLayout(), layout.NewSpacer(), input, layout.NewSpacer())
				input.Resize(fyne.NewSize(800, 40))
				splash.SetContent(inputWrapper)
				splash.Resize(fyne.NewSize(800, 40))
			}
		}

		input.OnSubmitted = func(s string) {
			if menu.CalculatorResult != "" || len(menu.SelectedItems) > 0 {
				list.Select(0)
			}
		}

		input.Resize(fyne.NewSize(800, 40))
		splash.SetContent(inputWrapper)
		splash.Canvas().Focus(input)
	}

	a.Run()
}
