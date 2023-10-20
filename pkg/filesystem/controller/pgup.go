package controller

func (c *Controller) pgUp() {
	if c.MovePgUp(c.DataAccessor.Len()) {
		c.render()
	}
}
