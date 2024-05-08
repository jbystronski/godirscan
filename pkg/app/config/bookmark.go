package config

func (c *Config) AddBookmark(groupName, bookmark string) {
	if group, ok := c.Bookmarks[groupName]; ok {
		for _, entry := range group {
			if bookmark == entry {
				return
			}
		}
		c.Bookmarks[groupName] = append(c.Bookmarks[groupName], bookmark)
		c.encode(c)

	}
}

func (c *Config) RemoveBookmark(groupName, bookmark string) {
	if group, ok := c.Bookmarks[groupName]; ok {
		position := -1

		for index, entry := range group {
			if bookmark == entry {
				position = index
				break

			}
		}

		if position != -1 {
			c.Bookmarks[groupName] = append(c.Bookmarks[groupName][:position], c.Bookmarks[groupName][position+1:]...)
			c.encode(c)
		}

	}
}

func (c *Config) AddBookmarkGroup(name string) {
	if _, ok := c.Bookmarks[name]; !ok {
		c.Bookmarks[name] = []string{}
		c.encode(c)
	}
}

func (c *Config) RemoveBookmarkGroup(name string) {
	if _, ok := c.Bookmarks[name]; ok {
		delete(c.Bookmarks, name)
		c.encode(c)
	}
}
