package common

type BaseNavigator struct {
	index, prevIndex int
}

func (n *BaseNavigator) MoveUp() {
	if n.index > 0 {
		n.prevIndex = n.index
		n.index--

	}
}

func (n *BaseNavigator) MoveDown(len int) {
	if len == 0 {
		return
	}

	if n.index < len-1 {
		n.prevIndex = n.index
		n.index++

	}
}

func (n *BaseNavigator) MoveToTop(len int) {
	if n.index != 0 || len != 0 {

		n.prevIndex = n.index
		n.index = 0
	}
}

func (n *BaseNavigator) MoveToBottom(len int) {
	if len != 0 || n.index != len-1 {
		n.prevIndex = n.index

		n.index = len - 1
	}
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

	// c.SetChunk()
}
