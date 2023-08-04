package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jbystronski/godirscan/pkg/terminal"
)

type Palette struct {
	Main        string `json:"color_1"`
	Accent      string `json:"color_2"`
	BgHighlight string `json:"bg_highlight"`
	BgPrompt    string `json:"bg_prompt"`
	BgGlobal    string `json:"bg_global"`
}

type Config struct {
	CurrentSchema        uint      `json:"set_palette"`
	DefaultRootDirectory string    `json:"default_root"`
	SilentMode           bool      `josn:"silent_mode"`
	DefaultEditor        string    `json:"default_editor"`
	MaxWorkers           int       `json:"max_concurrent_workers"`
	ColorSchemas         []Palette `json:"color_schemas"`
}

var Cfg = &Config{}

var theme Palette

var defaultConfig = Config{
	CurrentSchema:        0,
	SilentMode:           false,
	DefaultEditor:        "nano",
	DefaultRootDirectory: terminal.GetUserDirectory(),
	MaxWorkers:           1500,
	ColorSchemas: []Palette{{
		Main:        "magenta",
		Accent:      "bright_cyan",
		BgHighlight: "bright_white",
		BgPrompt:    "cyan",
		BgGlobal:    "",
	}, {
		Main:        "blue",
		Accent:      "yellow",
		BgHighlight: "bright_white",
		BgPrompt:    "yellow",
		BgGlobal:    "",
	}, {
		Main:        "bright_black",
		Accent:      "bright_white",
		BgHighlight: "bright_white",
		BgPrompt:    "bright_black",
		BgGlobal:    "",
	}, {
		Main:        "red",
		Accent:      "bright_yellow",
		BgHighlight: "white",
		BgPrompt:    "red",
		BgGlobal:    "bright_blue",
	}},
}

func ParseColorSchema(num uint) {
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

	parse([]string{s.Main, s.Accent}, []*string{&theme.Main, &theme.Accent}, terminal.ColorsMap, "white")

	parse([]string{s.BgHighlight, s.BgPrompt, s.BgGlobal}, []*string{&theme.BgHighlight, &theme.BgPrompt, &theme.BgGlobal}, terminal.BackgroundsMap, "white")
}

const (
	configFileName = "godirscan.json"
	// fmtDir           = terminal.BgYellow
	// fmtFile          = terminal.BrightMagenta
	// promptBackground = terminal.BgYellow
	// promptColor      = terminal.Black

	pSeparator = string(os.PathSeparator)
)

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
	err := decoder.Decode(&config)
	if err != nil {
		if err == io.EOF {
			fmt.Println("Config file is empty")
		} else {
			fmt.Println("Error decoding config file:", err)
		}
	}
}

func updateConfigFile(config *Config) {
	cfg := findOrCreateConfigFile()
	err := cfg.Truncate(0)
	if err != nil {
		panic(err)
	}

	defer cfg.Close()

	encodeConfig(cfg, config)
}

func switchTheme(num uint) {
	if num < uint(len(Cfg.ColorSchemas)) {

		Cfg.CurrentSchema = num

		ParseColorSchema(num)
		updateConfigFile(Cfg)
	}
}
