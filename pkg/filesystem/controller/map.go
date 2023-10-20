package controller

import k "github.com/eiannone/keyboard"

func (c *Controller) MapKey(action k.Key) {
	var activeController *Controller

	if c.Active() {
		activeController = c
	} else {
		activeController = c.Alt
	}

	switch action {
	case k.KeyF5:
		activeController.Refresh()

	case k.KeyCtrlSlash:
		activeController.updateTheme(activeController.config)

	case k.KeyTab:
		activeController.changeController()

	case k.KeyEsc:
		activeController.ctx.Cancel()

	case k.KeyHome:
		activeController.top()

	case k.KeyEnd:
		activeController.bottom()

	case k.KeyPgdn:
		activeController.pgDown()

	case k.KeyPgup:
		activeController.pgUp()

	case k.KeyEnter:
		activeController.execute()

	case k.KeyCtrlE:
		activeController.executeCmd()

	case k.KeyDelete, k.KeyF8:
		activeController.delete()

	case k.KeyCtrlR:
		activeController.rename()

	case k.KeyF7:
		activeController.newDirectory()

	case k.KeyCtrlW:
		activeController.newFile()

	case k.KeyCtrlF:
		activeController.search()

	case k.KeyCtrlS:
		activeController.sort()

	case k.KeyCtrlV:

		activeController.Copy()

	case k.KeyF6:
		activeController.Move()

	case k.KeyF9:

		activeController.scanDir()

	case k.KeyArrowLeft:
		activeController.left()

	case k.KeyArrowRight:
		activeController.right()

	case k.KeyInsert:
		activeController.selectEntry()

	case k.KeyCtrlA:
		activeController.selectAllEntries()

	case k.KeyCtrlG:
		activeController.goTo()

	case k.KeyArrowDown:
		activeController.down()

	case k.KeyArrowUp:
		activeController.up()

	case k.KeyCtrl4:
		activeController.edit(c.config.DefaultEditor)

	case k.KeyCtrl5:
		activeController.printHelp()

	}
}
