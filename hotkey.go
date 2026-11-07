package main

import (
	"fmt"
	"runtime"

	hook "github.com/robotn/gohook"
)

func RegisterHotkey(hotkeyPressedChan chan<- bool) {
	modifier := "ctrl"
	if runtime.GOOS == "darwin" {
		modifier = "cmd" // Use Command key on macOS
	}

	fmt.Printf("--- Listening z + %s + alt ---\n", modifier)

	hook.Register(hook.KeyDown, []string{"z", modifier, "alt"}, func(e hook.Event) {
		hotkeyPressedChan <- true
	})

	go func() {
		s := hook.Start()
		<-hook.Process(s)
	}()
}
