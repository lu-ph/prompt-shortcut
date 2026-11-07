package main

import (
	"image/color"

	"fyne.io/fyne/v2"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var customFontRes = fyne.NewStaticResource("NotoSansSC.ttf", fontData)

type InputWindowTheme struct{}

func (m InputWindowTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 0x22, G: 0x22, B: 0x22, A: 0xff}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 0x77, G: 0x77, B: 0x77, A: 0xff}
	case theme.ColorNameInputBackground, theme.ColorNameInputBorder, theme.ColorNameFocus, theme.ColorNamePrimary:
		return color.Transparent
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (m InputWindowTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}
func (m InputWindowTheme) Font(style fyne.TextStyle) fyne.Resource {
	if len(fontData) > 0 {
		return customFontRes
	}
	return theme.DefaultTheme().Font(style)
}
func (m InputWindowTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return 19
	}
	if name == theme.SizeNamePadding {
		return 9
	}
	return theme.DefaultTheme().Size(name)
}

type inputEntry struct {
	widget.Entry
	onEscape func()
	onSubmit func()
}

func newInputEntry() *inputEntry {
	e := &inputEntry{}
	e.ExtendBaseWidget(e)
	e.MultiLine = true
	e.Wrapping = fyne.TextWrapWord
	return e
}

func (e *inputEntry) TypedKey(k *fyne.KeyEvent) {
	switch k.Name {
	case fyne.KeyEscape:
		if e.onEscape != nil {
			e.onEscape()
		}
	case fyne.KeyReturn, fyne.KeyEnter:
		if e.onSubmit != nil {
			e.onSubmit()
		}
	default:
		e.Entry.TypedKey(k)
	}
}

func CreateInputWindow(app fyne.App, onResult func(text string, canceled bool)) fyne.Window {
	app.Settings().SetTheme(InputWindowTheme{})
	w := app.NewWindow("input")
	w.Resize(fyne.NewSize(550, 150))

	entry := newInputEntry()
	entry.SetPlaceHolder("Type something and press Enter...")

	canceled := false

	entry.onEscape = func() {
		canceled = true
		onResult("", true)
		w.Close()
	}

	entry.onSubmit = func() {
		onResult(entry.Text, false)
		w.Close()
	}

	w.SetCloseIntercept(func() {
		if !canceled && entry.Text != "" {
			onResult("", true)
		}
		w.Hide()
	})

	content := container.NewPadded(entry)
	w.SetContent(content)

	w.Canvas().Focus(entry)

	return w
}
