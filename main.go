package main

import (
	_ "embed"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
)

//go:embed assets/icon.ico
var trayIconICO []byte

//go:embed assets/icon.png
var trayIconPNG []byte

//go:embed assets/NotoSansSC-Regular.ttf
var fontData []byte

func main() {
	a := app.New()

	if desk, ok := a.(desktop.App); ok {
		iconRes := fyne.NewStaticResource("icon.png", trayIconPNG)
		a.SetIcon(iconRes)

		menu := fyne.NewMenu(
			"Prompt Shortcut",
			fyne.NewMenuItem("Prompt Templates", func() {
				OpenConfigWindow()
			}),
		)
		desk.SetSystemTrayMenu(menu)
	}

	hotkeyPressedChan := make(chan bool, 1)
	RegisterHotkey(hotkeyPressedChan)

	go func() {
		for pressed := range hotkeyPressedChan {
			if pressed {
				fmt.Println("hotkey pressed")

				fyne.Do(func() {
					inputWindow := CreateInputWindow(a, HandleUserInput)

					inputWindow.Show()
					inputWindow.RequestFocus()
				})
			}
		}
	}()

	a.Run()
}

func HandleUserInput(input string, canceled bool) {
	if canceled {
		fmt.Println("user canceled input")
		return
	}

	fmt.Printf("user input received: %s\n", input)
	go OpenBrowser(input)
}
