package main

import (
	_ "embed"
	"fmt"
	"runtime"

	"fyne.io/systray"
)

//go:embed icon.ico
var trayIconICO []byte

//go:embed icon.png
var trayIconPNG []byte

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	fmt.Println("--- start ---")
	if runtime.GOOS == "windows" {
		systray.SetIcon(trayIconICO)
	} else {
		systray.SetTemplateIcon(trayIconPNG, trayIconPNG)
	}
	systray.SetTitle("Prompr Shortcut")
	systray.SetTooltip("Prompr Shortcut")

	mSettings := systray.AddMenuItem("Prompt Templates", "Open prompt configuration settings")
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	go func() {
		for {
			select {
			case <-mSettings.ClickedCh:
				fmt.Println("Menu clicked: Settings")
				go OpenConfigWindow()
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()

	hotkeyPressedChan := make(chan bool, 1)
	inputWindowClosedChan := make(chan struct{}, 1)

	RegisterHotkey(hotkeyPressedChan)
	inputWindowOpen := false
	for {
		select {
		case pressed := <-hotkeyPressedChan:
			if !pressed || inputWindowOpen {
				continue
			}
			fmt.Println("hotkey pressed")
			inputWindowOpen = true
			go func() {
				userInput, submitted := OpenInputWindow()
				if submitted {
					fmt.Println("received input: " + userInput)
					go OpenBrowser(userInput)
				}
				inputWindowClosedChan <- struct{}{}
			}()
		case <-inputWindowClosedChan:
			inputWindowOpen = false
		}
	}
}

func onExit() {
	systray.Quit()
}
