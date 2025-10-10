package sound

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"github.com/sourena-kazemi/App-Launcher/components"
)

func Sound(app fyne.App) {
	if drv, ok := app.Driver().(desktop.Driver); ok {
		splash := drv.CreateSplashWindow()
		splash.SetFixedSize(true)
		splash.Resize(fyne.NewSize(800, 400))
		splash.Show()

		input := components.NewSoundCommandEntry(splash, "amixer sset Master 5%+", "amixer sset Master 5%-")
		result, err := exec.Command("bash", "-c", "awk -F'[][]' '/Left:/ { print $2 }' <(amixer sget Master)").Output()
		if err == nil {
			input.SetPlaceHolder(strings.TrimSuffix(strings.TrimSpace(string(result)), "%"))
		} else {
			input.SetPlaceHolder("50")
		}

		input.OnSubmitted = func(s string) {
			volume, err := strconv.Atoi(s)
			var cmd string
			if err == nil {
				cmd = fmt.Sprintf("amixer sset Master %d%%", volume)
			} else {
				cmd = fmt.Sprintf("amixer sset Master %s", s)

			}
			result, err := exec.Command("bash", "-c", "awk -F'[][]' '/Left:/ { print $2 }' <(amixer sget Master)").Output()
			if err == nil {
				input.SetPlaceHolder(strings.TrimSuffix(strings.TrimSpace(string(result)), "%"))
			}
			exec.Command("sh", "-c", cmd).Start()
			input.SetText("")
		}

		inputWrapper := container.New(layout.NewVBoxLayout(), layout.NewSpacer(), input, layout.NewSpacer())

		input.Resize(fyne.NewSize(800, 40))
		splash.SetContent(inputWrapper)
		splash.Resize(fyne.NewSize(800, 40))

		splash.Canvas().Focus(input)
	}

}
