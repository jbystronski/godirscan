package controller

func (c *Controller) bottom() {
	if c.MoveToBottom(c.DataAccessor.Len()) {
		c.render()
	}
}
