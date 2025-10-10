package settings

var Settings = map[string]string{
	"Applications":     "gnome-control-center applications",
	"Background":       "gnome-control-center background",
	"Bluetooth":        "gnome-control-center bluetooth",
	"Color":            "gnome-control-center color",
	"Display":          "gnome-control-center display",
	"Keyboard":         "gnome-control-center keyboard",
	"Mouse":            "gnome-control-center mouse",
	"Multitasking":     "gnome-control-center multitasking",
	"Network":          "gnome-control-center network",
	"Wifi":             "gnome-control-center wifi",
	"Notifications":    "gnome-control-center notifications",
	"Online Accounts":  "gnome-control-center online-accounts",
	"Power":            "gnome-control-center power",
	"Printers":         "gnome-control-center printers",
	"Privacy":          "gnome-control-center privacy",
	"Search":           "gnome-control-center search",
	"Sharing":          "gnome-control-center sharing",
	"Sound":            "gnome-control-center sound",
	"System":           "gnome-control-center system",
	"Ubuntu":           "gnome-control-center ubuntu",
	"Universal Access": "gnome-control-center universal-access",
	"Wacom":            "gnome-control-center wacom",
	"Wellbeing":        "gnome-control-center wellbeing",
	"Wwan":             "gnome-control-center wwan",
}

const SettingIcon = "/usr/share/icons/Papirus/48x48/apps/org.gnome.Settings.svg"

func GetSettingsNames() []string {
	result := []string{}
	for setting := range Settings {
		result = append(result, setting)
	}
	return result
}
