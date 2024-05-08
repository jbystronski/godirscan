package filesystem

import (
	"errors"
	"os"

	"github.com/jbystronski/godirscan/pkg/app/boxes"
	"github.com/jbystronski/godirscan/pkg/app/config"
	"github.com/jbystronski/godirscan/pkg/app/data"
	"github.com/jbystronski/godirscan/pkg/app/menu"
	"github.com/jbystronski/godirscan/pkg/global"

	"github.com/jbystronski/godirscan/pkg/lib/maps"
	"github.com/jbystronski/godirscan/pkg/lib/pubsub"
	"github.com/jbystronski/godirscan/pkg/lib/termui"
)

type FsController struct {
	*pubsub.Node
	termui.Navigator

	root        string
	ctx         *global.CancelContext
	pool        *data.Pool
	cache       *data.CacheStore
	selected    *data.Selected
	data        data.Accessor
	active      bool
	backtrace   maps.MapAccessor[string, int]
	alt         *FsController
	defaultSort func()
	header      termui.Section
	panel       termui.Section

	// ui Panel
}

func NewFsController(n *pubsub.Node) *FsController {
	pool := data.NewDataPool()
	selected := data.NewSelected()
	cache := data.NewCacheStore(104857600)

	c = &FsController{
		n,
		termui.Navigator{},
		"", &global.CancelContext{}, pool, cache, selected, data.NewFsData(), true, &maps.GenericMap[string, int]{}, nil, nil,
		termui.NewSection(),
		termui.NewSection(),
	}

	c.header.Top = 1
	c.header.SetHeight(1).SetWidth(cols() / 2)

	c.header.Left = 1
	c.panel.SetBorder()
	c.panel.Top = 2

	c.panel.SetPadding(0, 0, 1, 0).SetHeight(rows() - 2).SetWidth(cols() / 2)
	c.panel.Left = 1

	c.MinOffset = c.panel.OutputFirstLine()
	c.MaxOffset = c.panel.OutputLastLine() - 1

	c.alt = &FsController{
		n,
		termui.Navigator{},
		"", &global.CancelContext{}, pool, cache, selected, data.NewFsData(), false, &maps.GenericMap[string, int]{}, c, nil,
		termui.NewSection(),
		termui.NewSection(),
	}

	c.alt.header.Top = 1

	c.alt.header.Width = cols() / 2
	c.alt.header.Height = 1

	c.alt.header.Left = cols()/2 + 1

	c.alt.panel.SetBorder()
	c.alt.panel.Top = 2

	c.alt.panel.SetPadding(0, 0, 1, 0).SetHeight(rows() - 2).SetWidth(cols() / 2)
	c.alt.panel.Left = cols()/2 + 1

	c.alt.MinOffset = c.alt.panel.OutputFirstLine()
	c.alt.MaxOffset = c.alt.panel.OutputLastLine() - 1

	c.On(pubsub.INIT, func() {
		if c.root == "" {
			c.sendError(errors.New("empty root path"))
		}

		cls()

		c.UpdateSize()
		c.alt.UpdateSize()

		c.panel.Print(themeMain())
		c.alt.panel.Print(themeMain())

		err := c.setStore(c.root)
		if err != nil {
			c.sendError(err)
		}

		err = c.alt.setStore(c.root)
		if err != nil {
			c.alt.sendError(err)
		}
	})

	registerEvents(c)

	//	registerEvents(main.alt)

	return c
}

func registerEvents(c *FsController) {
	// c.On(pubsub.INIT, func() {
	// 	if c.root == "" {
	// 		c.sendError(errors.New("empty root path"))
	// 	}

	// 	cls()

	// 	c.ui.Print()
	// 	c.alt.ui.Print()

	// 	err := c.setStore(c.root)
	// 	if err != nil {
	// 		c.sendError(err)
	// 	}

	// 	err = c.alt.setStore(c.root)
	// 	if err != nil {
	// 		c.alt.sendError(err)
	// 	}
	// })

	c.On(pubsub.TAB, func() {
		main := c
		alt := c.alt
		c = alt
		c.alt = main

		// c, c.alt = c.alt, c
		c.active = true
		c.alt.active = false

		// temp := main
		// main = main.alt
		// temp.alt = main
		// main = main.alt
		c.render()
		c.alt.render()

		//	c.Active().changeActive()
	})

	c.Subscribe("add_bookmark_group", func(m pubsub.Message) {
		c.EnqueueMessage("add_bookmark_group", m)
	})

	c.Subscribe("bookmark", func(m pubsub.Message) {
		c.EnqueueMessage("bookmark", m)
	})

	c.Subscribe("open_bookmark", func(m pubsub.Message) {
		c.EnqueueMessage("open_bookmark", m)
	})

	c.Subscribe("remove_bookmark", func(m pubsub.Message) {
		c.EnqueueMessage("remove_bookmark", m)
	})

	c.Subscribe("bookmark_group", func(m pubsub.Message) {
		c.EnqueueMessage("bookmark_group", m)
	})

	c.On(pubsub.RENDER, func() {
		cls()
		hideCursor()
		c.UpdateSize()
		// c.panel.Print(themeMain())
		// c.alt.panel.Print(themeMain())
		c.fullRender()
		c.alt.fullRender()

		// c.restoreScreen()
		// c.alt.restoreScreen()
	})

	c.On(pubsub.CTRL_V, func() {
		answ := c.getInput("Copy selected into the current directory? :", "y")

		if answ != "y" || c.selected.Len() == 0 {
			return
		}
		err := c.Copy(false)
		if err != nil {
			c.sendError(err)
		}
	})

	c.On(pubsub.F6, func() {
		answ := c.getInput("Move selected into the current directory? :", "y")

		if answ != "y" || c.selected.Len() == 0 {
			return
		}
		err := c.Copy(true)
		if err != nil {
			c.sendError(err)
		}
	})

	c.On(pubsub.BOOKMARK_GROUP_LIST, func() {
		c.OpenMenu(menu.BookmarkGroupListMenu())
	})

	c.On(pubsub.BOOKMARK_ADD_GROUP, func() {
		c.addBookmarkGroup()

		c.OpenMenu(menu.BookmarkGroupListMenu())
	})

	c.On(pubsub.BOOKMARK_REMOVE_GROUP, func() {
		c.removeBookmarkGroup(string(c.DequeueMessage("bookmark_group")))

		c.OpenMenu(menu.BookmarkGroupListMenu())
	})

	c.On(pubsub.REMOVE_BOOKMARK, func() {
		bGroup := string(c.DequeueMessage("bookmark_group"))
		bName := string(c.DequeueMessage("remove_bookmark"))

		c.removeBookmark(bGroup, bName)
	})

	c.On(pubsub.BOOKMARK_OPEN, func() {
		c.openBookmark(string(c.DequeueMessage("bookmark")))
	})

	c.On(pubsub.BOOKMARK, func() {
		group := c.DequeueMessage("bookmark_group")

		c.bookmark(string(group))
	})

	c.On(pubsub.BOOKMARK_GROUP, func() {
		if c.HasNext() {
			c.Next.Unlink()
		}

		m := menu.BookmarkGroupMenu(string(c.DequeueMessage("bookmark_group")))

		m.Watch()
		c.LinkTo(m.Node)
		c.Passthrough(pubsub.RENDER, c.Next)
	})

	c.On(pubsub.END, func() {
		if c.LastEntry() {
			c.render()
		}
	})

	c.On(pubsub.HOME, func() {
		if c.FirstEntry() {
			c.render()
		}
	})

	c.On(pubsub.PG_UP, func() {
		if c.data.Len() > 0 && c.MovePgUp() {
			c.render()
		}
	})

	c.On(pubsub.PG_DOWN, func() {
		if c.MovePgDown(c.contentLines()) {
			c.render()
		}
	})

	c.On(pubsub.INSERT, func() {
		c.selectEntry()
	})

	c.On(pubsub.ARROW_DOWN, func() {
		if c.NextEntry() {
			c.render()
		}
	})

	c.On(pubsub.ARROW_UP, func() {
		if c.PrevEntry() {
			c.render()
		}
	})

	c.On(pubsub.SORT_NAME, func() {
		c.sortByName()
	})

	c.On(pubsub.SORT_TYPE, func() {
		c.sortByType()
	})

	c.On(pubsub.SORT_SIZE_ASC, func() {
		c.sortBySizeAsc()
	})

	c.On(pubsub.SORT_SIZE_DESC, func() {
		c.sortBySizeDesc()
	})

	c.On(pubsub.SETTINGS, func() {
		c.executeCmd(config.Running().DefaultEditor, config.Running().GetSettingsFilepath())
	})

	c.On(pubsub.ARROW_RIGHT, func() {
		c.right()
	})

	c.On(pubsub.ARROW_LEFT, func() {
		c.left()
	})

	c.On(pubsub.ENTER, func() {
		c.execute()
	})
	c.On(pubsub.C, func() {
		args := c.getInput("Execute command", "")

		if args == "" {
			return
		}

		path := c.getInput("Path", c.root)

		c.executeCmd(args, path)
	})

	c.On(pubsub.CTRL_F, func() {
		c.search()
	})

	c.On(pubsub.DELETE, func() {
		answ := c.getInput("Delete selected entries", "y")

		if answ != "y" || c.selected.Len() == 0 {
			return
		}

		err := c.delete()
		if err != nil {
			c.sendError(err)
		}

		c.cache.Clear()

		c.setStore(c.root)

		if c.alt.root == c.root {
			c.alt.setStore(c.root)
		}
	})

	c.On(pubsub.E, func() {
		c.edit()
	})

	c.On(pubsub.F, func() {
		c.newFile()
	})

	c.On(pubsub.G, func() {
		c.goTo()
	})

	c.On(pubsub.D, func() {
		c.newDirectory()
	})

	c.On(pubsub.CTRL_R, func() {
		c.rename()
	})

	c.On(pubsub.CTRL_A, func() {
		c.selectAll()
	})

	c.On(pubsub.SORT_MENU, func() {
		c.OpenMenu(menu.SortMenu())
	})

	c.On(pubsub.ESC, func() {
		c.OpenMenu(menu.QuitMenu())
	})

	c.On(pubsub.M, func() {
		c.OpenMenu(menu.MainMenu())
	})

	c.On(pubsub.I, func() {
		if c.data.Len() == 0 {
			return
		}

		entry, ok := c.activeEntry()

		if !ok {
			return
		}

		file, err := os.Stat(entry.FullPath())
		if err != nil {
			c.sendError(err)
		}

		if c.HasNext() {
			c.Next.Unlink()
		}

		info := boxes.NewFileInfo(file, entry.Size())

		info.Watch()
		c.LinkTo(info)
		c.Passthrough(pubsub.RENDER, c.Next)
	})

	c.On(pubsub.QUIT_APP, func() {
		c.Passthrough(pubsub.QUIT_APP, c.First())
	})

	// c.On(pubsub.Q, func() {
	// 	if c.HasNext() {
	// 		c.Next.Unlink()
	// 		c.Passthrough(pubsub.RENDER, c.Node)
	// 		// cls()
	// 		// c.restorePanels()
	// 	}
	// })

	c.On(pubsub.S, func() {
		c.scanDir(config.Running().DefaultRootDirectory)
	})

	c.OnGlobal(pubsub.T, func() {
		cls()
		c.fullRender()
		c.alt.fullRender()
	})

	c.OnGlobal(pubsub.RESIZE, func() {
		//	cls()
		//	timpubsub.Sleep(timpubsub.Millisecond * 100)

		c.UpdateSize()
		c.alt.UpdateSize()
		// c.panel.Print(themeMain())
		// c.alt.panel.Print(themeMain())
		c.fullRender()
		c.alt.fullRender()

		// c.restorePanels()
	})
}

func (c *FsController) OpenMenu(m *menu.MenuController) {
	if c.HasNext() {
		c.Next.Unlink()
	}

	m.Watch()
	c.LinkTo(m.Node)
	c.Passthrough(pubsub.RENDER, c.Next)
}
