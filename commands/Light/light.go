package light

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"github.com/sourena-kazemi/App-Launcher/components"
)

func getBacklightDevices() ([]string, error) {
	result, err := exec.Command("brightnessctl", "-l").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run brightnessctl -l: %v", err)
	}

	var backlightDevices []string

	re := regexp.MustCompile(`Device '([^']+)' of class 'backlight':`)

	scanner := bufio.NewScanner(strings.NewReader(string(result)))
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			backlightDevices = append(backlightDevices, matches[1])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading command output: %v", err)
	}

	return backlightDevices, nil
}

func Light(app fyne.App) {
	if drv, ok := app.Driver().(desktop.Driver); ok {
		splash := drv.CreateSplashWindow()
		splash.SetFixedSize(true)
		splash.Resize(fyne.NewSize(800, 400))
		splash.Show()

		devices, err := getBacklightDevices()
		if err != nil {
			log.Fatal(err)
		}

		var placeHolder string

		for _, device := range devices {
			lightLevel, err := exec.Command("bash", "-c", fmt.Sprintf("brightnessctl get -d %s", device)).Output()
			if err != nil {
				break
			}

			lightLevelInt, err := strconv.Atoi(strings.TrimSuffix(string(lightLevel), "\n"))
			if err != nil {
				break
			}
			maxLightLevel, err := exec.Command("bash", "-c", fmt.Sprintf("brightnessctl max -d %s", device)).Output()
			if err != nil {
				break
			}
			maxLightLevelInt, err := strconv.Atoi(strings.TrimSuffix(string(maxLightLevel), "\n"))
			if err != nil {
				break
			}
			percentage := int(float32(lightLevelInt) / float32(maxLightLevelInt) * 100)
			placeHolder += fmt.Sprintf(" %s : %d ", device, percentage)
		}

		input := components.NewLightCommandEntry(splash, devices)
		input.SetPlaceHolder(placeHolder)

		input.OnSubmitted = func(s string) {
			light, err := strconv.Atoi(s)
			for _, device := range devices {
				if err == nil {
					exec.Command("bash", "-c", fmt.Sprintf("brightnessctl set %d%% -d %s", light, device)).Run()
				} else {
					exec.Command("bash", "-c", fmt.Sprintf("brightnessctl set %s -d %s", s, device)).Run()
				}
				lightLevel, err := exec.Command("bash", "-c", fmt.Sprintf("brightnessctl get -d %s", device)).Output()
				if err != nil {
					break
				}

				lightLevelInt, err := strconv.Atoi(strings.TrimSuffix(string(lightLevel), "\n"))
				if err != nil {
					break
				}
				maxLightLevel, err := exec.Command("bash", "-c", fmt.Sprintf("brightnessctl max -d %s", device)).Output()
				if err != nil {
					break
				}
				maxLightLevelInt, err := strconv.Atoi(strings.TrimSuffix(string(maxLightLevel), "\n"))
				if err != nil {
					break
				}
				percentage := int(float32(lightLevelInt) / float32(maxLightLevelInt) * 100)
				placeHolder += fmt.Sprintf(" %s : %d ", device, percentage)
			}

			input.SetPlaceHolder(placeHolder)
			input.SetText("")
		}

		inputWrapper := container.New(layout.NewVBoxLayout(), layout.NewSpacer(), input, layout.NewSpacer())

		input.Resize(fyne.NewSize(800, 40))
		splash.SetContent(inputWrapper)
		splash.Resize(fyne.NewSize(800, 40))

		splash.Canvas().Focus(input)
	}

}
