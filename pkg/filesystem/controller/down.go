package controller

func (c *Controller) down() {
	if c.MoveDown(c.DataAccessor.Len()) {
		c.render()
	}
}
