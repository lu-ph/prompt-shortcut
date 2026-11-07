
# prompt-shortcut

prompt-shortcut is a small background app that quickly and conveniently prompts the Gemini webpage from anywhere using your custom prompt template. Press Ctrl+Alt+Z to open a quick-input window, wrap your text in your custom template, and instantly send it to the Google Gemini webpage in your browser.

## Usage

1. Launch the app. A new icon will appear in your system tray.
2. Press `Ctrl+Alt+Z` (or `Command+Alt+Z` on macos) to bring up the input window anywhere.
3. Type your text and wait for the browser to open.
4. Right-click the tray icon and select Prompt Templates to add, edit, or select your active template. Use `{{}}` in your template to mark the exact position where your typed input will be inserted.

## Build

You must have a C compiler installed on your machine to build the project.

### 1. Install Dependencies and Tools

Download the project dependencies
```bash
go mod download
```

Install the Fyne CLI packaging utility
```bash
go install fyne.io/tools/cmd/fyne@latest
```

### 2. Package the Application

**Windows:**
```bash
fyne package -icon assets/icon.png -name "Prompt Shortcut.exe"
```

**macOS:**
```bash
fyne package -icon assets/icon.png -name "Prompt Shortcut.app"
```

**Linux:**
```bash
fyne package -icon assets/icon.png -name "Prompt Shortcut"
```