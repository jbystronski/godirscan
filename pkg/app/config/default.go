package config

var defaultConfig = Config{
	CurrentSchema:        2,
	DefaultEditor:        "nano",
	DefaultRootDirectory: getUserDirectory(),
	Executors:            map[string]string{},
	MaxWorkers:           1500,
	ColorSchemas:         PrebuiltSchemas,
	Bookmarks:            map[string][]string{},
}
