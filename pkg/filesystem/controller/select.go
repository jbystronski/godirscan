package controller

func (c *Controller) selectEntry() {
	if c.DataAccessor.Len() > 0 {
		if c.path == c.selected.BasePath() {
			c.selected.Clear()
		}

		en, ok := c.Find(c.Index())

		if !ok {
			return
		}

		c.selected.Toggle(en.FullPath())
		if c.Alt.path == c.path {
			c.Alt.fullRender()
		}
		c.down()

	}
}

func (c *Controller) selectAllEntries() {
	if c.DataAccessor.Len() > 0 {
		if c.path == c.selected.BasePath() {
			c.selected.Clear()
		}

		if c.selected.Len() > 0 {
			c.selected.Clear()
		} else {
			for _, k := range c.DataAccessor.All() {
				c.selected.Set(k.FullPath(), struct{}{})
			}
		}

		c.fullRender()
		if c.Alt.path == c.path {
			c.Alt.fullRender()
		}
	}
}
