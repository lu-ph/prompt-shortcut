
# prompt-shortcut

prompt-shortcut is a small background app that quickly and conveniently prompts the Gemini webpage from anywhere using your custom prompt template. Press Ctrl+Alt+X to open a quick-input window, wrap your text in your custom template, and instantly send it to the Google Gemini webpage in your browser.

## Usage

1. Launch the app. A new icon will appear in your system tray.
2. Press Ctrl+Alt+X to bring up the input window anywhere.
3. Type your text and wait for the browser to open.
4. Right-click the tray icon and select Prompt Templates to add, edit, or select your active template. Use `{{}}` in your template to mark the exact position where your typed input will be inserted.

## Build

**Windows:**
```bash
GOOS=windows go build -ldflags "-H=windowsgui" -o promptshortcut.exe .
```

**Linux:**
```bash
GOOS=linux go build -o promptshortcut .
```

**macOS:**
```bash
GOOS=darwin go build -o promptshortcut .
```
*To run as a background app on macOS (hiding the terminal and Dock icon), package it into an `.app` bundle:*
```bash
mkdir -p promptshortcut.app/Contents/MacOS
mv promptshortcut promptshortcut.app/Contents/MacOS/
cp -r pages promptshortcut.app/Contents/MacOS/
cat <<EOF> promptshortcut.app/Contents/Info.plist
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "[http://www.apple.com/DTDs/PropertyList-1.0.dtd](http://www.apple.com/DTDs/PropertyList-1.0.dtd)">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>promptshortcut</string>
    <key>CFBundleIdentifier</key>
    <string>com.local.promptshortcut</string>
    <key>LSUIElement</key>
    <true/>
</dict>
</plist>
EOF
```
