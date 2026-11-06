package main

import (
	"fmt"

	hook "github.com/robotn/gohook"
)

func RegisterHotkey(hotkeyPressedChan chan<- bool) {
	fmt.Println("--- Listening x + ctrl + alt ---")

	hook.Register(hook.KeyDown, []string{"x", "ctrl", "alt"}, func(e hook.Event) {
		hotkeyPressedChan <- true
	})

	go func() {
		s := hook.Start()
		<-hook.Process(s)
	}()
}