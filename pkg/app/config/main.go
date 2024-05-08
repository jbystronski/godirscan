package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"sync"
)

type Config struct {
	CurrentSchema        uint                `json:"schema_set"`
	DefaultRootDirectory string              `json:"default_root"`
	DefaultEditor        string              `json:"default_editor"`
	MaxWorkers           int                 `json:"max_concurrent_workers"`
	ColorSchemas         []Schema            `json:"color_schemas"`
	Executors            map[string]string   `json:"executors"`
	Bookmarks            map[string][]string `json:"bookmarks"`
}

var runningConfig *Config

var CurrentTheme = &Theme{}

var once sync.Once

func Running() *Config {
	once.Do(func() {
		runningConfig = &Config{}
	})

	return runningConfig
}

func (c *Config) GetSettingsFilepath() string {
	return filepath.Join(getConfigDirName(), "godirscan.json")
}

func (c *Config) getFile() *os.File {
	path := c.GetSettingsFilepath()

	if _, err := os.Stat(path); err != nil && errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(path)
		if err != nil {
			panic(err)
		}

		c.encode(&defaultConfig)
	}

	file, err := os.OpenFile(path, os.O_RDWR, 0o644)
	if err != nil {
		panic(err)
	}

	return file
}

func getConfigDirName() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	return dir
}

func getUserDirectory() string {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	return currentUser.HomeDir
}

func (c *Config) Parse() {
	file := c.getFile()

	defer file.Close()

	if err := json.NewDecoder(file).Decode(Running()); err != nil {
		panic(err)
	}
}

func (c *Config) encode(source *Config) {
	file := c.getFile()

	defer file.Close()

	if err := file.Truncate(0); err != nil {
		panic(err)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(source); err != nil {
		panic(err)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		panic(err)
	}
}
