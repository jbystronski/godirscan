package controller

func (c *Controller) top() {
	if c.MoveToTop(c.DataAccessor.Len()) {
		c.render()
	}
}
