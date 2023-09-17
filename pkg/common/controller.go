package common

type BaseController struct {
	actionMap map[ControllerAction]func()
	active    bool
}

// func (b *BaseController) Alt() Controlable {
// 	return b.altController
// }

// func (b *BaseController) SetAlt(c Controlable) {
// 	b.altController = c
// }

func (b *BaseController) ActionMap() map[ControllerAction]func() {
	return b.actionMap
}

func (b *BaseController) SetActionMap(m map[ControllerAction]func()) {
	b.actionMap = m
}

func (b *BaseController) Map(action ControllerAction) {
	if fn, ok := b.ActionMap()[action]; ok {
		fn()
	}
}

func (b *BaseController) Active() bool {
	return b.active
}

func (b *BaseController) SetActive(active bool) {
	b.active = active
}
