package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type Theme struct{}

func (Theme) Font(s fyne.TextStyle) fyne.Resource {
	return resourceInter24ptRegularTtf
}

func (Theme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return 20
	}
	return theme.DefaultTheme().Size(name)
}

func (Theme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (Theme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}
