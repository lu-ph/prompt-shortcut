package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

const CONFIG_FOLDER_NAME string = ".promptshortcut"

const DEFAULT_PROMPT_TEMPLATE = `You are a professional native-level English translator and editor.
Translate the following sentence into English and provide 3 different versions or natural alternatives. Match the tone of the user's original sentence. Do not include any introductory or concluding text.

Sentence to translate: "{{}}"`

type Config struct {
	Index          int
	PromptTemplate string
}

type ConfigSettings struct {
	SelectedIndex int
	Configs       []Config
}

// GetConfigJsonPath returns the full path to config.json inside CONFIG_FOLDER_NAME
func GetConfigJsonPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("could not get home directory: %v", err)
	}

	dirPath := filepath.Join(homeDir, CONFIG_FOLDER_NAME)

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		log.Fatalf("could not create config directory: %v", err)
	}

	return filepath.Join(dirPath, "config.json")
}

func CreateConfigJson() {
	path := GetConfigJsonPath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		defaultSettings := ConfigSettings{
			SelectedIndex: 1,
			Configs: []Config{
				{Index: 1, PromptTemplate: DEFAULT_PROMPT_TEMPLATE},
			},
		}

		data, err := json.MarshalIndent(defaultSettings, "", "  ")
		if err != nil {
			log.Fatalf("could not marshal default config: %v", err)
		}

		err = os.WriteFile(path, data, 0644)
		if err != nil {
			log.Fatalf("could not write config file: %v", err)
		}
		log.Println("Created default config.json successfully.")
	}
}

func Add(newConfig Config) {
	path := GetConfigJsonPath()
	CreateConfigJson()
	settings := readConfigSettings(path)

	for _, cfg := range settings.Configs {
		if cfg.Index == newConfig.Index {
			log.Printf("Error: Config with Index %d already exists. Use Edit() instead.", newConfig.Index)
			return
		}
	}

	settings.Configs = append(settings.Configs, newConfig)
	if len(settings.Configs) == 1 {
		settings.SelectedIndex = newConfig.Index
	}
	if err := writeConfigSettings(path, settings); err != nil {
		log.Fatalf("could not write config file: %v", err)
	}

	log.Printf("Successfully added new config with Index %d", newConfig.Index)
}

func Delete(indexToDelete int) {
	path := GetConfigJsonPath()
	CreateConfigJson()
	settings := readConfigSettings(path)

	var updatedConfigs []Config
	found := false
	for _, cfg := range settings.Configs {
		if cfg.Index == indexToDelete {
			found = true
			continue
		}
		updatedConfigs = append(updatedConfigs, cfg)
	}

	if !found {
		log.Printf("Config with Index %d not found.", indexToDelete)
		return
	}

	settings.Configs = updatedConfigs
	if settings.SelectedIndex == indexToDelete {
		settings.SelectedIndex = 0
		if len(updatedConfigs) > 0 {
			settings.SelectedIndex = updatedConfigs[0].Index
		}
	}
	if err := writeConfigSettings(path, settings); err != nil {
		log.Fatalf("could not write config file: %v", err)
	}

	log.Printf("Successfully deleted config with Index %d", indexToDelete)
}

func Read() []Config {
	path := GetConfigJsonPath()
	CreateConfigJson()
	return readConfigSettings(path).Configs
}

func readConfigSettings(path string) ConfigSettings {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("could not read config file: %v", err)
	}

	var value json.RawMessage
	if err := json.Unmarshal(data, &value); err != nil {
		log.Fatalf("could not parse config json: %v", err)
	}

	if len(value) > 0 && value[0] == '[' {
		var legacyConfigs []struct {
			Config
			Selected bool
		}
		if err := json.Unmarshal(value, &legacyConfigs); err != nil {
			log.Fatalf("could not parse legacy config json: %v", err)
		}
		settings := ConfigSettings{Configs: make([]Config, 0, len(legacyConfigs))}
		for _, legacy := range legacyConfigs {
			settings.Configs = append(settings.Configs, legacy.Config)
			if legacy.Selected {
				settings.SelectedIndex = legacy.Index
			}
		}
		if settings.SelectedIndex == 0 && len(settings.Configs) > 0 {
			settings.SelectedIndex = settings.Configs[0].Index
		}
		if err := writeConfigSettings(path, settings); err != nil {
			log.Fatalf("could not migrate config file: %v", err)
		}
		return settings
	}

	var settings ConfigSettings
	if err := json.Unmarshal(value, &settings); err != nil {
		log.Fatalf("could not parse config json: %v", err)
	}
	if settings.SelectedIndex == 0 && len(settings.Configs) > 0 {
		settings.SelectedIndex = settings.Configs[0].Index
		if err := writeConfigSettings(path, settings); err != nil {
			log.Fatalf("could not update config selection: %v", err)
		}
	}
	return settings
}

func writeConfigSettings(path string, settings ConfigSettings) error {
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func SelectedPromptTemplate() string {
	path := GetConfigJsonPath()
	CreateConfigJson()
	settings := readConfigSettings(path)
	for _, config := range settings.Configs {
		if config.Index == settings.SelectedIndex {
			return config.PromptTemplate
		}
	}
	if len(settings.Configs) > 0 {
		return settings.Configs[0].PromptTemplate
	}
	return ""
}

func BuildFullPrompt(template string, input string) string {
	re := regexp.MustCompile(`\{\{[^}]*\}\}`)

	return re.ReplaceAllStringFunc(template, func(string) string {
		return input
	})
}
