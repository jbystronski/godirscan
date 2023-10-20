package controller

func (c *Controller) changeController() {
	c.SetActive(false)
	c.Alt.SetActive(true)
	c.fullRender()
	c.Alt.fullRender()
}
