package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

const configFolderName = ".promptshortcut"

const defaultPromptTemplate = `You are a professional native-level English translator and editor.
Translate the following sentence into English and provide different versions or natural alternatives. Match the tone of the user's original sentence. Do not include any introductory or concluding text.

Sentence to translate: "{{}}"`

var promptRe = regexp.MustCompile(`\{\{[^}]*\}\}`)

type Config struct {
	Index          int
	PromptTemplate string
}

type ConfigSettings struct {
	SelectedIndex int
	Configs       []Config
}

func getConfigJsonPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %w", err)
	}
	return filepath.Join(homeDir, configFolderName, "config.json"), nil
}

func ensureConfigDir(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}
	return nil
}

func ReadConfigSettings() (ConfigSettings, error) {
	path, err := getConfigJsonPath()
	if err != nil {
		return ConfigSettings{}, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		defaultSettings := ConfigSettings{
			SelectedIndex: 1,
			Configs: []Config{
				{Index: 1, PromptTemplate: defaultPromptTemplate},
			},
		}
		err = SaveConfigSettings(defaultSettings)
		return defaultSettings, err
	} else if err != nil {
		return ConfigSettings{}, fmt.Errorf("could not read config file: %w", err)
	}

	return parseConfig(path, data)
}

func parseConfig(path string, data []byte) (ConfigSettings, error) {
	var value json.RawMessage
	if err := json.Unmarshal(data, &value); err != nil {
		return ConfigSettings{}, fmt.Errorf("could not parse config json: %w", err)
	}

	var settings ConfigSettings

	if len(value) > 0 && value[0] == '[' {
		var legacyConfigs []struct {
			Config
			Selected bool
		}
		if err := json.Unmarshal(value, &legacyConfigs); err != nil {
			return ConfigSettings{}, fmt.Errorf("could not parse legacy config: %w", err)
		}

		for _, legacy := range legacyConfigs {
			settings.Configs = append(settings.Configs, legacy.Config)
			if legacy.Selected {
				settings.SelectedIndex = legacy.Index
			}
		}
		_ = SaveConfigSettings(settings)
		return settings, nil
	}

	if err := json.Unmarshal(value, &settings); err != nil {
		return ConfigSettings{}, fmt.Errorf("could not unmarshal config: %w", err)
	}

	if settings.SelectedIndex == 0 && len(settings.Configs) > 0 {
		settings.SelectedIndex = settings.Configs[0].Index
		_ = SaveConfigSettings(settings)
	}

	return settings, nil
}

func SaveConfigSettings(settings ConfigSettings) error {
	path, err := getConfigJsonPath()
	if err != nil {
		return err
	}

	if err := ensureConfigDir(path); err != nil {
		return err
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal config: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

func AddConfig(newConfig Config) error {
	settings, err := ReadConfigSettings()
	if err != nil {
		return err
	}

	for _, cfg := range settings.Configs {
		if cfg.Index == newConfig.Index {
			return fmt.Errorf("config with Index %d already exists", newConfig.Index)
		}
	}

	settings.Configs = append(settings.Configs, newConfig)
	if len(settings.Configs) == 1 {
		settings.SelectedIndex = newConfig.Index
	}

	return SaveConfigSettings(settings)
}

func DeleteConfig(indexToDelete int) error {
	settings, err := ReadConfigSettings()
	if err != nil {
		return err
	}

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
		return fmt.Errorf("config with Index %d not found", indexToDelete)
	}

	settings.Configs = updatedConfigs
	if settings.SelectedIndex == indexToDelete {
		settings.SelectedIndex = 0
		if len(updatedConfigs) > 0 {
			settings.SelectedIndex = updatedConfigs[0].Index
		}
	}

	return SaveConfigSettings(settings)
}

func SelectedPromptTemplate() (string, error) {
	settings, err := ReadConfigSettings()
	if err != nil {
		return "", err
	}
	for _, config := range settings.Configs {
		if config.Index == settings.SelectedIndex {
			return config.PromptTemplate, nil
		}
	}
	if len(settings.Configs) > 0 {
		return settings.Configs[0].PromptTemplate, nil
	}
	return defaultPromptTemplate, nil
}

func BuildFullPrompt(template string, input string) string {
	return promptRe.ReplaceAllStringFunc(template, func(string) string {
		return input
	})
}
