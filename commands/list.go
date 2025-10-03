package commands

import (
	"fyne.io/fyne/v2"
	todo "github.com/sourena-kazemi/App-Launcher/commands/Todo"
)

var Commands = map[string]func(fyne.App){
	"Todo": todo.Todo,
}

const CommandIcon = "/usr/share/icons/Papirus/48x48/apps/org.gnome.Settings.svg"

func GetCommandsNames() []string {
	result := []string{}
	for command := range Commands {
		result = append(result, command)
	}
	return result
}
