package settings

import (
	"fmt"
	"os"

	"dario.cat/mergo"
	"github.com/BurntSushi/toml"
	"github.com/creasty/defaults"
	msettings "github.com/mikros-dev/protoc-gen-mikros-extensions/pkg/settings"
)

type Settings struct {
	General *General `toml:"general" default:"{}"`
	Mikros  *Mikros  `toml:"mikros" default:"{}"`
}

type Mikros struct {
	UseSuffixes      bool   `toml:"use_suffixes"`
	SettingsFilename string `toml:"settings_filename"`
}

type General struct {
	Debug                 bool   `toml:"debug" default:"false"`
	RemoveEnumPrefix      bool   `toml:"remove_enum_prefix" default:"true"`
	RemoveUnspecifiedEnum bool   `toml:"remove_unspecified_enum" default:"true"`
	OutputPath            string `toml:"output_path" default:"openapi"`
}

func LoadSettings(filename string) (*Settings, error) {
	var settings Settings

	if filename != "" {
		file, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		if err := toml.Unmarshal(file, &settings); err != nil {
			return nil, err
		}
	}

	defaultSettings, err := loadDefaultSettings()
	if err != nil {
		return nil, err
	}

	if err := mergo.Merge(&settings, defaultSettings); err != nil {
		return nil, err
	}

	return &settings, nil
}

func loadDefaultSettings() (*Settings, error) {
	s := &Settings{}
	if err := defaults.Set(s); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Settings) MikrosSettings() (*msettings.Settings, error) {
	cfg, err := msettings.LoadSettings(s.Mikros.SettingsFilename)
	if err != nil {
		return nil, fmt.Errorf("could not load mikros plugin settings file: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid Settings: %w", err)
	}

	return cfg, nil
}