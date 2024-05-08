package filesystem

func (c *FsController) Active() *FsController {
	if c.active {
		return c
	}

	if c.alt.active {
		return c.alt
	}

	return nil
}
