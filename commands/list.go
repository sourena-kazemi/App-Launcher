package commands

import (
	"fyne.io/fyne/v2"
	light "github.com/sourena-kazemi/App-Launcher/commands/Light"
	sound "github.com/sourena-kazemi/App-Launcher/commands/Sound"
	todo "github.com/sourena-kazemi/App-Launcher/commands/Todo"
)

var Commands = map[string]func(fyne.App){
	"Todo":         todo.Todo,
	"Adjust Sound": sound.Sound,
	"Adjust Light": light.Light,
}

const CommandIcon = "/usr/share/icons/Papirus/48x48/apps/terminal.svg"

func GetCommandsNames() []string {
	result := []string{}
	for command := range Commands {
		result = append(result, command)
	}
	return result
}
