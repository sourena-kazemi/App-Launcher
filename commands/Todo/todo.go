package todo

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func Todo(app fyne.App) {
	if drv, ok := app.Driver().(desktop.Driver); ok {
		splash := drv.CreateSplashWindow()
		splash.SetFixedSize(true)
		splash.Resize(fyne.NewSize(800, 400))
		splash.Show()
		splash.Resize(fyne.NewSize(800, 40))

		label := widget.NewLabel("Test")
		splash.SetContent(label)
	}

}
