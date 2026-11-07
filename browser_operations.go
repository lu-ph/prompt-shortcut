package main

import (
	"log"

	"github.com/mxschmitt/playwright-go"
)

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
	promptTemplate, err := SelectedPromptTemplate()
	if err != nil {
		log.Fatalf("could not read the prompt template: %v", err)
	}

	text = BuildFullPrompt(promptTemplate, text)
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
