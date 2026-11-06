package main

import (
	_ "embed"
	"log"

	"github.com/mxschmitt/playwright-go"
)

func OpenInputWindow() (string, bool) {
	runOption := &playwright.RunOptions{
		SkipInstallBrowsers: true,
	}
	err := playwright.Install(runOption)
	if err != nil {
		log.Printf("could not install playwright dependencies: %v", err)
		return "", false
	}
	pw, err := playwright.Run()
	if err != nil {
		log.Printf("could not start playwright: %v", err)
		return "", false
	}
	defer pw.Stop()

	option := playwright.BrowserTypeLaunchOptions{
		Channel:  playwright.String("chrome"),
		Headless: playwright.Bool(false),
		Args: []string{
			"--window-size=300,200",
		},
	}
	browser, err := pw.Chromium.Launch(option)
	if err != nil {
		log.Printf("could not launch browser: %v", err)
		return "", false
	}
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		Viewport: &playwright.Size{
			Width:  600,
			Height: 200,
		},
	})

	defer browser.Close()

	page, err := context.NewPage()
	if err != nil {
		log.Printf("could not create page: %v", err)
		return "", false
	}
	if err := page.SetContent(string(inputPage)); err != nil {
		log.Printf("could not load input page: %v", err)
		return "", false
	}

	if err := page.BringToFront(); err != nil {
		log.Printf("Could not bring page to front: %v", err)
	}

	result, err := page.Evaluate(`() => {
		return new Promise((resolve) => {
			const textarea = document.querySelector('textarea');
			textarea.focus();
			
			textarea.addEventListener('keydown', (e) => {
				if (e.key === 'Enter' && !e.shiftKey) {
					e.preventDefault();
					const userInput = textarea.value;
					resolve(userInput);
				}
			});
		});
	}`)
	if err != nil {
		if !page.IsClosed() {
			log.Printf("Error reading input from page: %v", err)
		}
		return "", false
	}
	text, ok := result.(string)
	if !ok {
		log.Printf("input page returned an unexpected result")
		return "", false
	}
	return text, true
}

func OpenBrowser(text string) {
	runOption := &playwright.RunOptions{
		SkipInstallBrowsers: true,
	}
	err := playwright.Install(runOption)
	if err != nil {
		log.Fatalf("could not install playwright dependencies: %v", err)
	}
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	option := playwright.BrowserTypeLaunchOptions{
		Channel:  playwright.String("chrome"),
		Headless: playwright.Bool(false),
	}
	browser, err := pw.Chromium.Launch(option)
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	if _, err = page.Goto("https://gemini.google.com/app"); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	InputGemini(text, page)
}

func InputGemini(text string, page playwright.Page) {
	text = BuildFullPrompt(SelectedPromptTemplate(), text)
	inputBox := page.Locator("div[contenteditable='true'][role='textbox']").First()

	// Wait for it to be visible and ready
	if err := inputBox.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateVisible,
	}); err != nil {
		log.Fatalf("could not find input box: %v", err)
	}

	if err := inputBox.Click(); err != nil {
		log.Fatalf("could not click input box: %v", err)
	}

	if err := inputBox.Fill(text); err != nil {
		log.Fatalf("could not type text into input box: %v", err)
	}

	if err := page.Keyboard().Press("Enter"); err != nil {
		log.Fatalf("could not press enter: %v", err)
	}

	select {}
}

func OpenConfigWindow() {
	runOption := &playwright.RunOptions{
		SkipInstallBrowsers: true,
	}
	if err := playwright.Install(runOption); err != nil {
		log.Fatalf("could not install playwright dependencies: %v", err)
	}

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	defer pw.Stop()

	option := playwright.BrowserTypeLaunchOptions{
		Channel:  playwright.String("chrome"),
		Headless: playwright.Bool(false),
	}

	browser, err := pw.Chromium.Launch(option)
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	// ---------- expose Go functions to the page ----------
	// Load all configs
	err = page.ExposeBinding("loadConfigs", func(source *playwright.BindingSource, args ...interface{}) interface{} {
		path := GetConfigJsonPath()
		CreateConfigJson()
		settings := readConfigSettings(path)
		configs := make([]map[string]interface{}, 0, len(settings.Configs))
		for _, config := range settings.Configs {
			configs = append(configs, map[string]interface{}{
				"Index":          config.Index,
				"PromptTemplate": config.PromptTemplate,
			})
		}
		return map[string]interface{}{
			"SelectedIndex": settings.SelectedIndex,
			"Configs":       configs,
		}
	})
	if err != nil {
		log.Fatalf("expose loadConfigs: %v", err)
	}

	// overwrite the whole list
	err = page.ExposeBinding("saveConfigs", func(source *playwright.BindingSource, args ...interface{}) interface{} {
		if len(args) == 0 {
			return map[string]interface{}{"ok": false, "error": "no data"}
		}

		// args[0] arrives as a map containing the selected index and config list.
		payload, ok := args[0].(map[string]interface{})
		if !ok {
			return map[string]interface{}{"ok": false, "error": "invalid data format"}
		}
		selectedIndex, ok := configNumber(payload["SelectedIndex"])
		if !ok || selectedIndex < 0 || selectedIndex != float64(int(selectedIndex)) {
			return map[string]interface{}{"ok": false, "error": "invalid selected template index"}
		}
		raw, ok := payload["Configs"].([]interface{})
		if !ok {
			return map[string]interface{}{"ok": false, "error": "invalid config list"}
		}

		var configs []Config
		seenIndexes := make(map[int]bool, len(raw))
		for _, item := range raw {
			m, ok := item.(map[string]interface{})
			if !ok {
				return map[string]interface{}{"ok": false, "error": "invalid config entry"}
			}
			index, ok := configNumber(m["Index"])
			if !ok || index < 1 || index != float64(int(index)) {
				return map[string]interface{}{"ok": false, "error": "each config needs a positive integer index"}
			}
			promptTemplate, ok := m["PromptTemplate"].(string)
			if !ok {
				return map[string]interface{}{"ok": false, "error": "each config needs a prompt template"}
			}
			configIndex := int(index)
			if seenIndexes[configIndex] {
				return map[string]interface{}{"ok": false, "error": "duplicate config indexes are not allowed"}
			}
			seenIndexes[configIndex] = true
			configs = append(configs, Config{
				Index:          configIndex,
				PromptTemplate: promptTemplate,
			})
		}
		if len(configs) > 0 && (selectedIndex == 0 || !seenIndexes[int(selectedIndex)]) {
			return map[string]interface{}{"ok": false, "error": "selected template does not exist"}
		}
		if len(configs) == 0 && selectedIndex != 0 {
			return map[string]interface{}{"ok": false, "error": "empty config must not select a template"}
		}

		path := GetConfigJsonPath()
		settings := ConfigSettings{SelectedIndex: int(selectedIndex), Configs: configs}
		if err := writeConfigSettings(path, settings); err != nil {
			return map[string]interface{}{"ok": false, "error": err.Error()}
		}
		return map[string]interface{}{"ok": true}
	})
	if err != nil {
		log.Fatalf("expose saveConfigs: %v", err)
	}

	// Close the window
	err = page.ExposeBinding("closeWindow", func(source *playwright.BindingSource, args ...interface{}) interface{} {
		go func() {
			if err := browser.Close(); err != nil {
				log.Printf("could not close config browser: %v", err)
			}
		}()
		return nil
	})
	if err != nil {
		log.Fatalf("expose closeWindow: %v", err)
	}

	// ---------- load the HTML ----------
	if err := page.SetContent(string(configPage)); err != nil {
		log.Fatalf("could not load config page: %v", err)
	}

	if err := page.BringToFront(); err != nil {
		log.Printf("could not bring config page to front: %v", err)
	}

	// Keep the window open until the user closes it
	if _, err := page.WaitForEvent("close", playwright.PageWaitForEventOptions{
		Timeout: playwright.Float(0),
	}); err != nil {
		log.Printf("config window closed with error: %v", err)
	}
}

func configNumber(value interface{}) (float64, bool) {
	switch number := value.(type) {
	case float64:
		return number, true
	case float32:
		return float64(number), true
	case int:
		return float64(number), true
	case int32:
		return float64(number), true
	case int64:
		return float64(number), true
	default:
		return 0, false
	}
}

//go:embed pages/config.html
var configPage []byte

//go:embed pages/input.html
var inputPage []byte
