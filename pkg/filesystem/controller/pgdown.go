package controller

func (c *Controller) pgDown() {
	if c.MovePgDown(c.DataAccessor.Len()) {
		c.render()
	}
}
