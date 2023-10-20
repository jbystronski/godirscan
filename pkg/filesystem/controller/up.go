package controller

func (c *Controller) up() {
	if c.MoveUp(c.DataAccessor.Len()) {
		c.render()
	}
}
