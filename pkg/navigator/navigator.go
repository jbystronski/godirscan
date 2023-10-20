package navigator

type BaseNavigator struct {
	index, prevIndex int
}

func (n *BaseNavigator) MoveUp() bool {
	if n.index > 0 {
		n.prevIndex = n.index
		n.index--
		return true

	}

	return false
}

func (n *BaseNavigator) MoveDown(len int) bool {
	if len == 0 {
		return false
	}

	if n.index < len-1 {
		n.prevIndex = n.index
		n.index++
		return true

	}

	return false
}

func (n *BaseNavigator) MoveToTop(len int) bool {
	if n.index == 0 || len == 0 {
		return false
	}

	n.prevIndex = n.index
	n.index = 0
	return true
}

func (n *BaseNavigator) MoveToBottom(len int) bool {
	if len == 0 || n.index == len-1 {
		return false
	}

	n.prevIndex = n.index

	n.index = len - 1
	return true
}

func (n *BaseNavigator) Index() int {
	return n.index
}

func (n *BaseNavigator) SetIndex(i int) {
	n.prevIndex = n.index
	n.index = i
}

func (n *BaseNavigator) PrevIndex() int {
	return n.prevIndex
}

func (n *BaseNavigator) SetPrevIndex(i int) {
	n.prevIndex = i
}

func (n *BaseNavigator) Reset() {
	n.index = 0
	n.prevIndex = 0
}
