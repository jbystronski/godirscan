package filesystem

func (c *FsController) selectEntry() {
	if c.data.Len() > 0 {
		if c.root == c.selected.BasePath() {
			c.selected.Clear()
		}

		en, ok := c.activeEntry()

		if !ok {
			return
		}

		c.selected.Toggle(en.FullPath())
		if c.alt.root == c.root {
			c.alt.render()
		}
		if c.NextEntry() {
			c.render()
		}

	}
}

func (c *FsController) selectAll() {
	if c.data.Len() > 0 {
		if c.root == c.selected.BasePath() {
			c.selected.Clear()
		}

		if c.selected.Len() > 0 {
			c.selected.Clear()
		} else {
			for _, k := range c.data.All() {
				c.selected.Set(k.FullPath(), struct{}{})
			}
		}

		c.fullRender()
		if c.alt.root == c.root {
			c.alt.fullRender()
		}
	}
}
