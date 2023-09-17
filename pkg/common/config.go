package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
)

type Config struct {
	CurrentSchema        uint    `json:"set_palette"`
	DefaultRootDirectory string  `json:"default_root"`
	DefaultEditor        string  `json:"default_editor"`
	MaxWorkers           int     `json:"max_concurrent_workers"`
	ColorSchemas         []Theme `json:"color_schemas"`
}

const configFileName = "godirscan.json"

var Cfg = &Config{}

var CurrentTheme = &Theme{}

var defaultConfig = Config{
	CurrentSchema: 0,

	DefaultEditor:        "nano",
	DefaultRootDirectory: getUserDirectory(),
	MaxWorkers:           1500,
	ColorSchemas: []Theme{
		{
			Main:        "magenta",
			Accent:      "bright_cyan",
			Highlight:   "black",
			BgHighlight: "bright_white",
			BgHeader:    "cyan",
			Header:      "black",
			Select:      "bright_white",
			BgSelect:    "magenta",
			Prompt:      "bright_white",
			BgPrompt:    "cyan",
		},
		{
			Main:        "blue",
			Accent:      "bright_yellow",
			BgHighlight: "bright_white",
			Highlight:   "black",
			BgHeader:    "bright_blue",
			Header:      "bright_white",
			Select:      "blue",
			BgSelect:    "bright_yellow",
			Prompt:      "bright_white",
			BgPrompt:    "yellow",
		},
		{
			Main:        "bright_black",
			Accent:      "bright_white",
			BgHighlight: "bright_white",
			Highlight:   "black",
			Header:      "black",
			BgHeader:    "bright_black",
			Select:      "bright_white",
			BgSelect:    "bright_black",
			Prompt:      "bright_white",
			BgPrompt:    "bright_black",
		},
		{
			Main:        "red",
			Accent:      "bright_yellow",
			BgSelect:    "bright_yellow",
			Select:      "black",
			BgHighlight: "bright_white",
			Highlight:   "black",
			Header:      "bright_yellow",
			BgHeader:    "bright_red",
			Prompt:      "red",
			BgPrompt:    "bright_yellow",
		},
		{
			Main:        "bright_cyan",
			Accent:      "white",
			BgSelect:    "bright_black",
			Select:      "bright_white",
			BgHighlight: "cyan",
			Highlight:   "white",
			Header:      "bright_white",
			BgHeader:    "cyan",
			Prompt:      "bright_white",
			BgPrompt:    "cyan",
		},
		{
			Main:        "yellow",
			Accent:      "white",
			BgSelect:    "yellow",
			Select:      "black",
			BgHighlight: "white",
			Highlight:   "black",
			Header:      "bright_white",
			BgHeader:    "bright_black",
			Prompt:      "bright_white",
			BgPrompt:    "magenta",
		},
	},
}

func getUserDirectory() string {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	return currentUser.HomeDir
}

func ParseColorSchema(num uint, currentTheme *Theme) {
	parse := func(schemaValues []string, themeValues []*string, srcValues map[string]string, defaultValue string) {
		if len(schemaValues) != len(themeValues) {
			panic("number of values of a schema must match the number of values of the theme")
		}

		for i := range schemaValues {
			value, ok := srcValues[schemaValues[i]]
			if ok {
				*themeValues[i] = value
			} else {
				*themeValues[i] = srcValues[defaultValue]
			}
		}
	}

	s := Cfg.ColorSchemas[num]

	parse([]string{s.Main, s.Accent, s.Highlight, s.Select, s.Prompt, s.Header}, []*string{&currentTheme.Main, &currentTheme.Accent, &currentTheme.Highlight, &currentTheme.Select, &currentTheme.Prompt, &currentTheme.Header}, ColorsMap, "white")

	parse([]string{s.BgHighlight, s.BgHeader, s.BgSelect, s.BgPrompt}, []*string{&currentTheme.BgHighlight, &currentTheme.BgHeader, &currentTheme.BgSelect, &currentTheme.BgPrompt}, BackgroundsMap, "white")
}

func findOrCreateConfigFile() *os.File {
	var cFile *os.File

	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	pathToConfig := filepath.Join(configDir, configFileName)

	_, err = os.Stat(pathToConfig)
	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(pathToConfig)
		if err != nil {
			panic(err)
		}

		cFile = populateDefaultConfig(f)

	} else {

		f, err := os.OpenFile(pathToConfig, os.O_RDWR, 0o644)
		if err != nil {
			panic(err)
		}
		cFile = f

	}

	return cFile
}

func encodeConfig(configFile *os.File, config *Config) *os.File {
	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "\t")
	err := encoder.Encode(&config)
	if err != nil {
		panic(err)
	}

	_, err = configFile.Seek(0, io.SeekStart)

	if err != nil {
		panic(err)
	}

	return configFile
}

func populateDefaultConfig(configFile *os.File) *os.File {
	return encodeConfig(configFile, &defaultConfig)
}

func ParseConfigFile(config *Config) {
	configFile := findOrCreateConfigFile()

	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	err := decoder.Decode(Cfg)
	if err != nil {
		if err == io.EOF {
			fmt.Println("Config file is empty")
		} else {
			fmt.Println("Error decoding config file:", err)
		}
	}
}

func UpdateConfigFile(config *Config) {
	cfg := findOrCreateConfigFile()
	err := cfg.Truncate(0)
	if err != nil {
		panic(err)
	}

	defer cfg.Close()

	encodeConfig(cfg, config)
}

func (c *Config) ChangeTheme() {
	num := c.CurrentSchema
	if num < uint(len(c.ColorSchemas)-1) {
		num++
	} else {
		num = 0
	}

	c.CurrentSchema = num
	ParseColorSchema(num, CurrentTheme)
	UpdateConfigFile(Cfg)
}
